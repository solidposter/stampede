package main

// The client sets Key and Id
// Server responds if Key is valid
// with the same Id and Lport and Hport of the server

type message struct {
	Key   string // secret key
	Id    int    // message ID
	Lport int    // lowest port in range
	Hport int    // highest port in range
}
