package main

import (
	"context"
	"encoding/json"
	"fmt"
	"funcie/tunnel"
	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	Greeting string `json:"greeting"`
}

func HandleRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	resp := Response{
		Greeting: fmt.Sprintf("Hello %s!", event.QueryStringParameters["name"]),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		return events.LambdaFunctionURLResponse{}, err
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	tunnel.Start(HandleRequest)
}
