package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	redistransport "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/aws/aws-lambda-go/events"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
)

type Response struct {
	Greeting string `json:"greeting"`
}

func HandleRequest(_ context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	name := event.QueryStringParameters["name"]
	resp := Response{
		Greeting: fmt.Sprintf("Hello %s -- yay!", name),
	}

	if strings.ToLower(event.QueryStringParameters["error"]) == "true" {
		return events.LambdaFunctionURLResponse{}, errors.New("error :(")
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
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("FUNCIE_REDIS_ADDR"),
	})
	publisher := redistransport.NewPublisher(redisClient, "funcie:requests")
	consumer := redistransport.NewConsumer(redisClient, "funcie:requests")
	tunnel := funcie.NewLambdaTunnel("lambda-url-lib", HandleRequest, publisher, consumer)
	tunnel.Start()
}
