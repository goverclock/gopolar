package main

import (
	"context"
	"net"
	"net/http"
)

// connection between gopolar core and cli
type CLIEnd struct {
	client http.Client
}

func NewCLIEnd() *CLIEnd {
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/tmp/gopolar.sock")
			},
		},
	}

	return &CLIEnd{
		client: client,
	}
}

func (ce *CLIEnd) GET(url string) map[string]interface{} {
	response, err := ce.client.Get("http://unix" + url)
	check(err)
	buf := make([]byte, 1024)
	response.Body.Read(buf)

	return nil
}

func (ce *CLIEnd) Quit() {
	// ce.conn.Close()
}
