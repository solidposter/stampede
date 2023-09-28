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
		c.start(targetIP, serverconfig)
	}
	fmt.Println(" to port", *portPtr+*numPtr-1)
	<-(chan int)(nil) // wait forever

	/*
		serverconfig := message{
			Key:   "hemlignyckel",
			Id:    0,
			Lport: 10000,
			Hport: 10999,
		}
	*/

	/*
		fmt.Println("main serverconfig:", serverconfig)
		time.Sleep(1 * time.Second)

		fmt.Println("start servers on ports 10000-109999")
		for i := 10000; i < 11000; i++ {
			s := newServer(strconv.Itoa(i))
			go s.start(serverconfig)
		}
		// s := newServer("9001")
		// go s.start(serverconfig)

		c := newClient("2001")
		msg := c.probe("127.0.0.1:10000", *keyPtr)
		fmt.Println("main server response:", msg)
		time.Sleep(5 * time.Second)

		c.start("127.0.0.1", msg)

		fmt.Println("main done, sleeping 10s")
		time.Sleep(10 * time.Second)
	*/
}
