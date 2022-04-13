package main

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/server"
	"log"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
