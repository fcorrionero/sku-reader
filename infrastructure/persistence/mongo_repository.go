package persistence

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sku-reader/domain"
)

type MessageMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
}

func NewMongoRepository(client *mongo.Client, collection *mongo.Collection, ctx context.Context) MessageMongoRepository {
	return MessageMongoRepository{
		client:     client,
		collection: collection,
		ctx:        ctx,
	}
}

func (repository MessageMongoRepository) Save(message *domain.Message) error {
	if repository.client == nil {
		return errors.New("error trying to connect to database")
	}

	_, err := repository.collection.InsertOne(repository.ctx, message)

	return err
}

func (repository MessageMongoRepository) FindAll(id string) []*domain.Message {
	var messages []*domain.Message

	filter := bson.M{"sessionId": id}
	cursor, err := repository.collection.Find(repository.ctx, filter)
	if err != nil {
		log.Println(err)
		return messages
	}

	for cursor.Next(repository.ctx) {
		var message domain.Message
		err := cursor.Decode(&message)
		if err != nil {
			log.Println(err)
			return messages
		}
		messages = append(messages, &message)
	}

	return messages
}
