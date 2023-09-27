package main

import (
	"log"
	"net"
)

type client struct {
	srcport string
}

func newClient(srcport string) *client {
	return &client{
		srcport: srcport,
	}
}

func (c *client) start(target string) {
	data := make([]byte, 100) // test junk

	targetAddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("udp", ":"+c.srcport)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.WriteTo(data, targetAddr)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
