package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Austinhamilton1/dnsproxy/internal/server"
)

func main() {
	portPtr := flag.Int("port", 5353, "port to run the server on")
	blockFilePtr := flag.String("blocked", "", "name of the file to add blocked IPs from")

	flag.Parse()

	port := *portPtr
	connStr := fmt.Sprintf("127.0.0.1:%d", port)

	s := server.New(connStr, *blockFilePtr)

	if err := s.Run(); err != nil {
		log.Fatalf("could not create DNS proxy: %s", err)
	}
}
