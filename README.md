## echo server

a toy udp server made to echo your message back to your socket.

when running the program, it automatically initializes a "client" socket connection which sends random words to the "server" socket.

the server uses go's goroutines for concurrency, you can test this out by sending a datagram to the socket the server is running on:

```bash
echo "hey im a new client" | nc -u 127.0.0.1 PORT_THE_SERVER_IS_LISTENING_ON
```

graceful termination is handled. (with context lib)

## stuff learned

context.Context, graceful cancellation, goroutines, concurrency, error handling, socket programming, net lib, communication via udp, fmt package, when to use log vs fmt,
