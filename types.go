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

// The client sets Key and Id
// Server responds if Key is valid
// with the same Id and Lport and Hport of the server

// The server configuration and the req/resp uses the same struct.
// Server starts listeners on the from Lport to Hport.
// The clients send req with Key and ID set, server responds and
// sets Lport and Hport in the response.

type message struct {
	Key   string // secret key
	Id    int    // message ID
	Lport int    // lowest port in range
	Hport int    // highest port in range
}
