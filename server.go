package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type server struct {
	port string
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

		dec := gob.NewDecoder(bytes.NewBuffer(nbuf[:length]))
		err = dec.Decode(&m)
		if err != nil {
			fmt.Println("Server decode error:", err, addr)
			continue
		}
		if m.Key != config.Key {
			fmt.Println(" Server key mismatch", addr)
			continue
		}

		m.Lport = config.Lport
		m.Hport = config.Hport
		buffer := new(bytes.Buffer)
		enc := gob.NewEncoder(buffer)
		err = enc.Encode(m)
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.WriteTo(buffer.Bytes(), addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
