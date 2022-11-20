package gurl

import "github.com/m-mizutani/goerr"

var (
	ErrUnexpectedCode = goerr.New("unexpected status code")
	ErrInvalidOption  = goerr.New("invalid option")
)
