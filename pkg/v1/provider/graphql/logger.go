package graphql

import (
	"context"
	"github.com/graph-gophers/graphql-go/log"
	"github.com/sirupsen/logrus"
)

// GraphQL-Go logger that uses Logrus.
type graphqlLogger struct {
	log.Logger
}

// Performs Panic logging.
func (l *graphqlLogger) LogPanic(_ context.Context, value interface{}) {
	logrus.Panicf("GraphQL panic occurred: %v", value)
}
