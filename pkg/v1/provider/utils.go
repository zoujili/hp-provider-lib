package provider

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

// Utility function that allows waiting for a provider to run.
// Mostly usable from other providers that have a dependency.
func WaitForRunningProvider(p RunProvider, timeoutSeconds time.Duration) error {
	if p.IsRunning() {
		// No need to wait if provider is already running.
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	name := Name(p)
	logrus.WithField("timeout", timeoutSeconds).Debugf("Waiting for %s to run...", name)
	for {
		if p.IsRunning() {
			return nil
		}
		if ctx.Err() != nil {
			return fmt.Errorf("time exceeded for %s to run", name)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Utility function to get a provider's name.
func Name(provider Provider) string {
	return reflect.TypeOf(provider).Elem().String()
}
