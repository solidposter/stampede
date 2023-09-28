package main

import (
	"fmt"
	"time"
)

func main() {

	s := newServer("hemlignyckel", "9001")
	go s.start()

	c := newClient("2000")
	c.start("127.0.0.1:9001")
	time.Sleep(1 * time.Second)

	c = newClient("3000")
	c.start("127.0.0.1:9001")
	time.Sleep(1 * time.Second)

	c = newClient("4000")
	msg := c.test("127.0.0.1:9001", "hemlignyckel")
	fmt.Println(msg)
	time.Sleep(1 * time.Second)
}
