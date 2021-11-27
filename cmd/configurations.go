package main

import "time"

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
