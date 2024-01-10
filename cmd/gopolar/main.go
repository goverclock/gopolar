package main

import (
	"log"
)

var tunnelManager *TunnelManager

func main() {
	log.SetPrefix("[core]")
	log.SetFlags(0)

	// TODO: read config to initialize tunnel manager
	// cfg := NewConfig("config.toml")
	tunnelManager = NewTunnelManager( /* cfg */ )
	sock := setupSock()
	router := setupRouter()
	router.RunListener(sock)
}
