package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/goverclock/gopolar/internal/core"

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
		Timeout: 3 * time.Second,
	}
	return &CLIEnd{
		client: client,
	}
}

func (ce *CLIEnd) GetTunnelList() ([]core.Tunnel, error) {
	response, err := ce.GET("/tunnels/list")
	if err != nil {
		return nil, err
	}
	ret := []core.Tunnel{}
	mapstructure.Decode(response["tunnels"], &ret)
	return ret, nil
}

// returns ID of the new tunnel
func (ce *CLIEnd) CreateTunnel(name string, source string, dest string) (uint64, error) {
	body := core.CreateTunnelBody{
		Name:   name,
		Source: source,
		Dest:   dest,
	}
	response, err := ce.POST("/tunnels/create", body)
	if err != nil {
		return 0, err
	}
	id := uint64(response["id"].(float64))
	return id, nil
}

func (ce *CLIEnd) EditTunnel(id uint64, newName string, newSource string, newDest string) error {
	body := core.EditTunnelBody{
		NewName:   newName,
		NewSource: newSource,
		NewDest:   newDest,
	}
	_, err := ce.POST("/tunnels/edit/"+fmt.Sprint(id), body)
	return err
}

func (ce *CLIEnd) ToggleTunnel(id int64) error {
	_, err := ce.POST("/tunnels/toggle/"+strconv.FormatInt(id, 10), nil)
	return err
}

func (ce *CLIEnd) DeleteTunnel(id int64) error {
	_, err := ce.DELETE("/tunnels/delete/" + strconv.FormatInt(id, 10))
	return err
}

func (ce *CLIEnd) GetAboutInfo() (core.AboutInfo, error) {
	ret := core.AboutInfo{}
	response, err := ce.GET("/about")
	if err != nil {
		return ret, err
	}
	mapstructure.Decode(response["about"], &ret)
	return ret, nil
}

func bodyToJSON(body io.ReadCloser) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	jsonBytes, err := io.ReadAll(body)
	// WriteTTY("/dev/ttys018", fmt.Sprintln(string(jsonBytes)))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (ce *CLIEnd) GET(url string) (map[string]interface{}, error) {
	response, err := ce.client.Get("http://unix" + url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %v responses code %v", url, response.StatusCode)
	}

	ret, err := bodyToJSON(response.Body)
	if err != nil {
		return nil, err
	}
	if !ret["success"].(bool) {
		return nil, fmt.Errorf("Operation failed: " + ret["err_msg"].(string))
	}
	ret = ret["data"].(map[string]interface{})
	return ret, nil
}

func (ce *CLIEnd) POST(url string, data interface{}) (map[string]interface{}, error) {
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

	ret, err := bodyToJSON(response.Body)
	if err != nil {
		return nil, err
	}
	if !ret["success"].(bool) {
		return nil, fmt.Errorf("Operation failed: " + ret["err_msg"].(string))
	}
	ret = ret["data"].(map[string]interface{})
	return ret, nil
}

func (ce *CLIEnd) DELETE(url string) (map[string]interface{}, error) {
	// mimic of http.Client.Get
	response, err := func(c *http.Client, url string) (resp *http.Response, err error) {
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return nil, err
		}
		return c.Do(req)
	}(&ce.client, "http://unix"+url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DELETE %v responses code %v", url, response.StatusCode)
	}

	ret, err := bodyToJSON(response.Body)
	if err != nil {
		return nil, err
	}
	if !ret["success"].(bool) {
		return nil, fmt.Errorf("Operation failed: " + ret["err_msg"].(string))
	}
	ret = ret["data"].(map[string]interface{})
	return ret, nil
}
