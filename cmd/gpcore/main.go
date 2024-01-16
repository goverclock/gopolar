package main

import (
	"gopolar/internal/core"
)

func main() {
	tm := core.NewTunnelManager()
	tm.Run()
}
