package gurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/m-mizutani/goerr"
)

type RequestBody func() (io.Reader, string, error)

// EncodeAsJSON will send body as JSON format and add Content-Type header with "application/json" automatically if Content-Type does not exist.
func EncodeAsJSON(data any) RequestBody {
	return func() (io.Reader, string, error) {
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, "", goerr.Wrap(err, "marshaling body as json")
		}
		return bytes.NewReader(raw), "application/json", nil
	}
}

// EncodeAsURL will sends body as URL encoded format and add Content-Type header with "application/x-www-form-urlencoded" automatically if Content-Type does not exist. A value will be converted to string with %v of [fmt.Sprintf].
func EncodeAsURL(data map[string]any) RequestBody {
	return func() (io.Reader, string, error) {
		values := url.Values{}
		for key, value := range data {
			values.Add(key, fmt.Sprintf("%v", value))
		}
		return bytes.NewReader([]byte(values.Encode())), "application/x-www-form-urlencoded", nil
	}
}

func ByReader(r io.Reader) RequestBody {
	return func() (io.Reader, string, error) {
		return r, "", nil
	}
}

type ResponseReader func(r io.ReadCloser) error

func DecodeAsJSON(out any) ResponseReader {
	return func(r io.ReadCloser) error {
		if err := json.NewDecoder(r).Decode(out); err != nil {
			return goerr.Wrap(err, "decoding body as json")
		}
		return nil
	}
}

func CopyStream(w io.Writer) ResponseReader {
	return func(r io.ReadCloser) error {
		if _, err := io.Copy(w, r); err != nil {
			return goerr.Wrap(err, "coping body to writer")
		}
		return nil
	}
}
