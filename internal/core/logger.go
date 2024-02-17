package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ConnLogger struct {
	send *log.Logger
	recv *log.Logger
}

// log file at ~/.config/gopolar/logs/
func NewConnLogger(source string, dest string) *ConnLogger {
	current := time.Now()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("fail to read $HOME dir: ", err)
	}
	logDir := fmt.Sprintf("%v/.config/gopolar/logs/%v-%v/", homeDir, source, dest)

	sendLogName := fmt.Sprintf("%v-send", current.Format("2006-01-02 15:04:05.000000"))
	recvLogName := fmt.Sprintf("%v-recv", current.Format("2006-01-02 15:04:05.000000"))

	sendLogFile, err := os.OpenFile(logDir+sendLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("fail to create log file: %v ,err=%v \n", sendLogName, err)
	}
	recvLogFile, err := os.OpenFile(logDir+recvLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("fail to create log file: %v ,err=%v \n", recvLogName, err)
	}

	return &ConnLogger{
		send: log.New(sendLogFile, "", 0),
		recv: log.New(recvLogFile, "", 0),
	}
}

func (cl *ConnLogger) LogSend(b []byte) {
	// TODO: if !doLog { return }
	cl.send.Print(b) // TODO: should output raw bytes
}

func (cl *ConnLogger) LogRecv(b []byte) {
	// TODO: if !doLog { return }
	cl.recv.Print(b) // TODO: should output raw bytes
}
