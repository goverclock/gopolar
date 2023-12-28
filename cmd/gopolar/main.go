package main

import (
	"log"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setup() {
	log.SetPrefix("[core]")
	log.SetFlags(0)
}

func main() {
	setup()
	// homeDir, err := os.UserHomeDir()
	// configDir := homeDir + "/.gopolar"
	// check(err)
	// err = os.MkdirAll(configDir, os.ModePerm)
	// check(err)
	// f, err := os.Create(configDir + "/gopolar.sock")
	// check(err)
	// defer f.Close()
	// sock := os.Getpid()	// todo
	// f.WriteString(strconv.Itoa(sock))

	go handler()

	for {
	}
}
