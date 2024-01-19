package core

import (
	"log"
	"os"
)

var debug bool = false

func init() {
	if os.Getenv("DEBUG") == "1" {
		debug = true
	}
}

func Debugln(args ...interface{}) {
	if debug {
		log.Println(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
}
