package main

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	key  string // secret key
	port string // low port in range
}

func newServer(key, port string) *server {
	return &server{
		key:  key,
		port: port,
	}
}

func (s *server) start() {
	conn, err := net.ListenPacket("udp", ":"+s.port)
	if err != nil {
		log.Fatal(err)
	}

	nbuf := make([]byte, 1500)
	for {
		length, addr, err := conn.ReadFrom(nbuf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("received packet", addr, length)
	}
}
