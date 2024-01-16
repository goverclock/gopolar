package core

import (
	"fmt"
	"net/netip"
	"strings"
)

type Tunnel struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
	Source string `json:"source"` // e.g. localhost:xxxx
	Dest   string `json:"dest"`   // e.g. 192.168.1.0:7878, localhost:7878
}

func (t Tunnel) String() string {
	ret := fmt.Sprintf("Tunnel %v:\n", t.ID)
	ret += fmt.Sprintf("\tName: %v\n", t.Name)
	ret += fmt.Sprintf("\tEnable: %v\n", t.Enable)
	ret += fmt.Sprintf("\tSource: %v\n", t.Source)
	ret += fmt.Sprintf("\tDest: %v\n", t.Dest)
	return ret
}

func (t Tunnel) ParseSource() (netip.AddrPort, error) {
	s := strings.ReplaceAll(t.Source, "localhost", "127.0.0.1")
	return netip.ParseAddrPort(s)
}

func (t Tunnel) MustParseSource() netip.AddrPort {
	s := strings.ReplaceAll(t.Source, "localhost", "127.0.0.1")
	return netip.MustParseAddrPort(s)
}

func (t Tunnel) ParseDest() (netip.AddrPort, error) {
	s := strings.ReplaceAll(t.Dest, "localhost", "127.0.0.1")
	return netip.ParseAddrPort(s)
}

func (t Tunnel) MustParseDest() netip.AddrPort {
	s := strings.ReplaceAll(t.Dest, "localhost", "127.0.0.1")
	return netip.MustParseAddrPort(s)
}
