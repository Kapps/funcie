package bastion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"strconv"
)

// Server is the interface for a bastion server that can filter and forward requests.
type Server interface {
	// Listen starts the server and begins processing requests.
	// This function never returns unless an error occurs.
	Listen() error
}

type ResultCode int

var (
	// ResultCodeForwardSuccess indicates that the request was handled successfully and a response was forwarded back to the caller.
	ResultCodeForwardSuccess ResultCode = 1
	// ResultCodeNoConsumers indicates that the request was not forwarded to the target.
	ResultCodeNoConsumers ResultCode = 2
	// ResultCodeInternalError indicates that an internal error occurred while processing the request.
	ResultCodeInternalError ResultCode = 3
	// ResultCodeInvalidRequest indicates that the request was invalid.
	ResultCodeInvalidRequest ResultCode = 4
)

type server struct {
	httpServer *http.Server
	handler    RequestHandler
}

// NewServer creates a new Server along with a backing HTTP server.
func NewServer(address string, handler RequestHandler) Server {
	return &server{
		httpServer: &http.Server{
			Addr: address,
		},
		handler: handler,
	}
}

// NewServerWithHTTPServer creates a new Server with the given HTTP server instead of creating one.
func NewServerWithHTTPServer(httpServer *http.Server, handler RequestHandler) Server {
	return &server{
		httpServer: httpServer,
		handler:    handler,
	}
}

// Request is an incoming request from the caller to the bastion server.
type Request struct {
	// RequestId is a caller-specified unique ID for this request (for example, the request ID of a Lambda invocation).
	RequestId string `json:"requestId"`
	// MessageKind is the kind of message that is being forwarded.
	MessageKind funcie.MessageKind `json:"messageKind"`
	// Application is the name of the application that is making the request.
	Application string `json:"application"`
	// Payload is the JSON payload that is potentially being forwarded.
	Payload *json.RawMessage `json:"payload"`
	// RequestParameters are any specific parameters that the caller wants to pass to the bastion server.
	// This can be used for purposes such as filtering to only forward certain types of requests.
	RequestParameters map[string]string
}

// NewRequest creates a new Request.
func NewRequest(application string, payload *json.RawMessage, requestParameters map[string]string) *Request {
	return &Request{
		RequestId:         uuid.New().String(),
		Payload:           payload,
		RequestParameters: requestParameters,
		Application:       application,
	}
}

// Response is the response from the bastion server sent back to the caller.
type Response struct {
	Error      string      `json:"error,omitempty"`
	ResultCode ResultCode  `json:"resultCode"`
	Data       interface{} `json:"data,omitempty"`
}

// NewErrorResponse creates a new Response that indicates an error.
func NewErrorResponse(resultCode ResultCode, error string) *Response {
	return &Response{
		Error:      error,
		Data:       nil,
		ResultCode: resultCode,
	}
}

// NewDataResponse creates a new Response that indicates a processed message.
func NewDataResponse(resultCode ResultCode, data interface{}) *Response {
	return &Response{
		Error:      "",
		Data:       data,
		ResultCode: resultCode,
	}
}

func (s *server) Listen() error {
	s.httpServer.Handler = http.HandlerFunc(s.handleRequest)
	if err := s.httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("failed to listen or serve: %w", err)
	}

	return nil
}

func (s *server) handleRequest(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "method", r.Method, "path", r.URL.Path)

	if r.Method != http.MethodPost {
		writeErrorResult(w, http.StatusMethodNotAllowed, ResultCodeInvalidRequest, "invalid method")
		return
	}

	if r.URL.Path != "/dispatch" {
		writeErrorResult(w, http.StatusNotFound, ResultCodeInvalidRequest, "invalid path")
		return
	}

	ctx := r.Context()
	if err := s.handleDispatchRequest(ctx, w, r); err != nil {
		writeErrorResult(w, http.StatusInternalServerError, ResultCodeInternalError, err.Error())
		return
	}
}

func (s *server) handleDispatchRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	bodyContents, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	var request Request
	if err := json.Unmarshal(bodyContents, &request); err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}

	slog.Info("received request", "requestId", request.RequestId, "payload", string(*request.Payload))

	res, err := s.handler.Dispatch(ctx, &request)
	if errors.Is(err, funcie.ErrNoActiveConsumer) {
		slog.Debug("no consumers found for request", "requestId", request.RequestId)

		resp := NewDataResponse(ResultCodeNoConsumers, nil)
		if err := writeSuccessResponse(w, resp); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}

		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to dispatch request: %w", err)
	}

	slog.Info("dispatch request succeeded", "requestId", request.RequestId)

	resp := NewDataResponse(ResultCodeForwardSuccess, res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if err := writeSuccessResponse(w, resp); err != nil {
		return fmt.Errorf("failed to write success response: %w", err)
	}

	return nil
}

func writeSuccessResponse(w http.ResponseWriter, response *Response) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		slog.Error("failed to write response", err)
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func writeErrorResult(w http.ResponseWriter, statusCode int, resultCode ResultCode, message string) {
	resp := errorResponse(resultCode, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(resp)
	if err != nil {
		slog.Error("failed to write error response", err)
	}
}

func errorResponse(resultCode ResultCode, message string) []byte {
	response := NewErrorResponse(resultCode, message)
	data, err := json.Marshal(response)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal error response: %v", err))
	}
	return data
}
