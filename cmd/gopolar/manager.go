package main

import (
	"gopolar"
	"math/rand"
	"net/netip"
	"sync"
)

type TunnelManager struct {
	tunnels   []gopolar.Tunnel                        // ID -> source
	forwarder map[netip.AddrPort]chan gopolar.Command // source -> tunnel routine chan

	mu sync.Mutex
}

// init tunnels from config file, exit if any error occurs
func NewTunnelManager( /* cfg *Config */ ) *TunnelManager {
	// TODO:
	// 1. read tunnel list from config file
	// 2. build forward routines for tunnels with command channel, and store the channels
	return &TunnelManager{}
}

func (tm *TunnelManager) GetTunnels() []gopolar.Tunnel {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.tunnels
}

// returns error if tunnel already exists
func (tm *TunnelManager) AddTunnel(t gopolar.Tunnel) (uint64, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// TODO:
	tm.tunnels = append(tm.tunnels, t)
	return rand.Uint64() % 100, nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) ChangeTunnel(id uint64, newName string, newSource string, newDest string) error {
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
func (tm *TunnelManager) ToggleTunnel(id uint64) error {
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
func (tm *TunnelManager) RemoveTunnel(id uint64) error {
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
