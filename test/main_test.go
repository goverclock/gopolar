package gopolar_test

import (
	"os"
	"testing"

	"github.com/goverclock/gopolar/internal/core"
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
