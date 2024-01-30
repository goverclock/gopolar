# gopolar

![gopolar](./gopolar.png)

TCP port forwarding tool with both TUI and Web UI support.

> The [gopher image](https://go.dev/blog/gopher) is [Creative Commons Attribution 4.0](https://creativecommons.org/licenses/by/4.0/) licensed, credit to Renee French.

# Install

Simply `go install github.com/goverclock/gopolar/cmd/...@latest` would install gopolar core & TUI, which is enough for basic usage.

The Web UI has to be built with [ install.sh ](./install.sh).

# Usage

After install, run `gpcore` in terminal. In another window run `gptui` to start TUI, which is used to interact with the core.

If the Web UI is installed, visit `localhost:7070`. The Web UI offers exactly same functionality with TUI.

You may want to [ create a system service ](https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6)for gpcore if you are using systemd.

# RESTful API

You can also integrate gopolar easily with its RESTful API.

### Types

```
type Tunnel struct {
    id      uint64
    name    string
    enable  bool
    source  string  // always localhost:xxxx
    dest    string  // e.g. 192.168.10.1:7878
}
```

### Response

```
{
    success: true/false,
    err_msg: "error message",
    data: {...},
}
```

### API

**GET /tunnels/list**

Get all tunnels from core, sorted by tunnel ID.

```
response("data"):
{
    tunnels []Tunnels
}
```

**POST /tunnels/create**

Create a new tunnel.

```
body:
{
    name    string
    source  string
    dest    string
}
response("data"):
{
    id      uint64
}
```

**POST /tunnels/edit/:id**

Edit tunnel with ID.

```
body:
{
    name    string
    source  string
    dest    string
}
```

**POST /tunnels/toggle/:id**

Enable/disable tunnel with ID.

**DELETE /tunnels/delete/:id**

Delete tunnel with ID.

**GET /about**

Information about gopolar.

```
response("data"):
{
    version string  // e.g. 1.0.0
}
```
