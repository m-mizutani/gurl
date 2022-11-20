package examples

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gurl"
)

type errorClient struct{}

func (x *errorClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewReader([]byte("crashed!"))),
	}, nil
}

func ExampleError() {
	if err := gurl.Get("https://example.com",
		gurl.WithClient(&errorClient{}),
	); err != nil {
		var goErr *goerr.Error
		if errors.As(err, &goErr) {
			fmt.Println(string(goErr.Values()["body"].([]byte)))
		}
	}
	//Output: crashed!
}
