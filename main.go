package main

import (
	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var HostPort = "127.0.0.1:7933"
var Domain = "samples-domain"
var TaskListName = "helloworld-worker"
var ClientName = "helloworld-worker"
var CadenceService = "cadence-frontend"

func main() {
	startWorkers(buildLogger(), buildCadenceClient())
	select {}

}

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		panic("Failed to setup logger")
	}

	return logger
}

func buildCadenceClient() workflowserviceclient.Interface {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService))
}

func startWorkers(logger *zap.Logger, service workflowserviceclient.Interface) {
	// TaskListName identifies set of client workflows, activities, and workers.
	// It could be your group or client or application name.
	workerOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(TaskListName, map[string]string{}),
	}
	// DisableWorkflowWorker: false,
	// DisableActivityWorker: true,

	workerOptions.DisableWorkflowWorker = false
	workerOptions.Identity = "workflow-worker"
	workflowWorker := worker.New(
		service,
		Domain,
		TaskListName,
		workerOptions)
	err := workflowWorker.Start()
	if err != nil {
		panic("Failed to start worker")
	}

	workerOptions.DisableActivityWorker = false
	workerOptions.DisableWorkflowWorker = true
	workerOptions.Identity = "activity-worker"
	activityWorker := worker.New(
		service,
		Domain,
		TaskListName,
		workerOptions)

	err = activityWorker.Start()
	if err != nil {
		panic("Failed to start worker")
	}
	logger.Info("Started Workers.", zap.String("worker", TaskListName))
}
