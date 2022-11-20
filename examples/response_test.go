package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/m-mizutani/gurl"
)

func ExampleDecodeAsJSON() {
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

func ExampleCopyStream() {
	f, err := os.CreateTemp("", "*.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := gurl.Get("https://api.github.com/",
		gurl.WithReader(gurl.CopyStream(f)),
	); err != nil {
		log.Fatal(err.Error())
	}

	if err := f.Close(); err != nil {
		log.Fatal(err.Error())
	}
}
