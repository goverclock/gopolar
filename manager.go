package gopolar

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
)

type TunnelManager struct {
	tunnels   map[uint64]*Tunnel              // ID -> source
	forwarder map[netip.AddrPort]chan Command // source -> tunnel routine chan
	sock      net.Listener
	router    *gin.Engine
	cfg       *Config

	mu sync.Mutex
}

// init tunnels from config file, exit if any error occurs
func NewTunnelManager() *TunnelManager {
	log.SetPrefix("[core]")
	log.SetFlags(0)
	// TODO:
	// 1. read tunnel list from config file
	// 2. build forward routines for tunnels with command channel, and store the channels

	ret := &TunnelManager{
		tunnels:   make(map[uint64]*Tunnel),
		forwarder: make(map[netip.AddrPort]chan Command),
	}
	ret.setupSock()
	ret.setupRouter()
	cfg := NewConfig()
	ret.cfg = cfg
	for _, t := range cfg.tunnels {
		ret.addTunnel(t)
	}

	return ret
}

func (tm *TunnelManager) Run() {
	tm.router.RunListener(tm.sock)
}

// always return a list sorted by tunnel ID,  never errors
func (tm *TunnelManager) getTunnels() []Tunnel {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	list := make([]Tunnel, 0, len(tm.tunnels))
	for _, t := range tm.tunnels {
		list = append(list, *t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	return list
}

// returns error if tunnel already exists
func (tm *TunnelManager) addTunnel(nt Tunnel) (uint64, error) {
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
	// TODO: optimize by using Command to get forwarder routine's dest list
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
	tm.tunnels[newID] = &nt

	// TODO: create forwarder routine
	return newID, nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) changeTunnel(id uint64, newName string, newSource string, newDest string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tunnels[id]
	if !ok {
		return fmt.Errorf("tunnel %v does not exist", id)
	}
	t.Name = newName
	t.Source = newSource
	t.Dest = newDest
	tm.tunnels[id] = t

	// TODO: update forwarder routine
	return nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) toggleTunnel(id uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tunnels[id]
	if !ok {
		return fmt.Errorf("tunnel %v does not exist", id)
	}
	t.Enable = !t.Enable

	// TODO: update forwarder routine
	return nil
}

// returns error if tunnel with id does not exist
func (tm *TunnelManager) removeTunnel(id uint64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	_, ok := tm.tunnels[id]
	if !ok {
		return fmt.Errorf("tunnel %v does not exist", id)
	}
	delete(tm.tunnels, id)

	// TODO: update forward routine
	return nil
}
