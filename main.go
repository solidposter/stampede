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
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {

	modePtr := flag.Bool("s", false, "Set server mode")
	portPtr := flag.Int("p", 20000, "Base port")
	numPtr := flag.Int("n", 10, "Number ports")
	keyPtr := flag.String("k", "hemlignyckel", "Key")
	flag.Parse()

	// Server mode
	if *modePtr {
		serverconfig := message{
			Key:   *keyPtr,
			Id:    0,
			Lport: *portPtr,
			Hport: *portPtr + *numPtr - 1,
		}

		fmt.Printf("Starting servers on ports %v-", *portPtr)
		for i := *portPtr; i < *portPtr+*numPtr; i++ {
			s := newServer(strconv.Itoa(i))
			go s.start(serverconfig)
		}
		fmt.Printf("%v\n", *portPtr+*numPtr-1)
		<-(chan int)(nil) // wait forever
	}

	if len(flag.Args()) != 1 {
		fmt.Println("please specify target, ip:port")
		return
	}

	// Probe server for port configuration
	target := flag.Args()[0]
	fmt.Printf("Starting probe of %v: ", target)
	c := newClient(strconv.Itoa(*portPtr))
	serverconfig := c.probe(target, *keyPtr)
	fmt.Printf("ports %v-%v active\n", serverconfig.Lport, serverconfig.Hport)

	ip, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		log.Panic(err)
	}
	targetIP := ip.IP.String()

	// start the clients
	fmt.Printf("Starting clients on ports %v-", *portPtr)
	for i := *portPtr; i < *portPtr+*numPtr; i++ {
		c := newClient(strconv.Itoa(i))
		go c.start(targetIP, serverconfig)
	}
	fmt.Printf("%v: %v UDP sessions\n", *portPtr+*numPtr-1, (serverconfig.Hport-serverconfig.Lport+1)**numPtr)
	<-(chan int)(nil) // wait forever
}
