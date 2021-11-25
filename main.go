package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	host     = "localhost"
	port     = "4000"
	connType = "tcp"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	fmt.Println("Starting " + connType + " server on " + host + ":" + port)

	l, err := net.Listen(connType, host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	for i := 0; i < 5; i++ {
		go createClient(l, cancel, ctx)
	}
	<-ctx.Done()
	fmt.Println("PROCESS FINISHED")
}

func createClient(l net.Listener, cancel context.CancelFunc, ctx context.Context) {
	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	fmt.Println("Client connected.")
	fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

	handleConnection(c, cancel, l, ctx)
}

func handleConnection(conn net.Conn, cancel context.CancelFunc, l net.Listener, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buffer, err := bufio.NewReader(conn).ReadBytes('\n')

			if err != nil {
				fmt.Println("Client left.")
				conn.Close()

				createClient(l, cancel, ctx)
				return
			}

			log.Println("Client message:", string(buffer[:len(buffer)-1]))

			conn.Write(buffer)
			if string(buffer[:len(buffer)-1]) == "k" {
				cancel()
			}
		}
	}

}
