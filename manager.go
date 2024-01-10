package gopolar

import (
	"log"
	"math/rand"
	"net"
	"net/netip"
	"sync"

	"github.com/gin-gonic/gin"
)

type TunnelManager struct {
	tunnels   []Tunnel                        // ID -> source
	forwarder map[netip.AddrPort]chan Command // source -> tunnel routine chan
	sock      net.Listener
	router    *gin.Engine

	mu sync.Mutex
}

// init tunnels from config file, exit if any error occurs
func NewTunnelManager( /* cfg *Config */ ) *TunnelManager {
	log.SetPrefix("[core]")
	log.SetFlags(0)
	// TODO:
	// 1. read tunnel list from config file
	// 2. build forward routines for tunnels with command channel, and store the channels

	ret := &TunnelManager{}
	ret.setupSock()
	ret.setupRouter()
	return ret
}

func (tm *TunnelManager) getTunnels() []Tunnel {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.tunnels
}

func (tm *TunnelManager) Run() {
	tm.router.RunListener(tm.sock)
}

// returns error if tunnel already exists
func (tm *TunnelManager) addTunnel(t Tunnel) (uint64, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// TODO:
	tm.tunnels = append(tm.tunnels, t)
	return rand.Uint64() % 100, nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) changeTunnel(id uint64, newName string, newSource string, newDest string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// TODO:
	for i, t := range tm.tunnels {
		if t.ID != id {
			continue
		}
		tm.tunnels[i].Name = newName
		tm.tunnels[i].Source = newSource
		tm.tunnels[i].Dest = newDest
		break
	}
	return nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) toggleTunnel(id uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// TODO:
	for i, t := range tm.tunnels {
		if t.ID != id {
			continue
		}
		tm.tunnels[i].Enable = !t.Enable
	}

	return nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) removeTunnel(id uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// TODO:
	for i, t := range tm.tunnels {
		if t.ID == id {
			tm.tunnels = append(tm.tunnels[:i], tm.tunnels[i+1:]...)
			break
		}
	}

	return nil
}
