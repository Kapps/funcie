package funcie_tunnel

import (
	"github.com/Kapps/funcie/pkg/funcie"
)

// Start is a replacement to lambda.Start that configures the proxy from environment variables.
// It will start the proxy if the application is running in a Lambda, or start the receiver if it is running locally.
// The handler is the handler that will be invoked when a request is received.
// It is subject to the same restrictions as the handler for the underlying serverless function provider (such as lambda.Start).
// See NewConfigFromEnvironment for the environment variables that are used.
func Start(handler interface{}) {
	config := NewConfigFromEnvironment()
	StartWithConfig(*config, handler)
}

// StartWithConfig is a replacement to lambda.Start that configures the proxy from the given config.
// See `Start` for more information.
func StartWithConfig(config FuncieConfig, handler interface{}) {
	if funcie.IsRunningWithLambda() {
		// In a Lambda, we wait for the Lambda runtime to call the handler and forward that request to the bastion.
		client := NewHTTPBastionClient(config.ServerBastionEndpoint)
		proxy := NewLambdaFunctionProxy(config.ApplicationId, client, handler)
		proxy.Start()
	} else {
		// Locally, we receive the request from the bastion.
		receiver := NewLambdaBastionReceiver(config.ApplicationId, config.ListenAddress, config.ClientBastionEndpoint, handler)
		receiver.Start()
	}
}
