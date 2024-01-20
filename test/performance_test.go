package gopolar_test

import (
	"fmt"
	"gopolar/internal/core"
	"gopolar/test/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// returns n MB bytes without zero
func MakeDataMB(n uint) []byte {
	size := n * 1024 * 1024
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte('a' + (i % 26))
	}
	return buf
}

func TestDirect100MB(t *testing.T) {
	assert := assert.New(t)

	prefix := "hello"
	serv88 := testutil.NewEchoServer(88, prefix)
	defer serv88.Quit()

	clnt88 := testutil.NewEchoClient(88)
	assert.Nil(clnt88.Connect())

	data := string(MakeDataMB(100))
	msg := fmt.Sprintf("data from client: %v\n", data)
	assert.Nil(clnt88.Send(msg))
	reply := clnt88.Recv()
	assert.Equal(prefix+msg, reply)
}

// same with TestOne2One, but with 100 MB data
func TestForward100MB(t *testing.T) {
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
	assert.Nil(clnt33.Connect())

	data := string(MakeDataMB(100))
	msg := fmt.Sprintf("data from client: %v\n", data)
	assert.Nil(clnt33.Send(msg))
	reply := clnt33.Recv()
	assert.Equal(prefix+msg, reply)
}

// TODO: TestManyManyConnections()
