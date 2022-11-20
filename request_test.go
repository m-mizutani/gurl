package gurl_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/m-mizutani/gurl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testURL = "https://emhkq5vqrco2fpr6zqlctbjale0eyygt.lambda-url.ap-northeast-1.on.aws"
)

type testServerResponse struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Params  map[string]string `json:"params"`
	Body    []byte            `json:"body"` // []byte in json will decode base64
}

func TestGetRequest(t *testing.T) {
	var resp testServerResponse
	require.NoError(t, gurl.Get(testURL+"/hello_lambda",
		gurl.WithExpectedCode(http.StatusOK),
		gurl.WithHeader("Gurl-Test", "blue"),
		gurl.WithReader(gurl.DecodeAsJSON(&resp)),
	))

	assert.Equal(t, "GET", resp.Method)
	assert.Equal(t, "/hello_lambda", resp.Path)
	assert.Equal(t, "blue", resp.Headers["gurl-test"]) // AWS API gateway converts HTTP header key to lower case.
}

func TestPostRequest(t *testing.T) {
	var resp testServerResponse
	require.NoError(t, gurl.Post(testURL+"/hello_lambda",
		gurl.WithBody(gurl.ByReader(bytes.NewReader([]byte("orange")))),
		gurl.WithExpectedCode(http.StatusOK),
		gurl.WithReader(gurl.DecodeAsJSON(&resp)),
	))

	assert.Equal(t, "POST", resp.Method)
	assert.Equal(t, "/hello_lambda", resp.Path)
	assert.Equal(t, "orange", string(resp.Body))
}

func Test(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
