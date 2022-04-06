package main

import (
	"flag"

	"github.com/rafaelfuchsbr/tcp-service/internal/server"
)

func main() {

	portPtr := flag.Int("port", 12345, "")
	flag.Parse()

	server := server.NewServer(*portPtr)
	err := server.Serve()
	if err != nil {
		panic(err)
	}
}
