package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
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
	nbuf := make([]byte, 1500)

	conn, err := net.ListenPacket("udp", ":"+c.srcport)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for dport := config.Lport; dport < config.Hport; dport++ {

			config.Id += 1
			buffer := new(bytes.Buffer)
			enc := gob.NewEncoder(buffer)
			err := enc.Encode(config)
			if err != nil {
				log.Fatal(err)
			}

			s := strconv.Itoa(dport)
			if err != nil {
				log.Panic(err)
			}
			targetAddr, err := net.ResolveUDPAddr("udp", targetIP+":"+s)
			if err != nil {
				log.Fatal(err)
			}
			for {
				conn.SetReadDeadline((time.Now().Add(1000 * time.Millisecond)))
				_, err = conn.WriteTo(buffer.Bytes(), targetAddr)
				if err != nil {
					log.Fatal(err)
				}

				length, addr, err := conn.ReadFrom(nbuf)
				if err != nil {
					fmt.Println("client error", err, addr, length)
					continue
				}
				break
			}
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
