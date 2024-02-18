package core

import (
	"errors"
	"fmt"
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
	connLoggers map[*net.Conn]map[string]*ConnLogger

	quit bool
	mu   sync.Mutex
}

func NewForwarder(source netip.AddrPort) (*Forwarder, error) {
	src, err := net.Listen("tcp", ":"+fmt.Sprint(source.Port()))
	if err != nil {
		return nil, fmt.Errorf("fail to listen localhost:%v", source.Port())
	}
	Debugf("[forward] new forward listening localhost:%v\n", source.Port())
	fwd := &Forwarder{
		src:         src,
		connections: make(map[*net.Conn]map[string]*net.Conn),
		connLoggers: make(map[*net.Conn]map[string]*ConnLogger),
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
	Debugf("[forward] new dest=%v\n", d)

	// dial new dest for all existing connS
	for cs := range fwd.connections {
		connD, err := net.Dial("tcp", d)
		if err != nil {
			Debugf("[forward] fail to dial dest=%v for src=%v, err=%v\n", d, fwd.src.Addr(), err)
		}
		fwd.connections[cs][d] = &connD
		fwd.connLoggers[cs][d] = NewConnLogger(fwd.src.Addr().String(), d)
		Debugf("[forward] added dest=%v for src=%v\n", d, (*cs).RemoteAddr())
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
			Debugf("[forward] removed dest=%v\n", d)
			break
		}
	}

	// stop existing connections to this dest
	for cs := range fwd.connections {
		if fwd.connections[cs][d] != nil {
			(*fwd.connections[cs][d]).Close() // this should stops io.Copy
			delete(fwd.connections[cs], d)
			delete(fwd.connLoggers[cs], d)
			Debugf("[forward] ended existing connection: dest=%v for src=%v\n", d, fwd.src.Addr())
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
			Debugf("[forward] source=%v quitted\n", src.Addr())
			fwd.mu.Unlock()
			break
		}
		fwd.mu.Unlock()

		connS, err := src.Accept() // stop this by src.Close()
		if err != nil {
			Debugf("[forward] source=%v quitted(error omitted)\n", src.Addr())
			break
		}
		Debugf("[forward] new client connects to source=%v\n", src.Addr())
		Debugf("[forward] dest=%v\n", fwd.dest)

		fwd.mu.Lock()
		fwd.connections[&connS] = make(map[string]*net.Conn)
		fwd.connLoggers[&connS] = make(map[string]*ConnLogger)
		established := false
		for _, d := range fwd.dest { // dial all dest for connS
			connD, err := net.Dial("tcp", d)
			if err != nil {
				Debugf("[forward] fail to dial dest=%v for src=%v, err=%v\n", d, src.Addr(), err)
				continue
			}
			Debugf("[forward] src=%v dialed %v\n", src.Addr(), d)
			fwd.connections[&connS][d] = &connD
			fwd.connLoggers[&connS][d] = NewConnLogger(fwd.src.Addr().String(), d)
			established = true
		}
		if !established {
			connS.Close()
		}
		fwd.mu.Unlock()
	}
}

func (fwd *Forwarder) copyRoutine() {
	size := 1024 * 1024 // TODO: doc this, or this should be adjustable
	buf := make([]byte, size)
	for {
		fwd.mu.Lock()

		closedConnS := []*net.Conn{}
		for connS, mde := range fwd.connections {
			// read connS, write to many connD
			(*connS).SetReadDeadline(time.Now().Add(time.Microsecond)) // TODO: doc this in paper, see https://github.com/golang/go/issues/36973
			nr, err := (*connS).Read(buf)
			if nr != 0 {
				for d, connD := range mde {
					(*connD).Write(buf[0:nr])
					fwd.connLoggers[connS][d].LogSend(buf[0:nr])
				}
			}
			if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) { // connS is down
				closedConnS = append(closedConnS, connS)
				continue
			}

			// read many connD, write connS
			totNr := 0
			for d, connD := range mde {
				(*connD).SetReadDeadline(time.Now().Add(time.Microsecond))
				nr, _ := (*connD).Read(buf[totNr:]) // ignore errors from connD
				if nr != 0 {
					fwd.connLoggers[connS][d].LogRecv(buf[totNr : totNr+nr])
				}
				// Debugf("[forward] read %v bytes from dest=%v, err=%v", nr, dest, err)
				totNr += nr
			}
			if totNr != 0 {
				_, _ = (*connS).Write(buf[:totNr])
				// Debugf("[forward] write %v bytes to src=%v, err=%v", nw, fwd.src.Addr(), err)
			}
		}

		// remove closed connS and close relevant connections
		for _, ccs := range closedConnS {
			(*ccs).Close()
			Debugf("[forward] connS closed for src=%v", fwd.src.Addr())
			for _, connD := range fwd.connections[ccs] {
				(*connD).Close()
			}
			delete(fwd.connections, ccs)
			delete(fwd.connLoggers, ccs)
		}

		if fwd.quit {
			Debugf("[forward] copyRoutine for src=%v quitted", fwd.src.Addr())
			fwd.mu.Unlock()
			return
		}
		fwd.mu.Unlock()
	}
}
