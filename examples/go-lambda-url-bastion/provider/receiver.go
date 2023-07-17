package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/exp/slog"
	"io"
	"net"
	"net/http"
	"net/url"
)

// BastionReceiver represents a receiver that can be used to receive requests from a bastion.
type BastionReceiver interface {
	// Start starts the tunnel. This function never returns.
	// The handler is the handler that will be invoked when a request is received.
	// It is subject to the same restrictions as the handler for the serverless function provider (such as lambda.Start).
	Start()
}

type bastionReceiver struct {
	applicationId   string
	bastionEndpoint url.URL
	listenAddress   string
	server          *http.Server
	client          *http.Client
	handler         lambda.Handler
}

// NewLambdaBastionReceiver creates a new BastionReceiver for AWS Lambda operations.
func NewLambdaBastionReceiver(
	applicationId string,
	listenAddress string,
	bastionEndpoint url.URL,
	handler interface{},
) BastionReceiver {
	lambdaHandler := lambda.NewHandler(handler)

	return &bastionReceiver{
		applicationId:   applicationId,
		bastionEndpoint: bastionEndpoint,
		handler:         lambdaHandler,
		listenAddress:   listenAddress,
		client:          &http.Client{},
		server: &http.Server{
			Addr: listenAddress,
		},
	}
}

func (r *bastionReceiver) Start() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(r.handleRequest))
	r.server.Handler = mux

	// Before we start the server, we need to subscribe to the bastion.
	listener, err := net.Listen("tcp", r.listenAddress)
	if err != nil {
		slog.Error("failed to listen", err)
		panic(err)
	}

	err = r.subscribe(listener.Addr())
	if err != nil {
		slog.Error("failed to subscribe", err)
		panic(err)
	}

	slog.Info("starting bastion receiver", "applicationId", r.applicationId, "listenAddress", r.server.Addr)

	err = r.server.ListenAndServe()
	slog.Warn("server stopped", "err", err)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (r *bastionReceiver) subscribe(addr net.Addr) error {
	localEndpoint := funcie.MustNewEndpointFromAddress(fmt.Sprintf("http://%s/", addr))
	registerEndpoint := fmt.Sprintf("%s/dispatch", r.bastionEndpoint.String())
	payload := messages.NewRegistrationRequestPayload(r.applicationId, localEndpoint)
	message := funcie.NewMessageWithPayload(r.applicationId, messages.MessageKindRegister, payload)

	marshaled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	slog.Info("sending registration request",
		"message", message, "bastionEndpoint", r.bastionEndpoint.String(), "registerEndpoint", registerEndpoint)

	resp, err := r.client.Post(registerEndpoint, "application/json", bytes.NewReader(marshaled))
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	var response funcie.Message
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	slog.Info("received response", "response", response)

	return nil
}

func (r *bastionReceiver) handleRequest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to read request body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message funcie.Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to unmarshal request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.InfoCtx(ctx, "received request", "message", message)
	if message.Kind != messages.MessageKindForwardRequest {
		slog.WarnCtx(ctx, "received message with invalid kind", "kind", message.Kind)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unmarshaled, err := funcie.UnmarshalMessagePayload[messages.ForwardRequestMessage](&message)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to unmarshal request message", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload := []byte(unmarshaled.Payload.Body)
	invokeResponse, err := r.handler.Invoke(ctx, payload)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to handle message", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := messages.NewForwardRequestResponsePayload(invokeResponse)

	slog.InfoCtx(ctx, "sending response", "response", response)
	responseBody, err := json.Marshal(response)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to marshal response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to write response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.InfoCtx(ctx, "sent response")
}
