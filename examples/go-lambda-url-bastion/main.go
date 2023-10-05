package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/clients/go/funcie-tunnel"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type Response struct {
	Greeting string `json:"greeting"`
}

func HandleRequest(_ context.Context, event events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
	name := event.QueryStringParameters["name"]
	resp := Response{
		Greeting: fmt.Sprintf("Hello %s -- yay! :)", name),
	}

	if strings.ToLower(event.QueryStringParameters["error"]) == "true" {
		return &events.LambdaFunctionURLResponse{}, errors.New("error :(")
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return &events.LambdaFunctionURLResponse{}, err
	}

	return &events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	funcie.ConfigureLogging()
	funcie_tunnel.Start(HandleRequest)
}
