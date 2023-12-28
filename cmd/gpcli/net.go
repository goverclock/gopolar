package main

import (
	"context"
	"encoding/json"
	"io"
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
	ret := make(map[string]interface{})
	response, err := ce.client.Get("http://unix" + url)
	check(err)
	jsonData, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(jsonData, &ret)
	check(err)
	ret = ret["data"].(map[string]interface{})
	return ret
}
