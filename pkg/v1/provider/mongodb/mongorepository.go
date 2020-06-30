package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/middleware/tenant"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetDatabaseNameFromTenantID func(tenantID string) string

func DefalutDatabaseNameGetter(databaseNameSuffix string) GetDatabaseNameFromTenantID {
	return func(tenantID string) string {
		return fmt.Sprintf("%s-%s", tenantID, databaseNameSuffix)
	}
}

type IMongoRepository interface {
	provider.Provider
	MongoClient() *mongo.Client
	MongoDatabase(ctx context.Context) *mongo.Database
	RunTransaction(ctx context.Context, txnFunc func(sessionContext mongo.SessionContext) error) error
}

type MongoRepository struct {
	provider.AbstractProvider
	mongoProvider      *MongoDB
	databaseNameGetter GetDatabaseNameFromTenantID
}

func NewMongoRepositoryWithDatabaseNameGetter(mongoProvider *MongoDB, databaseNameGetter GetDatabaseNameFromTenantID) IMongoRepository {
	return &MongoRepository{
		mongoProvider:      mongoProvider,
		databaseNameGetter: databaseNameGetter,
	}
}

func NewMongoRepository(mongoProvider *MongoDB) IMongoRepository {
	return &MongoRepository{
		mongoProvider: mongoProvider,
	}
}

func (m MongoRepository) MongoClient() *mongo.Client {
	return m.mongoProvider.Client
}

func (m MongoRepository) MongoDatabase(ctx context.Context) *mongo.Database {
	// multi-tenant support
	if tenantID, ok := tenant.FromTenantInterceptorContext(ctx); ok && m.databaseNameGetter != nil {
		return m.mongoProvider.Client.Database(m.databaseNameGetter(tenantID))
	}

	// MONGODB_DATABASE set
	if len(os.Getenv("MONGODB_DATABASE")) > 0 {
		return m.mongoProvider.Database
	}

	// MONGODB_URI  contains database_name
	uri, err := connstring.Parse(m.mongoProvider.Config.URI)
	if err == nil && len(uri.Database) > 0 {
		return m.mongoProvider.Client.Database(uri.Database)
	}

	panic(status.Error(codes.Internal, "No database specified"))
}

func (m MongoRepository) RunTransaction(ctx context.Context, txnFn func(mongo.SessionContext) error) error {
	return m.mongoProvider.Client.UseSessionWithOptions(ctx,
		options.Session().SetDefaultReadPreference(readpref.Primary()), func(sessionContext mongo.SessionContext) error {
			return m.runTransactionWithRetry(sessionContext, txnFn)
		})
}

func (m MongoRepository) runTransactionWithRetry(sctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
	txID := uuid.NewV4().String()
	log.Println("Begin transaction: ", txID)
	for {
		// start transaction
		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			log.Printf("Start transaction: %s. Error: %s\n", txID, err)
			return err
		}

		err = txnFn(sctx)
		if err != nil {
			log.Printf("Abort transaction: %s. Error: %s\n", txID, err)
			if abortErr := sctx.AbortTransaction(sctx); abortErr != nil {
				log.Println("Failed to abort transaction: ", txID)
			}
			return err
		}

		err = m.commitWithRetry(sctx)
		switch e := err.(type) {
		case nil:
			log.Printf("End transaction: %s, successful. \n", txID)
			return nil
		case mongo.CommandError:
			// If transient error, retry the whole transaction
			if e.HasErrorLabel("TransientTransactionError") {
				log.Printf("TransientTransactionError: %s, retrying transaction...\n", txID)
				continue
			}
			return e
		default:
			log.Printf("End transaction: %s. Error: %s\n", txID, err)
			return e
		}
	}
}

func (m MongoRepository) commitWithRetry(sctx mongo.SessionContext) error {
	for {
		err := sctx.CommitTransaction(sctx)
		switch e := err.(type) {
		case nil:
			return nil
		case mongo.CommandError:
			// Can retry commit
			if e.HasErrorLabel("UnknownTransactionCommitResult") {
				continue
			}
			return e
		default:
			return e
		}
	}
}
