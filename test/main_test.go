package gopolar_test

import (
	"gopolar/internal/core"
	"os"
	"testing"
)

var tm *core.TunnelManager

func TestMain(t *testing.M) {
	tm = core.NewTunnelManager(false)

	os.Exit(t.Run())
}

// remove all tunnels in tm
func clear() {
	tunnels := tm.GetTunnels()
	for _, t := range tunnels {
		if tm.RemoveTunnel(t.ID) != nil {
			panic("fail to remove tunnel")
		}
	}
}
