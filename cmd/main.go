package main

import (
	"context"
	"log"
	"net"
	sku_reader "sku-reader"
	"strconv"
	"time"
)

const (
	timeReading    = 60 * time.Second
	socketHost     = "localhost"
	socketPort     = "4000"
	connType       = "tcp"
	connNumber     = 5
	endSequence    = "terminate"
	mongoHost      = "localhost"
	mongoPort      = "27017"
	username       = "user"
	password       = "password"
	collectionName = "messages"
	database       = "sku_reader"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), timeReading)
	defer cancel()
	log.Println("Starting " + connType + " server on " + socketHost + ":" + socketPort)

	listener, err := net.Listen(connType, socketHost+":"+socketPort)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err.Error())
		}
	}(listener)
	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}

	config := sku_reader.Config{
		Host:           mongoHost,
		Port:           mongoPort,
		UserName:       username,
		Password:       password,
		CollectionName: collectionName,
		Database:       database,
	}

	skuController := sku_reader.InitializeSkuController(context.Background(), listener, config)

	// UUID could be used but since we only can use standard library we use time instead
	sessionId := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	finishReading := make(chan bool)

	for i := 0; i < connNumber; i++ {
		go skuController.HandleConnections(sessionId, endSequence, finishReading)
	}
	for {
		select {
		case <-finishReading:
			skuController.GenerateReport(sessionId)
			log.Println("PROCESS FINISHED")
			return
		case <-ctx.Done():
			skuController.GenerateReport(sessionId)
			log.Println("PROCESS FINISHED")
			return
		}
	}

}
