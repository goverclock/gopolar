package gopolar_test

import (
	"fmt"
	"gopolar/internal/core"
	"gopolar/test/testutil"
	"testing"

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
	assert.Equal(nil, err)

	prefix := "hello"
	serv88 := testutil.NewEchoServer(t, 88, prefix)
	defer serv88.Quit()

	clnt33 := testutil.NewEchoClient(t, 33)
	err = clnt33.Connect()
	assert.Equal(nil, err)

	msg := "asdfg\n"
	err = clnt33.Send(msg)
	assert.Equal(nil, err)
	reply := clnt33.Recv()
	assert.Equal(prefix+msg, reply)
}

// forward request from 33 to 88, 89, 90, ... ,100
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
		assert.Equal(nil, err)

		// servers
		prefix := "serv" + p
		s := testutil.NewEchoServer(t, i, prefix)
		defer s.Quit()
		expectRecv += len(prefix) + len(msg)
	}

	clnt33 := testutil.NewEchoClient(t, 33)
	err := clnt33.Connect()
	assert.Equal(nil, err)

	err = clnt33.Send(msg)
	assert.Equal(nil, err)

	reply := clnt33.Recv()
	assert.Equal(expectRecv, len(reply))
}

// TODO: TestNoServer, client should fail to connect when no server is running

// TODO: test check closed, connection should be closed when forward is removed
