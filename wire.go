//go:build wireinject

package sku_reader

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"sku-reader/application"
	"sku-reader/domain"
	"sku-reader/infrastructure/persistence"
	"sku-reader/infrastructure/socket"
)

type Config struct {
	Host           string
	Port           string
	UserName       string
	Password       string
	CollectionName string
	Database       string
}

func createMongoDbClient(ctx context.Context, cfg Config) *mongo.Client {
	credentials := options.Credential{Username: cfg.UserName, Password: cfg.Password}

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)).SetAuth(credentials))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func initializeMongoDbCollection(mongoClient *mongo.Client, cfg Config) *mongo.Collection {
	return mongoClient.Database(cfg.Database).Collection(cfg.CollectionName)
}

func initializeMongoDbRepository(
	ctx context.Context,
	cfg Config) persistence.MessageMongoRepository {
	wire.Build(createMongoDbClient, initializeMongoDbCollection, persistence.NewMongoRepository)
	return persistence.MessageMongoRepository{}
}

func initializeCreateMessageCommandHandler(repository domain.MessageRepository) application.CreateMessageCommandHandler {
	wire.Build(application.NewCreateMessageCommandHandler)
	return application.CreateMessageCommandHandler{}
}

func initializeGenerateReportQueryHandler(repository domain.MessageRepository) application.GenerateReportQueryHandler {
	wire.Build(application.NewGenerateReportQueryHandler)
	return application.GenerateReportQueryHandler{}
}

func getConfig() Config {
	return Config{
		Host:           MongoHost,
		Port:           MongoPort,
		UserName:       Username,
		Password:       Password,
		CollectionName: CollectionName,
		Database:       Database,
	}
}

func InitializeSkuController(
	ctx context.Context,
	listener net.Listener) socket.SkuController {
	wire.Build(
		getConfig,
		initializeMongoDbRepository,
		wire.Bind(new(domain.MessageRepository), new(persistence.MessageMongoRepository)),
		initializeCreateMessageCommandHandler,
		initializeGenerateReportQueryHandler,
		socket.NewSkuController)
	return socket.SkuController{}
}
