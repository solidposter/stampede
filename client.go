package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

type client struct {
	srcport string
}

func newClient(srcport string) *client {
	return &client{
		srcport: srcport,
	}
}

func (c *client) start(targetIP string, config message) {
	var buffer bytes.Buffer
	nbuf := make([]byte, 1500)

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(config)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("udp", ":"+c.srcport)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for dport := config.Lport; dport < config.Hport; dport++ {
			s := strconv.Itoa(dport)
			if err != nil {
				log.Panic(err)
			}
			targetAddr, err := net.ResolveUDPAddr("udp", targetIP+":"+s)
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.WriteTo(buffer.Bytes(), targetAddr)
			if err != nil {
				log.Fatal(err)
			}

			length, addr, err := conn.ReadFrom(nbuf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("client received packet from", addr, length)
		}
	}

}

// Use probe to get server port range in a struct message
func (c *client) probe(target string, key string) message {
	var buffer bytes.Buffer
	nbuf := make([]byte, 1500)
	m := message{}

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
	length, addr, err := conn.ReadFrom(nbuf)
	if err != nil {
		log.Fatal(err)
	}

	dec := gob.NewDecoder(bytes.NewBuffer(nbuf[:length]))
	err = dec.Decode(&m)
	if err != nil {
		fmt.Println("Client decode error:", err, addr)
	}

	conn.Close()
	return m
}
