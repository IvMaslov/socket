package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvMaslov/socket"
)

var (
	dev = flag.String("i", "", "device for listening")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	ifce, err := socket.New(socket.WithDevice(*dev))
	if err != nil {
		log.Fatal(err)
	}

	go startReading(ctx, ifce)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	cancel()

	log.Println("Successfully exiting")

	ifce.Close()
}

func startReading(ctx context.Context, i *socket.Interface) {
	buf := make([]byte, 1500)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopped reading from interface")
			return
		default:
			n, err := i.Read(buf)
			if err != nil {
				log.Printf("Error while reading from interface: %v\n", err)
				return
			}

			log.Printf("Msg [%d] - %v\n", n, buf[:n])
		}
	}
}
