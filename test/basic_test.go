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
// func TestNoServer(t *testing.T) {
// 	assert := assert.New(t)
// 	clear()

// 	_, err := tm.AddTunnel(core.Tunnel{
// 		Name:   "dummy",
// 		Enable: true,
// 		Source: "localhost:3300",
// 		Dest:   "localhost:8800",
// 	})
// 	assert.Nil(err)

// 	clnt := testutil.NewEchoClient(3300)
// 	assert.NotEqual(nil, clnt.Connect())
// 	// clnt.Connect()
// }

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

// if server closes connection, client should note it
func TestSilentServer(t *testing.T) {
	assert := assert.New(t)
	clear()

	_, err := tm.AddTunnel(core.Tunnel{
		Name:   "3300to8800silent",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:8800",
	})
	assert.Nil(err)

	sserv8800 := testutil.NewSilentServer(8800)
	defer sserv8800.Quit()

	clnt3300 := testutil.NewEchoClient(3300)
	assert.Nil(clnt3300.Connect())
	// the server disconnects immediately
	assert.False(clnt3300.IsConnected())
}

// TunnelManager should return error when creating a tunnel
// with (source, dest) same with existing tunnel
func TestDenyDuplicateTunnel(t *testing.T) {
	assert := assert.New(t)
	clear()

	tn := core.Tunnel{
		Name:   "3300to8800",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:8800",
	}
	_, err := tm.AddTunnel(tn)
	assert.Nil(err)

	tn.Name = "haha"
	_, err = tm.AddTunnel(tn)
	assert.NotNil(err)
}

// should dial a new dest for existing connection when a new tunnel with
// existing source is created
func TestCreateTunnelOnline(t *testing.T) {
	assert := assert.New(t)
	clear()

	tn1 := core.Tunnel{
		Name:   "3300to8800",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:8800",
	}
	_, err := tm.AddTunnel(tn1)
	assert.Nil(err)

	prefix1 := "hello"
	serv8800 := testutil.NewEchoServer(8800, prefix1)
	defer serv8800.Quit()
	prefix2 := "bye"
	serv9900 := testutil.NewEchoServer(9900, prefix2)
	defer serv9900.Quit()

	clnt3300 := testutil.NewEchoClient(3300)
	err = clnt3300.Connect()
	assert.Nil(err)

	// one2one validate
	msg := "asdfg\n"
	err = clnt3300.Send(msg)
	assert.Nil(err)
	reply := clnt3300.Recv()
	assert.Equal(prefix1+msg, reply)

	// new dest
	tn2 := core.Tunnel{
		Name:   "3300to9900",
		Enable: true,
		Source: "localhost:3300",
		Dest:   "localhost:9900",
	}
	_, err = tm.AddTunnel(tn2)
	assert.Nil(err)

	// one2many validate
	clnt3300.TotRecv = 0
	serv8800.TotEcho = 0

	msg = "ashortmessage\n"
	err = clnt3300.Send(msg)
	assert.Nil(err)
	clnt3300.Recv()
	clnt3300.Recv()

	// t.Logf("clnt3300-%v serv8800-%v serv9900-%v\n", clnt3300.TotRecv, serv8800.TotEcho, serv9900.TotEcho)
	assert.Equal(clnt3300.TotRecv, serv8800.TotEcho+serv9900.TotEcho)
}

// TODO(pending): TestEditDuplicated