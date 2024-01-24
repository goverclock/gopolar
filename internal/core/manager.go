package core

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type TunnelManager struct {
	tunnels   map[uint64]*Tunnel            // ID -> source
	forwarder map[netip.AddrPort]*Forwarder // source -> forwarder, only maintains running tunnels
	router    *gin.Engine
	cfg       *Config

	mu sync.Mutex
}

// init tunnels from config file, exit if any error occurs
func NewTunnelManager(readConfig bool) *TunnelManager {
	log.SetFlags(0)

	ret := &TunnelManager{
		tunnels:   make(map[uint64]*Tunnel),
		forwarder: make(map[netip.AddrPort]*Forwarder),
	}
	ret.setupRouter()

	if readConfig {
		cfg := NewConfig()
		ret.cfg = cfg
		for _, t := range cfg.tunnels {
			ret.AddTunnel(t)
		}
	}

	return ret
}

func (tm *TunnelManager) Run() {
	go tm.router.Run()

	os.Remove("/tmp/gopolar.sock")
	// this creates the unix domain socket
	tm.router.RunUnix("/tmp/gopolar.sock")
}

// tm.mu must be held,
// save current tunnel list to gopolar.toml
func (tm *TunnelManager) saveL() {
	viper.Set("tunnels", tunnelMapToListL(tm.tunnels))
	viper.WriteConfig()
}

// tm.mu must be held,
// create new forwarder if needed,
// then add the forward
func (tm *TunnelManager) addForwardL(src netip.AddrPort, dest string) {
	if tm.forwarder[src] == nil {
		fwd, err := NewForwarder(src)
		if err != nil {
			Debugf("[manager] fail to create new forwarder for src=%v: %v\n", src, err)
			return
		}
		tm.forwarder[src] = fwd
	}
	tm.forwarder[src].Add(dest)
}

// tm.mu must be held,
func (tm *TunnelManager) removeForwardL(src netip.AddrPort, dest string) {
	if tm.forwarder[src].Remove(dest) {
		tm.forwarder[src] = nil
	}
}

// always return a list sorted by tunnel ID,  never errors
func (tm *TunnelManager) GetTunnels() []Tunnel {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tunnelMapToListL(tm.tunnels)
}

// tm.mu must be held
func tunnelMapToListL(m map[uint64]*Tunnel) []Tunnel {
	list := make([]Tunnel, 0, len(m))
	for _, t := range m {
		list = append(list, *t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	return list

}

// new tunnel is enabled by default
// returns error if tunnel already exists
func (tm *TunnelManager) AddTunnel(nt Tunnel) (uint64, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// validate souce, dest
	src, err := nt.ParseSource()
	if err != nil {
		return 0, err
	}
	dest, err := nt.ParseDest()
	if err != nil {
		return 0, err
	}
	if src == dest {
		return 0, fmt.Errorf("source and dest can not be the same: %v", dest)
	}

	// check if a forwarder routine is already running this mapping
	for id, t := range tm.tunnels {
		if t.MustParseSource() == src && t.MustParseDest() == dest {
			return 0, fmt.Errorf("tunnel from %v to %v already exists(ID=%v)", src, dest, id)
		}
	}

	// generate id for the new tunnel
	newID := uint64(1)
	for {
		_, ok := tm.tunnels[newID]
		if ok { // this ID is taken
			newID++
		} else {
			break
		}
	}

	// add the tunnel
	nt.ID = newID
	nt.Enable = true
	tm.tunnels[newID] = &nt

	// update forward
	tm.addForwardL(src, nt.Dest)

	tm.saveL()
	return newID, nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) ChangeTunnel(id uint64, newName string, newSource string, newDest string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tunnels[id]
	if !ok {
		return fmt.Errorf("tunnel %v does not exist", id)
	}

	t.Name = newName
	if t.Source != newSource || t.Dest != newDest {
		// need to update forwarder, by removing old & adding new
		tm.removeForwardL(t.MustParseSource(), t.Dest)
		t.Source = newSource
		t.Dest = newDest
		tm.addForwardL(t.MustParseSource(), newDest)
	}

	tm.saveL()
	return nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) ToggleTunnel(id uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tunnels[id]
	if !ok {
		return fmt.Errorf("tunnel %v does not exist", id)
	}
	t.Enable = !t.Enable

	// update forwarder routine
	if t.Enable {
		tm.addForwardL(t.MustParseSource(), t.Dest)
	} else {
		tm.removeForwardL(t.MustParseSource(), t.Dest)
	}

	tm.saveL()
	return nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) RemoveTunnel(id uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	t, ok := tm.tunnels[id]
	if !ok {
		return fmt.Errorf("tunnel %v does not exist", id)
	}

	// if tunnel is not enabled, it's already not in forwarder
	if t.Enable {
		tm.removeForwardL(t.MustParseSource(), t.Dest)
	}

	delete(tm.tunnels, id)

	tm.saveL()
	return nil
}
