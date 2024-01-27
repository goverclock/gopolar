package main

import (
	"github.com/goverclock/gopolar/internal/core"
)

func main() {
	tm := core.NewTunnelManager(true)
	tm.Run()
}
