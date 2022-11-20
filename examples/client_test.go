package examples

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/m-mizutani/gurl"
)

type testClient struct{}

func (x *testClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Println(req.URL)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil
}

func ExampleWithClient() {
	client := &testClient{}

	if err := gurl.Get("https://www.google.com/search",
		gurl.WithParam("q", "security"),
		gurl.WithClient(client),
	); err != nil {
		log.Fatal(err.Error())
	}

	//Output: https://www.google.com/search?q=security
}

type testServer struct{}

func (x testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.WriteHeader(200)
	w.Write([]byte("hello"))
}

func ExampleWithHandlerFunc() {
	server := &testServer{}

	if err := gurl.Get("https://www.google.com/search",
		gurl.WithParam("q", "security"),
		gurl.WithHandlerFunc(server.ServeHTTP),
	); err != nil {
		log.Fatal(err.Error())
	}

	//Output: https://www.google.com/search?q=security
}
