package main

import (
	"log"
	"net"

	"github.com/juliotorresmoreno/unravel-api/server"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	svr := server.SetupServer()

	svr.RunListener(listener)
}
