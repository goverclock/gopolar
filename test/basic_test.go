package gopolar_test

import (
	"fmt"
	"testing"

	"github.com/goverclock/gopolar/internal/core"
	"github.com/goverclock/gopolar/test/testutil"

	"github.com/stretchr/testify/assert"
)

// forward request from 3300 to 8800
func TestOne2One(t *testing.T) {
	assert := assert.New(t)
	clear()

	_, err := tm.AddTunnel(core.Tunnel{
		Name:   "tfrom 3300 to 8800",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:8800",
	})
	assert.Nil(err)

	prefix := "hello"
	serv8800 := testutil.NewEchoServer(8800, prefix)
	defer serv8800.Quit()

	clnt3300 := testutil.NewEchoClient(3300)
	err = clnt3300.Connect()
	assert.Nil(err)

	msg := "asdfg\n"
	err = clnt3300.Send(msg)
	assert.Nil(err)
	reply := clnt3300.Recv()
	assert.Equal(prefix+msg, reply)
}

// forward request from 3300 to 8800, 8801, 8802, ..., 8900
func TestOne2Many(t *testing.T) {
	assert := assert.New(t)
	clear()

	msg := "what good lol\n"
	expectRecv := 0
	start := 8800
	end := start + 100
	for i := uint64(start); i <= uint64(end); i++ {
		// tunnels
		p := fmt.Sprint(i)
		tn := core.Tunnel{
			Name:   "tfrom 3300 to " + p,
			Enable: true,
			Source: "localhost:3300",
			Dest:   "localhost:" + p,
		}
		_, err := tm.AddTunnel(tn)
		assert.Nil(err)

		// servers
		prefix := "serv" + p
		s := testutil.NewEchoServer(i, prefix)
		defer s.Quit()
		expectRecv += len(prefix) + len(msg)
	}

	clnt3300 := testutil.NewEchoClient(3300)
	err := clnt3300.Connect()
	assert.Nil(err)

	err = clnt3300.Send(msg)
	assert.Nil(err)

	// should read multiple times, since each Recv() returns on '\n'
	reply := ""
	for i := 0; i < end-start+1; i++ {
		reply += clnt3300.Recv()
	}
	assert.Equal(expectRecv, len(reply))
}

// forward request from 8800, 89, 90, ..., 100 to 3300
func TestMany2One(t *testing.T) {
	assert := assert.New(t)
	clear()

	prefix := "heyqasdjklqjweoi"
	serv3300 := testutil.NewEchoServer(3300, prefix)
	defer serv3300.Quit()

	start := 8800
	end := 100
	for i := uint64(start); i <= uint64(end); i++ {
		// tunnels
		p := fmt.Sprint(i)
		tn := core.Tunnel{
			Name:   fmt.Sprintf("tfrom %v to 3300", i),
			Enable: true,
			Source: "localhost:" + p,
			Dest:   "localhost:3300",
		}
		_, err := tm.AddTunnel(tn)
		assert.Nil(err)

		// clients
		msg := fmt.Sprintf("hahaha message from client %v\n", i)
		c := testutil.NewEchoClient(i)
		assert.Nil(c.Connect())
		assert.Nil(c.Send(msg))
		assert.Equal(prefix+msg, c.Recv())
	}
}

// client should fail to connect when no server is running
// TODO(pending): forwarder need a way to detect TCP dial while not accepting it,
// currently forwarder always accepts client dial, then close it if no,
// dest can be reached
func TestNoServer(t *testing.T) {
	assert := assert.New(t)
	clear()

	_, err := tm.AddTunnel(core.Tunnel{
		Name:   "dummy",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:8800",
	})
	assert.Nil(err)

	clnt := testutil.NewEchoClient(3300)
	// assert.NotEqual(nil, clnt.Connect())
	clnt.Connect()
}

// connection should be closed when tunnel is deleted/disabled
func TestDisconnect(t *testing.T) {
	assert := assert.New(t)
	clear()

	id, err := tm.AddTunnel(core.Tunnel{
		Name:   "3300to8800",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:8800",
	})
	assert.Nil(err)

	prefix := "this is a server running on port 8800"
	serv8800 := testutil.NewEchoServer(8800, prefix)
	defer serv8800.Quit()

	clnt3300 := testutil.NewEchoClient(3300)
	assert.False(clnt3300.IsConnected())
	clnt3300.Connect()
	msg := "this is a message from clnt 3300\n"
	assert.Nil(clnt3300.Send(msg))
	assert.Equal(prefix+msg, clnt3300.Recv())
	assert.True(clnt3300.IsConnected())

	// remove the tunnel
	assert.Nil(tm.RemoveTunnel(id))

	assert.False(clnt3300.IsConnected())
}

// TODO(pending): TestSilentServer, server always close connections immediately
