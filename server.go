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
	"log"
	"net"

	"github.com/solidposter/stampede/pb"
	"google.golang.org/protobuf/proto"
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
	pbreq := &pb.Payload{}
	nbuf := make([]byte, 128)

	conn, err := net.ListenPacket("udp", ":"+s.port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		length, addr, err := conn.ReadFrom(nbuf)
		if err != nil {
			log.Fatal(err)
		}

		if err := proto.Unmarshal(nbuf[0:length], pbreq); err != nil {
			log.Print("Server decode error:", err, addr, length)
			continue
		}
		if pbreq.Key != config.Key {
			log.Println("Server key mismatch", addr, pbreq.Key, config.Key)
			continue
		}

		pbreq.Lport = uint32(config.Lport)
		pbreq.Hport = uint32(config.Hport)
		nbuf, err := proto.Marshal(pbreq)
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.WriteTo(nbuf, addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
