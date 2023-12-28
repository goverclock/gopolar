package gopolar

type Tunnel struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
	Source string `json:"source"`	// e.g. localhost:xxxx
	Dest   string `json:"dest"`		// e.g. 192.168.1.0:7878
}
