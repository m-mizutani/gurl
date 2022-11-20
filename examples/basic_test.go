package examples

import (
	"bytes"
	"fmt"
	"log"

	"github.com/m-mizutani/gurl"
)

const (
	serverURL = "https://emhkq5vqrco2fpr6zqlctbjale0eyygt.lambda-url.ap-northeast-1.on.aws"
)

func ExampleGet() {
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
}

func ExampleWithParam() {
	// It will send request to https://www.google.com/search?q=security
	if err := gurl.Get("https://www.google.com/search",
		gurl.WithParam("q", "security"),
	); err != nil {
		log.Fatal(err.Error())
	}

	// It will send request to https://www.google.com/search?q=security&source=hp&sclient=gws-wiz
	if err := gurl.Get("https://www.google.com/search&q=security",
		gurl.WithParam("source", "hp"),
		gurl.WithParam("sclient", "gws-wiz"),
	); err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleWithHeader() {
	if err := gurl.Get(serverURL,
		gurl.WithHeader("Authorization", "Bearer XXXXXXXX"),
	); err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleByReader() {
	r := bytes.NewReader([]byte("my_test_data"))

	if err := gurl.Post(serverURL,
		gurl.WithBody(gurl.ByReader(r)),
	); err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEncodeAsJSON() {
	var req struct {
		Name string `json:"name"`
	}

	if err := gurl.Post(serverURL,
		gurl.WithBody(gurl.EncodeAsJSON(req)),
	); err != nil {
		log.Fatal(err.Error())
	}
}

func ExampleEncodeAsURL() {
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
}
