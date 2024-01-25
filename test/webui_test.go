package gopolar_test

import (
	"fmt"
	"gopolar/internal/core"
	"gopolar/internal/tui"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateManyTunnels(t *testing.T) {
	assert := assert.New(t)
	end := tui.NewCLIEnd()

	serverStart := uint64(10001)
	serverEnd := serverStart - 1 + 60
	tunnelStart := serverEnd + 1
	tunnelEnd := tunnelStart - 1 + 20
	m := make(map[uint64][]uint64) // tunnels, source to a list of dests
	for i := tunnelStart; i <= tunnelEnd; i++ {
		dest := []uint64{}
		n := rand.Intn(5) + 1 // number of dest
		for j := 0; j < n; j++ {
			d := serverStart + uint64(rand.Intn(int(serverEnd-serverStart)))
			assert.True(d >= serverStart && d <= serverEnd)
			dest = append(dest, d)
		}
		m[i] = dest
	}

	// create tunnels
	for s, ds := range m {
		for _, d := range ds {
			t := core.Tunnel{
				Name:   fmt.Sprintf("tfrom %v to %v", s, d),
				Source: fmt.Sprintf("localhost:%v", s),
				Dest:   fmt.Sprintf("localhost:%v", d),
			}
			end.CreateTunnel(t.Name, t.Source, t.Dest)
		}
	}
	list, _ := end.GetTunnelList()
	t.Logf("%v tunnels created\n", len(list))
}
