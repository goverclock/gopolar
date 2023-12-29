package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopolar"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

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

func (ce *CLIEnd) GetTunnelsList() ([]gopolar.Tunnel, error) {
	response, err := ce.GET("/tunnels/list")
	if err != nil {
		return nil, err
	}
	tunnelsList := response["tunnels"].([]interface{})
	ret := []gopolar.Tunnel{}
	mapstructure.Decode(tunnelsList, &ret)
	return ret, nil
}

func (ce *CLIEnd) CreateTunnel(name string, source string, dest string) error {
	body := gopolar.CreateTunnelBody{
		Name:   name,
		Source: source,
		Dest:   dest,
	}
	_, err := ce.POST("/tunnels/create", body)
	return err
}

func (ce *CLIEnd) EditTunnel(id int64, newName string, newSource string, newDest string) error {
	body := gopolar.EditTunnelBody{
		NewName:   newName,
		NewSource: newSource,
		NewDest:   newDest,
	}
	_, err := ce.POST("/tunnels/edit/"+strconv.FormatInt(id, 10), body)
	return err
}

// TODO
func (ce *CLIEnd) DeleteTunnel(id int64) error {
	return nil
}

func (ce *CLIEnd) GET(url string) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	response, err := ce.client.Get("http://unix" + url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %v responses code %v", url, response.StatusCode)
	}

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &ret)
	if err != nil {
		return nil, err
	}
	if !ret["success"].(bool) {
		log.Println("fail to GET", ret["err_msg"])
	}
	ret = ret["data"].(map[string]interface{})
	return ret, nil
}

func (ce *CLIEnd) POST(url string, data interface{}) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	response, err := ce.client.Post("http://unix"+url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST %v responses code %v", url, response.StatusCode)
	}

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &ret)
	if err != nil {
		return nil, err
	}
	if !ret["success"].(bool) {
		log.Println("fail to POST", ret["err_msg"])
	}
	ret = ret["data"].(map[string]interface{})
	return ret, nil
}
