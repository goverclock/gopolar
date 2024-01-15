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
	src, err := net.Listen("tcp", ":"+fmt.Sprint(source.Port()))
	if err != nil {
		return nil, fmt.Errorf("fail to listen localhost:%v", source.Port())
	}
	log.Printf("[forward] new forward listening localhost:%v\n", source.Port())
	fwd := &Forwarder{
		src:      src,
		forwards: make(map[*net.Conn]map[string]*net.Conn),
	}
	go fwd.run()
	return fwd, nil
}

// add a dest(e.g. 198.51.100.1:80) for this forwarder
func (fwd *Forwarder) Add(d string) {
	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	fwd.dest = append(fwd.dest, d)
	log.Printf("[forward] new dest=%v\n", d)

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
			log.Printf("[forward] remove dest=%v\n", d)
			break
		}
	}

	// stop existing connections to this dest
	for cs := range fwd.forwards {
		if fwd.forwards[cs][d] != nil {
			(*fwd.forwards[cs][d]).Close() // this should stops io.Copy
			fwd.forwards[cs][d] = nil
			log.Printf("[forward] removed existing connection: dest=%v for src=%v\n", d, fwd.src.Addr())
		}
	}
	if len(fwd.dest) == 0 {
		fwd.src.Close()
		return true
	}
	return false
}

func (fwd *Forwarder) run() {
	for {
		// for each connect to source
		src := fwd.src
		connS, err := src.Accept() // stop this by src.Close()
		if err != nil {
			log.Printf("[forward] source=%v quiting: %v", src.Addr(), err)
			break // kill the routine by using src.Close()
		}
		log.Printf("[forward] new client connects to source=%v\n", src.Addr())
		log.Printf("[forward] dest=%v\n", fwd.dest)

		fwd.mu.Lock()
		fwd.forwards[&connS] = make(map[string]*net.Conn) // TODO: when is this map deleted?
		for _, d := range fwd.dest {                      // dial all dest for connS
			connD, err := net.Dial("tcp", d)
			log.Printf("[forward] src=%v dialed %v\n", src.Addr(), d)
			if err != nil {
				log.Printf("[forward] fail to dial dest=%v for src=%v\n", d, src.Addr())
				continue
			}
			fwd.forwards[&connS][d] = &connD
			go biCopyIO(connS, connD)
		}
		fwd.mu.Unlock()
	}
}

// TODO: in order to copy to multiple dests, how to read src without consuming?
func biCopyIO(src, dest net.Conn) {
	log.Printf("[forward] start io copy(bidirectional), src=%v, dest=%v\n", src.LocalAddr(), dest.RemoteAddr())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(dest, src)
		dest.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(src, dest)
		src.Close()
	}()
	wg.Wait()
	log.Printf("[forward] end io copy, src=%v, dest=%v\n", src.LocalAddr(), dest.RemoteAddr())
}
