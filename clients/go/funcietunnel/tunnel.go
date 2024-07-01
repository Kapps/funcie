package funcietunnel

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"log/slog"
	"time"
)

// Start is a replacement to lambda.Start that configures the proxy from environment variables and SSM.
// It will start the proxy if the application is running in a Lambda, or start the receiver if it is running locally.
// The handler is the handler that will be invoked when a request is received.
// It is subject to the same restrictions as the handler for the underlying serverless function provider (such as lambda.Start).
// See NewConfig for the environment variables that are used.
// The application name is an arbitrary identifier to uniquely identify this application in order to route messages.
func Start(appName string, handler interface{}) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	fmt.Println("Loading AWS config")
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %s", err))
	}

	fmt.Println("Creating Funcie config")
	ssmClient := ssm.NewFromConfig(conf)
	fmt.Println("Created SSM client")
	funcieConfig := NewConfig(ctx, appName, ssmClient)
	fmt.Println("Created Funcie config")
	StartWithConfig(*funcieConfig, slog.Default(), handler)
}

// StartWithConfig is a replacement to lambda.Start that configures the proxy from the given config.
// See `Start` for more information.
func StartWithConfig(config FuncieConfig, logger *slog.Logger, handler interface{}) {
	if funcie.IsRunningWithLambda() {
		// In a Lambda, we wait for the Lambda runtime to call the handler and forward that request to the bastion.
		client := NewHTTPBastionClient(config.ServerBastionEndpoint, logger)
		proxy := NewLambdaFunctionProxy(config.ApplicationId, client, handler, logger)
		proxy.Start()
	} else {
		// Locally, we receive the request from the bastion.
		receiver := NewLambdaBastionReceiver(
			config.ApplicationId,
			config.ListenAddress,
			config.ClientBastionEndpoint,
			handler,
			logger,
		)
		receiver.Start()
	}
}
