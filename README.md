# gURL: Go URL Request Library <!-- omit in toc --> [![Go Reference](https://pkg.go.dev/badge/github.com/m-mizutani/gurl.svg)](https://pkg.go.dev/github.com/m-mizutani/gurl) [![test](https://github.com/m-mizutani/gurl/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/gurl/actions/workflows/test.yml) [![gosec](https://github.com/m-mizutani/gurl/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/gurl/actions/workflows/gosec.yml) [![trivy](https://github.com/m-mizutani/gurl/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/gurl/actions/workflows/trivy.yml)

Make a URL request in Go easy.

```go
	var resp struct {
		KeysURL string `json:"keys_url"`
	}

	if err := gurl.Get("https://api.github.com/",
		gurl.WithExpectedCode(200),
		gurl.WithReader(gurl.DecodeAsJSON(&resp)),
	); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(resp.KeysURL)
	//Output: https://api.github.com/user/keys
```

## Motivation <!-- omit in toc -->

Go official [net/http](https://pkg.go.dev/net/http) package is simple and powerful to control HTTP communication. However we need to write similar code (creating a request, error handling, decoding a response, etc.) in most cases. Major use cases should be able to implement easier like [requests](https://requests.readthedocs.io/en/latest/) in Python.

## How-to based Usages <!-- omit in toc -->

- [Passing body data](#passing-body-data)
- [Passing query string parameters in URL](#passing-query-string-parameters-in-url)
- [Custom Headers](#custom-headers)
- [JSON Response Content](#json-response-content)
- [Dump response body](#dump-response-body)
- [With context.Context](#with-contextcontext)
- [Error handling](#error-handling)
- [Testing with your http.Client or http.HandlerFunc](#testing-with-your-httpclient-or-httphandlerfunc)

### Passing body data

`WithBody()` and `ByReader()` simply set io.Reader to read HTTP request body.

```go
	r := bytes.NewReader([]byte("my_test_data"))

	if err := gurl.Post(serverURL,
		gurl.WithBody(gurl.ByReader(r)),
	); err != nil {
		log.Fatal(err.Error())
	}
```

`EncodeAsJSON()` supports JSON encoding from struct data to byte data and create a new `io.Reader`.

```go
	var req struct {
		Name string `json:"name"`
	}

	if err := gurl.Post(serverURL,
		gurl.WithBody(gurl.EncodeAsJSON(req)),
	); err != nil {
		log.Fatal(err.Error())
	}
```

Also, gurl has `EncodeAsURL()` to support URL encoding.

```go
	body := map[string]any{
		"color": "blue",
		"text":  "Hello Günter",
	}

	if err := gurl.Post(serverURL,
		// It will send "color=blue&text=Hello+G%C3%BCnter" as request body
		gurl.WithBody(gurl.EncodeAsURL(body)),
	); err != nil {
		log.Fatal(err.Error())
	}
```

### Passing query string parameters in URL

`WithParam()` option adds query parameter to URL.

```go
	// It will send request to https://www.google.com/search?q=security
	if err := gurl.Get("https://www.google.com/search",
		gurl.WithParam("q", "security"),
	); err != nil {
		log.Fatal(err.Error())
	}
```

If original URL has query string, `WithParam()` keeps original query parameters and appends new parameters.

```go
	// It will send request to https://www.google.com/search?q=security&source=hp&sclient=gws-wiz
	if err := gurl.Get("https://www.google.com/search&q=security",
		gurl.WithParam("source", "hp"),
		gurl.WithParam("sclient", "gws-wiz"),
	); err != nil {
		log.Fatal(err.Error())
	}
```

### Custom Headers

`WithHeader()` can add HTTP header to a request.

```go
	if err := gurl.Get(serverURL,
		gurl.WithHeader("Authorization", "Bearer XXXXXXXX"),
	); err != nil {
		log.Fatal(err.Error())
	}
```

### JSON Response Content

`WithReader()` and `DecodeAsJSON()` support to decode JSON formatted response.

```go
	var resp struct {
		KeysURL string `json:"keys_url"`
	}

	if err := gurl.Get("https://api.github.com/",
		gurl.WithReader(gurl.DecodeAsJSON(&resp)),
	); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(resp.KeysURL)
	//Output: https://api.github.com/user/keys
```

### Dump response body

`CopyStream()` calls `io.Copy()` to read response body and to write given `io.Writer`. For example, following code saves result of HTTP request to a temp file.

```go
	f, err := os.CreateTemp("", "*.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := gurl.Get("https://api.github.com/",
		gurl.WithReader(gurl.CopyStream(f)),
	); err != nil {
		log.Fatal(err.Error())
	}
```

### With context.Context

`WithCtx()` will create `http.Request` with `context.Context`.

```go
	if err := gurl.Get("https://slow.example.com",
	    gurl.WithCtx(ctx),
	); err != nil {
		log.Fatal(err.Error())
	}
```

### Error handling

`gurl` uses [goerr](https://github.com/m-mizutani/goerr) and keeps response body content  to error structure when the request is failed. It can be extracted with `errors.As` method.

```go
	if err := gurl.Get("https://example.com",
		gurl.WithClient(&errorClient{}),
	); err != nil {
		var goErr *goerr.Error
		if errors.As(err, &goErr) {
			fmt.Println(string(goErr.Values()["body"].([]byte)))
		}
	}
	//Output: crashed!
```

### Testing with your http.Client or http.HandlerFunc

You can use own http.Client that is implemented `Do(req *http.Request) (*http.Response, error)` method.

```go
type testClient struct{}

func (x *testClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Println(req.URL)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil
}
```

Then, following code will output `https://www.google.com/search?q=security`.

```go
	client := &testClient{}

	if err := gurl.Get("https://www.google.com/search",
		gurl.WithParam("q", "security"),
		gurl.WithClient(client),
	); err != nil {
		log.Fatal(err.Error())
	}
```

Also `WithHandlerFunc()` can be used to test your HTTP server that is implemented `ServeHTTP`.

```go
type testServer struct{}

func (x testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.WriteHeader(200)
	w.Write([]byte("hello"))
}
```

```go
	server := &testServer{}

	if err := gurl.Get("https://www.google.com/search",
		gurl.WithParam("q", "security"),
		gurl.WithHandlerFunc(server.ServeHTTP),
	); err != nil {
		log.Fatal(err.Error())
	}
```

## License  <!-- omit in toc -->

Apache License 2.0
