package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type server struct {
	port string // low port in range
}

func newServer(port string) *server {
	return &server{
		port: port,
	}
}

func (s *server) start(config message) {
	m := message{}
	nbuf := make([]byte, 1500)

	conn, err := net.ListenPacket("udp", ":"+s.port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		length, addr, err := conn.ReadFrom(nbuf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("received packet", addr, length)

		dec := gob.NewDecoder(bytes.NewBuffer(nbuf[:length]))
		err = dec.Decode(&m)
		if err != nil {
			fmt.Println("decode error:", err) // do I care ?
			continue
		}
		fmt.Println(m)
		if m.Key != config.Key {
			fmt.Println("Key mismatch")
			continue
		}
		fmt.Println("key match")
	}
}
