package gurl

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"golang.org/x/exp/slog"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Request struct {
	method Method
	url    string

	ctx          context.Context
	expectedCode int
	headers      map[string][]string
	params       map[string][]string

	body   RequestBody
	reader ResponseReader

	client Client

	logger *slog.Logger
}

// New creates a new [Request] with default values and configures [Request] with options.
func New(method Method, url string, options ...Option) *Request {
	req := &Request{
		expectedCode: http.StatusOK,

		method:  method,
		url:     url,
		headers: make(map[string][]string),
		params:  make(map[string][]string),
		body:    func() (io.Reader, string, error) { return nil, "", nil },
		client:  http.DefaultClient,

		logger: slog.New(slog.HandlerOptions{
			Level: slog.ErrorLevel,
		}.NewTextHandler(os.Stdout)),
	}

	for _, opt := range options {
		opt(req)
	}

	return req
}

type Option func(req *Request)

// WithCtx will set given [context.Context] to [http.Request].
func WithCtx(ctx context.Context) Option {
	return func(req *Request) {
		req.ctx = ctx
	}
}

// WithExpectedCode sets expected status code. If returned status code from server is not matched with *statusCode*, Emit method will return ErrUnexpectedCode. Default value is 200 (HTTP OK). If you don't want to handle status code, set 0 and Emit method will return no error with any status code.
func WithExpectedCode(statusCode int) Option {
	return func(req *Request) {
		req.expectedCode = statusCode
	}
}

// WithHeader adds a pair of key and value as HTTP header. This option can be set multiply.
func WithHeader(key, value string) Option {
	return func(req *Request) {
		req.headers[key] = append(req.headers[key], value)
	}
}

// WithParam adds a pair of key and value for Query String Parameter. This method will not overwrite existing query string in original URL.
func WithParam(key, value string) Option {
	return func(req *Request) {
		req.params[key] = append(req.params[key], value)
	}
}

// WithBody sets body of HTTP request as [io.Reader].
func WithBody(body RequestBody) Option {
	return func(req *Request) {
		req.body = body
	}
}

// WithClient sets [Client] interface to handle [http.Request].
func WithClient(client Client) Option {
	return func(req *Request) {
		req.client = client
	}
}

type testClient struct {
	handler http.HandlerFunc
}

func (x *testClient) Do(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	x.handler(resp, req)
	return resp.Result(), nil
}

// WithHandlerFunc is a option for testing. The f [http.HandlerFunc] will be injected as [Client] and record response with [httptest.ResponseRecorder].
func WithHandlerFunc(f http.HandlerFunc) Option {
	return func(req *Request) {
		req.client = &testClient{
			handler: f,
		}
	}
}

// WithLogger sets *[slog.Logger] as logging interface.
func WithLogger(logger *slog.Logger) Option {
	return func(req *Request) {
		req.logger = logger
	}
}

// WithReader is a handler of response body. It sets a functions to be called with Body of [http.Response]. The function must call r.Close() to terminate communication explicitly.
func WithReader(f ResponseReader) Option {
	return func(req *Request) {
		req.reader = f
	}
}
