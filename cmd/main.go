package main

import (
	"context"
	"log"
	"net"
	sku_reader "sku-reader"
	"strconv"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), sku_reader.TimeReading)
	defer cancel()
	log.Println("Starting " + sku_reader.ConnType + " server on " + sku_reader.SocketHost + ":" + sku_reader.SocketPort)

	listener, err := net.Listen(sku_reader.ConnType, sku_reader.SocketHost+":"+sku_reader.SocketPort)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err.Error())
		}
	}(listener)

	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}

	ctxConnections, cancelConnections := context.WithCancel(context.Background())
	defer cancelConnections()
	skuController := sku_reader.InitializeSkuController(ctxConnections, listener)

	// UUID could be used but since we only can use standard library we use time instead
	sessionId := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	finishReading := make(chan bool)
	errorStream := make(chan interface{})

	for i := 0; i < sku_reader.ConnNumber; i++ {
		go skuController.HandleConnections(sessionId, sku_reader.EndSequence, finishReading, errorStream)
	}
	for {
		select {
		case err := <-errorStream:
			// Proper error handling should be added, metrics server, etc
			log.Fatal(err)
		case <-finishReading:
			log.Println(skuController.GenerateReport(sessionId))
			log.Println("PROCESS FINISHED")
			return
		case <-ctx.Done():
			log.Println(skuController.GenerateReport(sessionId))
			log.Println("PROCESS FINISHED")
			return
		}
	}

}
