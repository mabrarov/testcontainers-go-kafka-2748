package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <host:port>")
		return
	}
	address := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		fmt.Printf("Dial error: %#v\n", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	fmt.Printf("Connected to: %q\n", address)
}
