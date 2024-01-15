package gopolar

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"sync"
)

// forward one source to one or multiple dest
type Forwarder struct {
	src      net.Listener
	dest     []string                           // only for new connection to src to set up
	forwards map[*net.Conn]map[string]*net.Conn // map[connSrc]map[dest]connDest

	mu sync.Mutex
}

func NewForwarder(source netip.AddrPort) (*Forwarder, error) {
	src, err := net.Listen("tcp", fmt.Sprint(source.Port()))
	if err != nil {
		return nil, fmt.Errorf("fail to listen localhost:" + fmt.Sprint(source.Port()))
	}
	fwd := &Forwarder{
		src:      src,
		forwards: make(map[*net.Conn]map[string]*net.Conn),
	}
	go fwd.run(src)
	return fwd, nil
}

// add a dest(e.g. 198.51.100.1:80) for this forwarder
func (fwd *Forwarder) Add(d string) {
	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	fwd.dest = append(fwd.dest, d)

	// dial new dest for all existing connS
	for cs := range fwd.forwards {
		connD, err := net.Dial("tcp", d)
		if err != nil {
			log.Printf("[forward] fail to dial dest=%v for src=%v\n", d, fwd.src.Addr())
		}
		fwd.forwards[cs][d] = &connD
		log.Printf("[forward] added dest=%v for src=%v\n", d, (*cs).RemoteAddr())
	}
}

// stop forwarding to a dest, does nothing if not found,
// return true if no dest remains after the operation,
// in which case the forwarder should be deleted
func (fwd *Forwarder) Remove(d string) bool {
	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	for i, v := range fwd.dest {
		if v == d {
			fwd.dest = append(fwd.dest[:i], fwd.dest[i+1:]...)
			break
		}
	}

	// stop existing connections to this dest
	for cs := range fwd.forwards {
		if fwd.forwards[cs][d] != nil {
			(*fwd.forwards[cs][d]).Close() // this should stops io.Copy
			fwd.forwards[cs][d] = nil
			log.Printf("[forward] removed dest=%v for src=%v\n", d, fwd.src.Addr())
		}
	}
	if len(fwd.dest) == 0 {
		fwd.src.Close()
		return true
	}
	return false
}

func (fwd *Forwarder) run(src net.Listener) {
	for {
		// for each connect to source
		connS, err := src.Accept() // stop this by src.Close()
		if err != nil {
			log.Printf("[forward] source=%v fail to accept a connection: %v", src.Addr(), err)
			break // kill the routine by using src.Close()
		}
		log.Printf("[forward] new client connects to source=%v\n", src.Addr())

		fwd.mu.Lock()
		fwd.forwards[&connS] = make(map[string]*net.Conn)
		for _, d := range fwd.dest { // dial all dest for connS
			connD, err := net.Dial("tcp", d)
			if err != nil {
				log.Printf("[forward] fail to dial dest=%v for src=%v\n", d, src.Addr())
			}
			fwd.forwards[&connS][d] = &connD
			go copyIO(connS, connD)
			go copyIO(connD, connS)
		}
		fwd.mu.Unlock()
	}
}

func copyIO(src, dest net.Conn) {
	log.Printf("[forward] start io copy, src=%v, dest=%v\n", src.LocalAddr(), dest.RemoteAddr())
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
	log.Printf("[forward] end io copy, src=%v, dest=%v\n", src.LocalAddr(), dest.RemoteAddr())
}
