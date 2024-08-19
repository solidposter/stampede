# stampede

I use stampede to load-test concurrent active sessions in firewalls.

This example runs UDP traffic over 1M sessions.
### Server:
stampede -s -k secretkey -p 20000 -n 2000

### Client:
stampede -k secretkey -p 15000 -n 500 x.x.x.x:20000

This server listens on ports 20000-201999 (2000 ports).
This client runs 500 goroutines using source ports 15000-15499 (500 ports).

The client goroutines will iterate over the server port range at the speed of RTT.

The numbers can be tuned a lot, there is no problem running the server with 50k ports active as long as the ports are free.
