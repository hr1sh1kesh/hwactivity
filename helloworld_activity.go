package main

import (
	"context"

	"go.uber.org/cadence/activity"
)

func init() {
	activity.Register(HelloworldActivity)

}

func HelloworldActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "!", nil
}
