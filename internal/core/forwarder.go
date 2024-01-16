package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"
)

// forward one source to one or multiple dest
type Forwarder struct {
	src         net.Listener
	dest        []string                           // only for new connection to src to set up
	connections map[*net.Conn]map[string]*net.Conn // map[connSrc]map[dest]connDest

	quit bool
	mu   sync.Mutex
}

func NewForwarder(source netip.AddrPort) (*Forwarder, error) {
	src, err := net.Listen("tcp", ":"+fmt.Sprint(source.Port()))
	if err != nil {
		return nil, fmt.Errorf("fail to listen localhost:%v", source.Port())
	}
	log.Printf("[forward] new forward listening localhost:%v\n", source.Port())
	fwd := &Forwarder{
		src:         src,
		connections: make(map[*net.Conn]map[string]*net.Conn),
	}
	go fwd.listen()
	go fwd.copyRoutine()
	return fwd, nil
}

// add a dest(e.g. 198.51.100.1:80) for this forwarder
func (fwd *Forwarder) Add(d string) {
	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	fwd.dest = append(fwd.dest, d)
	log.Printf("[forward] new dest=%v\n", d)

	// dial new dest for all existing connS
	for cs := range fwd.connections {
		connD, err := net.Dial("tcp", d)
		if err != nil {
			log.Printf("[forward] fail to dial dest=%v for src=%v\n", d, fwd.src.Addr())
		}
		fwd.connections[cs][d] = &connD
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
			log.Printf("[forward] removed dest=%v\n", d)
			break
		}
	}

	// stop existing connections to this dest
	for cs := range fwd.connections {
		if fwd.connections[cs][d] != nil {
			(*fwd.connections[cs][d]).Close() // this should stops io.Copy
			delete(fwd.connections[cs], d)
			log.Printf("[forward] removed existing connection: dest=%v for src=%v\n", d, fwd.src.Addr())
		}
	}
	if len(fwd.dest) == 0 {
		fwd.src.Close() // close listener
		// close all existing connS
		for cs := range fwd.connections {
			(*cs).Close()
		}
		fwd.quit = true // notify copyRoutine() and listen() to quit
		return true
	}
	return false
}

func (fwd *Forwarder) listen() {
	for {
		// for each connect to source
		fwd.mu.Lock()
		src := fwd.src
		if fwd.quit {
			log.Printf("[forward] source=%v quitted\n", src.Addr())
			fwd.mu.Unlock()
			break
		}
		fwd.mu.Unlock()

		connS, err := src.Accept() // stop this by src.Close()
		if err != nil {
			log.Printf("[forward] source=%v quitted(error omitted)\n", src.Addr())
			break
		}
		log.Printf("[forward] new client connects to source=%v\n", src.Addr())
		log.Printf("[forward] dest=%v\n", fwd.dest)

		fwd.mu.Lock()
		fwd.connections[&connS] = make(map[string]*net.Conn)
		for _, d := range fwd.dest { // dial all dest for connS
			connD, err := net.Dial("tcp", d)
			if err != nil {
				log.Printf("[forward] fail to dial dest=%v for src=%v\n", d, src.Addr())
				continue
			}
			log.Printf("[forward] src=%v dialed %v\n", src.Addr(), d)
			fwd.connections[&connS][d] = &connD
		}
		fwd.mu.Unlock()
	}
}

func (fwd *Forwarder) copyRoutine() {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		fwd.mu.Lock()

		closedConnS := []*net.Conn{}
		for connS, mde := range fwd.connections {
			// read connS, write to many connD
			(*connS).SetReadDeadline(time.Now().Add(time.Millisecond)) // TODO: doc this in paper, see https://github.com/golang/go/issues/36973
			nr, err := (*connS).Read(buf)
			if nr != 0 {
				for _, connD := range mde {
					(*connD).Write(buf[0:nr])
				}
			}
			if !errors.Is(err, os.ErrDeadlineExceeded) && err != nil { // connS is down
				closedConnS = append(closedConnS, connS)
				continue
			}

			// read many connD, write connS
			totNr := 0
			for _, connD := range mde {
				(*connD).SetReadDeadline(time.Now().Add(time.Millisecond))
				nr, _ := (*connD).Read(buf[totNr:]) // ignore errors from connD
				totNr += nr
			}
			if totNr != 0 {
				(*connS).Write(buf[:totNr])
			}
		}

		// remove closed connS and close relevant connections
		for _, ccs := range closedConnS {
			(*ccs).Close()
			for _, connD := range fwd.connections[ccs] {
				(*connD).Close()
			}
			delete(fwd.connections, ccs)
		}

		if fwd.quit {
			log.Printf("[forward] copyRoutine for src=%v quitted", fwd.src.Addr())
			fwd.mu.Unlock()
			return
		}
		fwd.mu.Unlock()
	}
}