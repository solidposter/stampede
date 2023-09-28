package main

import (
	"bytes"
	"encoding/gob"
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

func (c *client) test(target string, key string) message {
	var buffer bytes.Buffer

	msg := message{
		Key:   key,
		Id:    0,
		Lport: 0,
		Hport: 0,
	}

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}

	targetAddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenPacket("udp", ":"+c.srcport)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteTo(buffer.Bytes(), targetAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn.Close()

	return msg
}
