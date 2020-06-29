package migrate

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/mongodb"
	"github.com/golang-migrate/migrate/v4"
	mongo "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file" // mattes migrate
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Migrate struct {
	provider.AbstractProvider
	Config  *Config
	Mongodb *mongodb.MongoDB
}

func New(config *Config, mongodb *mongodb.MongoDB) *Migrate {
	return &Migrate{
		Mongodb: mongodb,
		Config:  config,
	}
}

func (s *Migrate) Init() (err error) {
	directory := s.Config.Directory
	logrus.Infof("Run migrations under directory: %s", directory)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		logrus.Errorf("Failed to list migrations: %v", err)
		return err
	}
	logrus.Debug("Found migrations:")
	for _, f := range files {
		logrus.Debug(f.Name())
	}
	mongoCfg, err := s.getMongoConfig()
	if err != nil {
		logrus.Errorf("Failed to get mongo config: %v", err)
		return err
	}
	db, err := mongo.WithInstance(s.Mongodb.Client, mongoCfg)
	if err != nil {
		logrus.Errorf("Failed to WithInstance: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", directory), "", db)
	if err != nil {
		logrus.Errorf("Failed to NewWithDatabaseInstance: %v", err)
		return err
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logrus.Info("Schema already up to date")
		} else {
			logrus.Errorf("Failed to migrate: %v", err)
			return err
		}
	}
	logrus.Info("Schema migrated successfully!")
	return nil
}

func (s *Migrate) getMongoConfig() (*mongo.Config, error) {
	uri, err := connstring.Parse(s.Mongodb.Config.URI)
	if err != nil {
		return nil, err
	}
	var databaseName = uri.Database
	if len(uri.Database) == 0 {
		databaseName = s.Mongodb.Config.Database
	}
	unknown := url.Values(uri.UnknownOptions)
	migrationsCollection := unknown.Get("x-migrations-collection")
	transactionMode, _ := strconv.ParseBool(unknown.Get("x-transaction-mode"))

	return &mongo.Config{
		DatabaseName:         databaseName,
		MigrationsCollection: migrationsCollection,
		TransactionMode:      transactionMode,
	}, nil
}
