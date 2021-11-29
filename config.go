package sku_reader

import "time"

/*
We use this config file instead of .env file due to for this project only standard library is allowed and the package
github.com/joho/godotenv is not in it. Otherwise, we would have to implement custom functions to read .env files.
*/
const (
	TimeReading    = 60 * time.Second
	SocketHost     = "localhost"
	SocketPort     = "4000"
	ConnType       = "tcp"
	ConnNumber     = 5
	EndSequence    = "terminate"
	MongoHost      = "localhost"
	MongoPort      = "27017"
	Username       = "user"
	Password       = "password"
	CollectionName = "messages"
	Database       = "sku_reader"
)
