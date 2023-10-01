package main

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/examples/go-grpc-bastion/funciegrpc"
	"github.com/Kapps/funcie/pkg/funcie"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"time"
)

type Response struct {
	Greeting string `json:"greeting"`
}

type funcieGrpc interface {
	grpc.ServiceRegistrar
	Serve(net.Listener) error
}

type exampleServer struct {
	funciegrpc.UnimplementedGreeterServer
}

type registrar struct {
}

func (r *registrar) RegisterService(desc *grpc.ServiceDesc, impl any) {
	fmt.Printf("RegisterService %T: %+v\n", desc.HandlerType, desc.HandlerType)
}

func (r *registrar) Serve(listener net.Listener) error {
	return nil
}

func (s *exampleServer) SayHello(ctx context.Context, req *funciegrpc.HelloRequest) (*funciegrpc.HelloReply, error) {
	return &funciegrpc.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil

	/*name := event.QueryStringParameters["name"]
	resp := Response{
		Greeting: fmt.Sprintf("Hello %s -- yay! :)", name),
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
	}, nil*/
}

func NewFuncieUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		jsonBytes := funcie.MustSerialize(req)
		fmt.Printf("FuncieUnaryInterceptorFunc: %v\n", string(jsonBytes))
		//fmt.Printf("FuncieUnaryInterceptorFuncInfo: %+v\n", info.Server.(grpc.ClientConnInterface))
		//md := server.GetServiceInfo()[info.FullMethod].Metadata
		//fmt.Printf("FuncieUnaryInterceptorFuncMetadata: %+v\n", md)

		return handler(ctx, req)
	}
}

func makeServer(isLocal bool) funcieGrpc {
	if isLocal {
		return &registrar{}
	} else {
		return grpc.NewServer(grpc.ChainUnaryInterceptor(NewFuncieUnaryInterceptor()))
	}
}

func main() {
	port := 24191

	isLocal := os.Getenv("IS_LOCAL") == "1"
	s := makeServer(isLocal)
	funciegrpc.RegisterGreeterServer(s, &exampleServer{})
	if isLocal {
		log.Printf("server _not_ listening.")
	} else {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("server listening at %v", lis.Addr())
		go func() {
			time.Sleep(time.Second)
			dc, err := grpc.DialContext(
				context.Background(),
				fmt.Sprintf("localhost:%d", port),
				grpc.WithCredentialsBundle(insecure.NewBundle()),
			)
			if err != nil {
				log.Fatalf("failed to dial: %v", err)
			}
			client := funciegrpc.NewGreeterClient(dc)
			resp, err := client.SayHello(context.Background(), &funciegrpc.HelloRequest{Name: "World"})
			if err != nil {
				log.Fatalf("failed to call: %v", err)
			}
			log.Printf("response: %+v", resp)
		}()
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
	//funcie.ConfigureLogging()
	//funcie_tunnel.Start(HandleRequest)
}
