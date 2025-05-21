package main

import (
	"context"
	"errors"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			var syscallErr *os.SyscallError
			if errors.As(opErr.Err, &syscallErr) {
				fmt.Printf("Dial error: %#v\n", syscallErr)
				return
			}
			fmt.Printf("Dial error: %#v\n", opErr)
			return
		}
		fmt.Printf("Dial error: %#v\n", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	fmt.Printf("Connected to: %q\n", address)
}
