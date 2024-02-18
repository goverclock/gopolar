package core

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
)

var homeDir string

type ConnLogger struct {
	sendf *os.File
	recvf *os.File
}

// remove all existing logs when gpcore starts
func init() {
	hd, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("fail to read $HOME dir: ", err)
	}
	homeDir = hd
	os.RemoveAll(fmt.Sprintf("%v/.config/gopolar/logs", homeDir))
}

// log file at ~/.config/gopolar/logs/
func NewConnLogger(source string, dest string) *ConnLogger {
	current := time.Now()

	logDir := fmt.Sprintf("%v/.config/gopolar/logs/%v-%v/", homeDir, source, dest)
	os.MkdirAll(logDir, 0700)

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
		sendf: sendLogFile,
		recvf: recvLogFile,
	}
}

func (cl *ConnLogger) LogSend(b []byte) {
	// TODO: if !doLog { return }
	err := binary.Write(cl.sendf, binary.LittleEndian, b)
	if err != nil {
		log.Fatal("[logger] LogSend fail")
	}
}

func (cl *ConnLogger) LogRecv(b []byte) {
	// TODO: if !doLog { return }
	err := binary.Write(cl.recvf, binary.LittleEndian, b)
	if err != nil {
		log.Fatal("[logger] LogRecv fail")
	}
}
