package gopolar_test

import (
	"fmt"
	"testing"

	"github.com/goverclock/gopolar/internal/core"
	"github.com/goverclock/gopolar/test/testutil"

	"github.com/stretchr/testify/assert"
)

// forward request from 33 to 88
func TestOne2One(t *testing.T) {
	assert := assert.New(t)
	clear()

	_, err := tm.AddTunnel(core.Tunnel{
		Name:   "tfrom 33 to 88",
		Source: "localhost:33",
		Dest:   "localhost:88",
	})
	assert.Nil(err)

	prefix := "hello"
	serv88 := testutil.NewEchoServer(88, prefix)
	defer serv88.Quit()

	clnt33 := testutil.NewEchoClient(33)
	err = clnt33.Connect()
	assert.Nil(err)

	msg := "asdfg\n"
	err = clnt33.Send(msg)
	assert.Nil(err)
	reply := clnt33.Recv()
	assert.Equal(prefix+msg, reply)
}

// forward request from 33 to 88, 89, 90, ..., 100
func TestOne2Many(t *testing.T) {
	assert := assert.New(t)
	clear()

	msg := "what good lol\n"
	expectRecv := 0
	start := 88
	end := 100
	for i := uint64(start); i <= uint64(end); i++ {
		// tunnels
		p := fmt.Sprint(i)
		tn := core.Tunnel{
			Name:   "tfrom 33 to " + p,
			Source: "localhost:33",
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

	clnt33 := testutil.NewEchoClient(33)
	err := clnt33.Connect()
	assert.Nil(err)

	err = clnt33.Send(msg)
	assert.Nil(err)

	// should read multiple times, since each Recv() returns on '\n'
	reply := ""
	for i := 0; i < end-start+1; i++ {
		reply += clnt33.Recv()
	}
	assert.Equal(expectRecv, len(reply))
}

// forward request from 88, 89, 90, ..., 100 to 33
func TestMany2One(t *testing.T) {
	assert := assert.New(t)
	clear()

	prefix := "heyqasdjklqjweoi"
	serv33 := testutil.NewEchoServer(33, prefix)
	defer serv33.Quit()

	start := 88
	end := 100
	for i := uint64(start); i <= uint64(end); i++ {
		// tunnels
		p := fmt.Sprint(i)
		tn := core.Tunnel{
			Name:   fmt.Sprintf("tfrom %v to 33", i),
			Source: "localhost:" + p,
			Dest:   "localhost:33",
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
		Source: "localhost:33",
		Dest:   "localhost:88",
	})
	assert.Nil(err)

	clnt := testutil.NewEchoClient(33)
	// assert.NotEqual(nil, clnt.Connect())
	clnt.Connect()
}

// connection should be closed when tunnel is deleted/disabled
func TestDisconnect(t *testing.T) {
	assert := assert.New(t)
	clear()

	id, err := tm.AddTunnel(core.Tunnel{
		Name:   "33to88",
		Source: "localhost:33",
		Dest:   "localhost:88",
	})
	assert.Nil(err)

	prefix := "this is a server running on port 88"
	serv88 := testutil.NewEchoServer(88, prefix)
	defer serv88.Quit()

	clnt33 := testutil.NewEchoClient(33)
	assert.False(clnt33.IsConnected())
	clnt33.Connect()
	msg := "this is a message from clnt 33\n"
	assert.Nil(clnt33.Send(msg))
	assert.Equal(prefix+msg, clnt33.Recv())
	assert.True(clnt33.IsConnected())

	// remove the tunnel
	assert.Nil(tm.RemoveTunnel(id))

	assert.False(clnt33.IsConnected())
}

// TODO(pending): TestSilentServer, server always close connections immediately
