package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {

	modePtr := flag.Bool("s", false, "Set server mode")
	portPtr := flag.Int("p", 10000, "Base port")
	numPtr := flag.Int("n", 1000, "Number ports")
	keyPtr := flag.String("k", "hemlignyckel", "Key")
	flag.Parse()

	// Server mode
	if *modePtr {
		fmt.Println("server mode")
		serverconfig := message{
			Key:   *keyPtr,
			Id:    0,
			Lport: *portPtr,
			Hport: *portPtr + *numPtr,
		}

		fmt.Print("Starting servers from port ", *portPtr)
		for i := *portPtr; i < *portPtr+*numPtr; i++ {
			s := newServer(strconv.Itoa(i))
			go s.start(serverconfig)
		}
		fmt.Println(" to port", *portPtr+*numPtr-1)
		<-(chan int)(nil) // wait forever
	}

	if len(flag.Args()) != 1 {
		fmt.Println("please specify target, ip:port")
		return
	}

	// Probe server for port configuration
	target := flag.Args()[0]
	fmt.Println("Starting initial Probe of", target)
	c := newClient(strconv.Itoa(*portPtr))
	serverconfig := c.probe(target, *keyPtr)
	fmt.Println(serverconfig)

	ip, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		log.Panic(err)
	}
	targetIP := ip.IP.String()

	// start the clients
	fmt.Print("Starting clients from port ", *portPtr)
	for i := *portPtr; i < *portPtr+*numPtr; i++ {
		c := newClient(strconv.Itoa(i))
		go c.start(targetIP, serverconfig)
	}
	fmt.Println(" to port", *portPtr+*numPtr-1)
	<-(chan int)(nil) // wait forever
}
