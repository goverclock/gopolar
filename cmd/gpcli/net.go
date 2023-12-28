package main

import (
	"context"
	"encoding/json"
	"gopolar"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/mitchellh/mapstructure"
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

// TODO
func (ce *CLIEnd) GetTunnelsList() ([]gopolar.Tunnel, error) {
	data := ce.GET("/tunnels/list")
	tunnelsList := data["tunnels"].([]interface{})
	ret := []gopolar.Tunnel{}
	mapstructure.Decode(tunnelsList, &ret)
	return ret, nil
}

// TODO
func (ce *CLIEnd) CreateTunnel(name string, source string, dest string) error {
	return nil
}

// TODO
func (ce *CLIEnd) DeleteTunnel(id int64) error {
	return nil
}

func (ce *CLIEnd) GET(url string) map[string]interface{} {
	ret := make(map[string]interface{})
	response, err := ce.client.Get("http://unix" + url)
	check(err)
	if response.StatusCode != http.StatusOK {
		log.Printf("fail to GET %v", url)
		return nil
	}

	jsonData, err := io.ReadAll(response.Body)
	check(err)
	err = json.Unmarshal(jsonData, &ret)
	check(err)
	if !ret["success"].(bool) {
		log.Println("fail to GET", ret["err_msg"])
	}
	ret = ret["data"].(map[string]interface{})
	return ret
}
