package authservice

import (
	"fmt"
	"net/http"
)

type AuthResult struct {
	// Block indicates wether the request should be blocked or not.
	Block bool

	status int
	err    string

	content []byte
	header  http.Header
}

// AddHeaders sets additional headers which sent by the auth server.
func (ar AuthResult) AddHeaders(h http.Header) {
	for k, v := range ar.header {
		h[k] = v
	}
}

func (ar AuthResult) RenderError(wr http.ResponseWriter) error {
	wr.WriteHeader(ar.status)

	for key, v := range ar.header {
		for _, val := range v {
			wr.Header().Add(key, val)
		}
	}

	if len(ar.content) > 0 {
		wr.Write(ar.content)
	} else {
		wr.Write([]byte(ar.err))
	}

	return nil
}

func (ar AuthResult) String() string {
	block := "✅"

	if ar.Block {
		block = "⛔️"
	}

	return fmt.Sprintf("Block: %s Error: %s, StatusCode: %d", block, ar.err, ar.status)
}
