package gopolar_test

import (
	"fmt"
	"gopolar/internal/core"
	"gopolar/test/testutil"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -bench=. -benchmem -v

// returns n MB bytes without zero
func MakeDataMB(n uint) []byte {
	size := n * 1024 * 1024
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte('a' + (i % 26))
	}
	return buf
}

func RandomString(maxLen int) string {
	buf := make([]byte, rand.Intn(maxLen)+1)
	for i := range buf {
		buf[i] = byte('a' + (rand.Intn(26)))
	}
	return string(buf)
}

// connect, send, recv, verify, disconnect
// only works for One2One or Many2One
func transmit(msg *string, serv *testutil.EchoServer, clnt *testutil.EchoClient) {
	clnt.Connect()
	err := clnt.Send(*msg)
	if err != nil {
		panic(err)
	}
	reply := clnt.Recv()
	if reply != serv.Prefix+*msg {
		panic("Direct(): got inconsistent data")
	}
	clnt.Disconnect()
}

// as standard
func BenchmarkSingleDirect100MB(b *testing.B) {
	clear()

	serv88 := testutil.NewEchoServer(88, "hello")
	defer serv88.Quit()
	clnt88 := testutil.NewEchoClient(88)

	data := string(MakeDataMB(100))
	msg := fmt.Sprintf("data from client: %v\n", data)
	for i := 0; i < b.N; i++ {
		transmit(&msg, serv88, clnt88)
	}
}

// as standard
func BenchmarkSingleDirect500MB(b *testing.B) {
	clear()

	serv88 := testutil.NewEchoServer(88, "hello")
	defer serv88.Quit()
	clnt88 := testutil.NewEchoClient(88)

	data := string(MakeDataMB(500))
	msg := fmt.Sprintf("data from client: %v\n", data)
	for i := 0; i < b.N; i++ {
		transmit(&msg, serv88, clnt88)
	}
}

// same with TestOne2One, but with 100 MB data
func BenchmarkSingleForward100MB(b *testing.B) {
	assert := assert.New(b)
	clear()

	_, err := tm.AddTunnel(core.Tunnel{
		Name:   "tfrom 33 to 88",
		Source: "localhost:33",
		Dest:   "localhost:88",
	})
	assert.Nil(err)

	serv88 := testutil.NewEchoServer(88, "hahaha")
	defer serv88.Quit()
	clnt33 := testutil.NewEchoClient(33)

	data := string(MakeDataMB(100))
	msg := fmt.Sprintf("data from client: %v\n", data)
	for i := 0; i < b.N; i++ {
		transmit(&msg, serv88, clnt33)
	}
}

func BenchmarkSingleForward500MB(b *testing.B) {
	assert := assert.New(b)
	clear()

	_, err := tm.AddTunnel(core.Tunnel{
		Name:   "tfrom 33 to 88",
		Source: "localhost:33",
		Dest:   "localhost:88",
	})
	assert.Nil(err)

	serv88 := testutil.NewEchoServer(88, "hahaha")
	defer serv88.Quit()
	clnt33 := testutil.NewEchoClient(33)

	data := string(MakeDataMB(500))
	msg := fmt.Sprintf("data from client: %v\n", data)
	for i := 0; i < b.N; i++ {
		transmit(&msg, serv88, clnt33)
	}
}

// 600 servers listening on [10001,10600],
// forward [10601,10800] to 1-5 dest in [10001,10600],
// 600 clients connects to [10601,10800],
// echo client sends 1 MB data, validate response
func BenchmarkManyConnectionsForward1MB(b *testing.B) {
	assert := assert.New(b)
	clear()

	// setup servers
	servs := []*testutil.EchoServer{}
	serverStart := uint64(10001)
	serverEnd := serverStart - 1 + 600
	for i := serverStart; i <= serverEnd; i++ {
		s := testutil.NewEchoServer(i, RandomString(20))
		servs = append(servs, s)
		defer s.Quit()
	}
	b.Logf("%v servers created", len(servs))

	tunnelStart := serverEnd + 1
	tunnelEnd := tunnelStart - 1 + 200
	m := make(map[uint64][]uint64)      // tunnels, source to a list of dests
	serverCount := make(map[uint64]int) // number of dest for each source
	for i := tunnelStart; i <= tunnelEnd; i++ {
		dest := []uint64{}
		n := rand.Intn(5) + 1 // number of dest
		for j := 0; j < n; j++ {
			d := serverStart + uint64(rand.Intn(int(serverEnd-serverStart)))
			assert.True(d >= serverStart && d <= serverEnd)
			dest = append(dest, d)
		}
		serverCount[i] = n
		m[i] = dest
	}

	// create tunnels
	for s, ds := range m {
		for _, d := range ds {
			_, err := tm.AddTunnel(core.Tunnel{
				Name:   fmt.Sprintf("tfrom %v to %v", s, d),
				Source: fmt.Sprintf("localhost:%v", s),
				Dest:   fmt.Sprintf("localhost:%v", d),
			})
			if err != nil {
				serverCount[s]--
			}
		}
	}
	list := tm.GetTunnels()
	b.Logf("%v tunnels created\n", len(list))

	// create all clients, then send data on each, recv all response
	clnts := []*testutil.EchoClient{}
	data := string(MakeDataMB(1))
	msg := fmt.Sprintf("data from client: %v\n", data)
	var wg sync.WaitGroup
	for i := tunnelStart; i <= tunnelEnd; i++ {
		wg.Add(1)
		c := testutil.NewEchoClient(i)
		clnts = append(clnts, c)
		go func() {
			assert.Nil(c.Connect())
			assert.Nil(c.Send(msg))
			for j := 0; j < serverCount[c.Port]; j++ {
				c.Recv()
			}
			wg.Done()
		}()
	}
	b.Logf("%v clients created", len(clnts))
	wg.Wait()

	totClntRecv := uint64(0)
	for _, c := range clnts {
		totClntRecv += c.TotRecv
	}
	totServReply := uint64(0)
	for _, s := range servs {
		totServReply += s.TotEcho
	}

	b.Logf("total client recv=%v, total server reply=%v", totClntRecv, totServReply)
	assert.Equal(totClntRecv, totServReply)
}
