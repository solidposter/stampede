# stampede

I used stampede to load-test concurrent active sessions in firewalls.

This test runs UDP traffic over 1M sessions.
### Server:
stampede -s -k secretkey -p 20000 -n 2000

### Client:
stampede -k secretkey -p 15000 -n 500 x.x.x.x:20000

The server listens on ports 20000-201999 (2000 ports) and bounces back requests.
The client runs 500 goroutines using source ports 15000-15499 (500 ports).

Each of the client goroutines will iterate over the server port range.
The above configuration will run traffic at the speed of RTT across
1M active sessions (2000*500).

