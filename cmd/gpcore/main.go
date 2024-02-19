package main

import (
	"flag"

	"github.com/goverclock/gopolar/internal/core"
)

func main() {
	logPtr := flag.Bool("log", false, "enable logging for debugging, disable for better performance")
	nosavePtr := flag.Bool("nosave", false, "ignore saved tunnels in ~/.config/gopolar/tunnels.toml")
	flag.Parse()

	cfg := core.DefaultConfig
	cfg.DoLogs = *logPtr
	cfg.ReadSaved = !*nosavePtr
	tm := core.NewTunnelManager(cfg)
	tm.Run()
}
