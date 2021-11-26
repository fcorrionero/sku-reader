package main

import (
	"context"
	"log"
	"net"
	sku_reader "sku-reader"
	"time"
)

const (
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	skuController := sku_reader.InitializeSkuController(ctx, listener, cancel, config)

	// It is better to use a UUID but since we only can use standard library we use time instead
	sessionId := time.Now().String()

	for i := 0; i < connNumber; i++ {
		go skuController.HandleConnections(sessionId, endSequence)
	}
	<-ctx.Done()
	skuController.GenerateReport(sessionId)
	log.Println("PROCESS FINISHED")
}
