package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/vlaner/cache-server/server"
)

func main() {
	s := server.New(":8080")

	sigCh := make(chan os.Signal, 1)
	log.Println("Starting server...")

	s.Start()

	signal.Notify(sigCh, os.Interrupt)
	sig := <-sigCh

	s.Stop()

	log.Printf("Received %s at %s, shutting down...\n", sig, time.Now())
}
