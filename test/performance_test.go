package gopolar_test

import (
	"fmt"
	"gopolar/internal/core"
	"gopolar/test/testutil"
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

// connect, send, recv, verify, disconnect
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

// TODO: TestManyManyConnections()
