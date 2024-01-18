package main

import (
	"gopolar/internal/core"
)

func main() {
	tm := core.NewTunnelManager(true)
	tm.Run()
}
