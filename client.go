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
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/solidposter/stampede/pb"
	"google.golang.org/protobuf/proto"
)

type client struct {
	srcport string
}

func newClient(srcport string) *client {
	return &client{
		srcport: srcport,
	}
}

func (c *client) start(targetIP string, config configuration) {
	nbuf := make([]byte, 128)
	pbreq := &pb.Payload{
		Key:   config.Key,
		Id:    uint64(rand.Int63()),
		Hport: uint32(config.Hport),
		Lport: uint32(config.Lport),
	}
	pbresp := &pb.Payload{}

	conn, err := net.ListenPacket("udp", ":"+c.srcport)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for dport := config.Lport; dport <= config.Hport; dport++ {
			if pbreq.Id == math.MaxUint64 {
				pbreq.Id = 0
			} else {
				pbreq.Id += 1
			}

			outbytes, err := proto.Marshal(pbreq)
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

			success := false // set to true for valid response
			for {
				_, err = conn.WriteTo(outbytes, targetAddr)
				if err != nil {
					log.Fatal(err)
				}

				conn.SetReadDeadline((time.Now().Add(1000 * time.Millisecond)))
				for {
					length, addr, err := conn.ReadFrom(nbuf)
					if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
						fmt.Print(".")
						break
					}
					if err != nil {
						log.Fatal(err)
					}

					if addr.String() != targetAddr.String() {
						log.Printf("Packet received from invalid source %v with length %v", addr, length)
						continue
					}

					if err := proto.Unmarshal(nbuf[0:length], pbresp); err != nil {
						log.Print("Server decode error:", err, addr)
						continue
					}

					if pbresp.Key != pbreq.Key {
						log.Print("Invalid key in response:", pbresp.Key)
						continue
					}
					if pbresp.Id != pbreq.Id {
						log.Printf("Incorrect Id, expected %v got %v\n", pbreq.Id, pbresp.Id)
						continue
					}
					success = true // valid response
					break
				}
				if success {
					break // move on to next port
				}
			}
		}
	}
}

// Returns a message struct with server configuration
func (c *client) probe(target string, key string) configuration {
	inbytes := make([]byte, 128)
	pbresp := &pb.Payload{}

	pbreq := &pb.Payload{
		Key:   key,
		Id:    uint64(rand.Int63()),
		Lport: 0,
		Hport: 0,
	}
	outbytes, err := proto.Marshal(pbreq)
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

	success := false // set to true for valid response
	for {
		_, err = conn.WriteTo(outbytes, targetAddr)
		if err != nil {
			log.Fatal(err)
		}

		conn.SetReadDeadline((time.Now().Add(1000 * time.Millisecond)))
		for {
			length, addr, err := conn.ReadFrom(inbytes)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				fmt.Print(".")
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if addr.String() != targetAddr.String() {
				log.Printf("Packet received from invalid source %v with length %v", addr, length)
				continue
			}

			if err := proto.Unmarshal(inbytes[0:length], pbresp); err != nil {
				log.Print("Server decode error:", err, addr)
				continue
			}

			if pbresp.Key != pbreq.Key {
				log.Print("Invalid key in response:", pbresp.Key)
				continue
			}
			if pbresp.Id != pbreq.Id {
				log.Printf("Incorrect Id, expected %v got %v\n", pbreq.Id, pbresp.Id)
				continue
			}
			success = true // valid response
			break
		}
		if success {
			break
		}
	}
	conn.Close()

	resp := configuration{}
	resp.Key = key
	resp.Hport = int(pbresp.Hport)
	resp.Lport = int(pbresp.Lport)
	return resp
}
