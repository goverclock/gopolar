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