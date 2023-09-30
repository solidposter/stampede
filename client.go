package main

//
// Copyright (c) 2023 Tony Sarendal <tony@polarcap.org>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

import (
	"bytes"
	"encoding/json"
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

func (c *client) start(targetIP string, req message) {
	nbuf := make([]byte, 1500)

	conn, err := net.ListenPacket("udp", ":"+c.srcport)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for dport := req.Lport; dport <= req.Hport; dport++ {

			req.Id += 1
			buffer := new(bytes.Buffer)
			enc := json.NewEncoder(buffer)
			err := enc.Encode(req)
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
				_, err = conn.WriteTo(buffer.Bytes(), targetAddr)
				if err != nil {
					log.Fatal(err)
				}

				conn.SetReadDeadline((time.Now().Add(1000 * time.Millisecond)))
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
	resp := message{}

	req := message{
		Key:   key,
		Id:    0,
		Lport: 0,
		Hport: 0,
	}

	enc := json.NewEncoder(&buffer)
	err := enc.Encode(req)
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

	success := false // set to true with succesful decode of response
	for {
		_, err = conn.WriteTo(buffer.Bytes(), targetAddr)
		if err != nil {
			log.Fatal(err)
		}

		conn.SetReadDeadline((time.Now().Add(1000 * time.Millisecond)))
		for {
			length, addr, err := conn.ReadFrom(nbuf)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				// Timeout reached,print a dot and resend request
				fmt.Print(".")
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if addr.String() != targetAddr.String() {
				// some other packet hit the port, ignore and go back to listen
				log.Printf("Packet received from incorrect IP address %v with length %v", addr, length)
				continue
			}
			dec := json.NewDecoder(bytes.NewBuffer(nbuf[:length]))
			err = dec.Decode(&resp)
			if err != nil {
				log.Print("Client decode error:", err)
				continue
			}
			success = true
			break
		}

		if success {
			break
		}

	}
	conn.Close()
	return resp
}
