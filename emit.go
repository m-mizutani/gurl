package gurl

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

func (x *Request) Emit() error {
	x.logger.Debug("building HTTP request",
		slog.Any("params", x.params),
		slog.Any("header", x.headers),
		slog.Int("expected code", x.expectedCode),
	)

	baseURL, err := url.Parse(x.url)
	if err != nil {
		return goerr.Wrap(err, "parsing URL")
	}

	// Add query strings
	query := baseURL.Query()
	for key, values := range x.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	baseURL.RawQuery = query.Encode()

	// Setting up request body
	body, contentType, err := x.body()
	if err != nil {
		return err
	}

	if contentType != "" {
		hasContentType := func(headers map[string][]string) bool {
			for hdrKey := range headers {
				if strings.ToLower(hdrKey) == "content-type" {
					return true
				}
			}
			return false
		}
		if !hasContentType(x.headers) {
			x.headers["Content-Type"] = []string{contentType}
		}
	}

	// Create http.Request
	var req *http.Request
	if x.ctx == nil {
		req, err = http.NewRequest(methodMap[x.method], baseURL.String(), body)
	} else {
		req, err = http.NewRequestWithContext(x.ctx, methodMap[x.method], baseURL.String(), body)
	}

	if err != nil {
		return goerr.Wrap(err, "creating HTTP request")
	}

	// Add HTTP Headers
	for key, values := range x.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := x.client.Do(req)
	if err != nil {
		return goerr.Wrap(err, "sending HTTP request")
	}
	defer resp.Body.Close()

	if x.expectedCode > 0 && resp.StatusCode != x.expectedCode {
		body, _ := io.ReadAll(resp.Body)
		return goerr.Wrap(ErrUnexpectedCode).
			With("code", resp.StatusCode).
			With("body", body)
	}

	// Handling HTTP response body
	if x.reader != nil {
		if err := x.reader(resp.Body); err != nil {
			return err
		}
	}

	return nil
}
