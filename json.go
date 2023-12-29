package gopolar

type CreateTunnelBody struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

type EditTunnelBody struct {
	NewName   string `json:"name"`
	NewSource string `json:"source"`
	NewDest   string `json:"dest"`
}
