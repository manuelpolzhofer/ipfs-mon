package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Initializing ipfs-mon...")
	ctx, cancel := context.WithCancel(context.Background())
	shutdownHandler(cancel)
	c := NewCluster()
	c.Start(ctx)
}

func shutdownHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Starting Shutdown")
		cancel()
	}()
}
