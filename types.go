package main

type message struct {
	key   string // secret key
	id    int    // message ID
	lport int    // lowest port in range
	hport int    // highest port in range
}
