package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	mrand "math/rand/v2"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	serverAddr, err := echoUDPServer(ctx, "127.0.0.1:")
	if err != nil {
		log.Printf("starting echo udp server: %v", err)
		return
	}
	fmt.Printf("started server on endpoint %q:\n", serverAddr)

	err = udpClient(ctx, serverAddr)
	if err != nil {
		log.Printf("udp client: %v", err)
	}
}

func udpClient(ctx context.Context, addr net.Addr) error {
	s, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		return fmt.Errorf("connecting to server: %v", err)
	}
	fmt.Printf("hi i'm a client which sends udp packages\n my endpoint is %q\n", s.LocalAddr())

	// cancel this program whenever ctrl c is sent gracefully
	go func() {
		<-ctx.Done()
		log.Println("gracefully shutting down client...")
		_ = s.Close()
	}()

	for {
		time.Sleep(time.Second)
		_, err := s.WriteTo([]byte(randomWord()), addr)
		if err != nil {
			if ctx.Err() != nil && ctx.Err().Error() == "context canceled" {
				return ctx.Err()
			}

			log.Printf("sending msg to endpoint %q from %q: %v", s.LocalAddr(), addr, err)
			continue
			// not returning here cuz the program shouldnt stop unless i cancel it
		}

		buf := make([]byte, 1024)
		n, readAddr, err := s.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("reading msg from %q: %v", readAddr, err)
		}
		msg := buf[:n]

		if readAddr.String() != addr.String() {
			log.Printf("stray message, endpoint %q does not match server's addr %q", readAddr, addr)
			log.Printf("stray msg: %s", string(msg))
		}

		fmt.Printf("CLIENT: the echo server echoed back %q\n", msg)

	}
}

// ctx should preferrably be a cancellable context; addr = 127.0.0.1: makes the OS choose the port the server listens on
func echoUDPServer(ctx context.Context, addr string) (net.Addr, error) {
	// obviouslt i need to get the socket somehow, this is how we do it. addr is a form of light dependency injection
	s, err := net.ListenPacket("udp", addr)
	// handle the error, write only the part that we did rn (listening to the udp endpoint)
	if err != nil {
		return nil, fmt.Errorf("listening to udp endpoint %q: %w", addr, err)
	}

	// im starting a goroutine because I want this to not block the main thread
	go func() {
		// if the context is cancelled, then stop blocking on this goroutine and close the socket
		go func() {
			<-ctx.Done()
			// println so the ^C is above the logs and not shoved into them in the terminal
			fmt.Println("")
			log.Println("gracefully shutting down echo server...")
			_ = s.Close()
		}()

		// make allocates and initializes an object of type Type
		buf := make([]byte, 1024)

		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}
			fmt.Printf("SERVER: %q  (client endpoint) -- %q (msg)\n", clientAddr, string(buf[:n]))

			// TODO: add some state which checks if this is the first message you sent to this endpoint, if not, just send it
			// if its the first, send a hello message explaining what the server does
			if _, err = s.WriteTo(buf[:n], clientAddr); err != nil {
				return
			}

		}
	}()

	return s.LocalAddr(), nil
}

func gibberish() []byte {
	thing := make([]byte, 5)
	// the error crashes the program if it errors anyways
	_, _ = rand.Read(thing)

	return thing
}

func randomWord() string {
	words := []string{"grug", "ssj2", "ssj4", "concurrency", "go", "written in rust btw", "cloud native :geek emoji:", "cargo cult fan", "idempotent", "paralellism", "htmx", "reactjs", "k8s", "aposd", "usd"}

	return words[mrand.IntN(len(words))]
}
