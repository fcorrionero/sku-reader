package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	host     = "localhost"
	port     = "4000"
	connType = "tcp"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Starting " + connType + " server on " + host + ":" + port)

	l, err := net.Listen(connType, host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	for i:=0;i<5;i++ {
		createClient(l, cancel)
	}
	<- ctx.Done()
	fmt.Println("PROCESS FINISHED")
}

func createClient(l net.Listener,cancel context.CancelFunc) {
	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	fmt.Println("Client connected.")
	fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

	go handleConnection(c, cancel, l)
}

func handleConnection(conn net.Conn, cancel context.CancelFunc, l net.Listener) {
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			fmt.Println("Client left.")
			conn.Close()

			createClient(l, cancel)
			return
		}

		log.Println("Client message:", string(buffer[:len(buffer)-1]))

		conn.Write(buffer)
		if string(buffer[:len(buffer)-1]) == "k" {
			cancel()
		}
	}

}