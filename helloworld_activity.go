package main

import (
	"context"

	"github.com/hr1sh1kesh/hwworkflow/wf"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
)

func init() {
	workflow.Register(wf.HelloWorldWorkflow)
	activity.Register(HelloworldActivity)
}

func HelloworldActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "!", nil
}
