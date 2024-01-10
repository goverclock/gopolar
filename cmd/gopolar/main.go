package main

import (
	"gopolar"
)

func main() {
	// TODO: read config to initialize tunnel manager
	tm := gopolar.NewTunnelManager( /*cfg*/ )
	tm.Run()
}
