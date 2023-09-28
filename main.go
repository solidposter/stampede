package main

import (
	"fmt"
	"time"
)

func main() {

	serverconfig := message{
		Key:   "hemlignyckel",
		Id:    0,
		Lport: 10000,
		Hport: 10999,
	}
	fmt.Println(serverconfig)

	s := newServer("9001")
	go s.start(serverconfig)

	c := newClient("2001")
	msg := c.test("127.0.0.1:9001", "hemlignyckel")
	fmt.Println(msg)
	time.Sleep(1 * time.Second)

	c = newClient("2000")
	c.start("127.0.0.1:9001")
	time.Sleep(1 * time.Second)

	c = newClient("4000")
	msg = c.test("127.0.0.1:9001", "hemlignyckel")
	fmt.Println(msg)
	time.Sleep(1 * time.Second)
}
