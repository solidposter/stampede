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

		dec := json.NewDecoder(bytes.NewBuffer(nbuf[:length]))
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
		enc := json.NewEncoder(buffer)
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
