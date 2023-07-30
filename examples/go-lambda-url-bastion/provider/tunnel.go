package provider

import (
	"github.com/Kapps/funcie/pkg/funcie"
)

// Start is a replacement to lambda.Start that configures the proxy from environment variables.
func Start(handler interface{}) {
	if funcie.IsRunningWithLambda() {
		// In a Lambda, we wait for the Lambda runtime to call the handler and forward that request to the bastion.
		config := NewConfigFromEnvironment()

		client := NewHTTPBastionClient(config.BastionEndpoint)
		proxy := NewLambdaFunctionProxy(config.ApplicationId, client, handler)
		proxy.Start()
	} else {
		// Locally, we receive the request from the bastion.
	}
}
