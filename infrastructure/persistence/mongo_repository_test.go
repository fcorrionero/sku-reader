//go:build integration
// +build integration

package persistence_test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sku-reader/domain"
	"sku-reader/infrastructure/persistence"
	"testing"
)

func createMongoClient(ctx context.Context) *mongo.Client {
	credentials := options.Credential{Username: "user", Password: "password"}

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", "localhost", "27017")).SetAuth(credentials))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func createMongoCollection(mongoClient *mongo.Client) *mongo.Collection {
	return mongoClient.Database("sku_reader").Collection("messages_test")
}

func TestMessageShouldBeSaved(t *testing.T) {
	ctx := context.Background()
	mongoClient := createMongoClient(ctx)
	collection := createMongoCollection(mongoClient)

	defer cleanUp(mongoClient, collection, ctx)

	repository := persistence.NewMongoRepository(mongoClient, collection, ctx)

	testSku := "test-1111"
	message := domain.NewMessage("1", testSku)

	err := repository.Save(&message)
	if err != nil {
		t.Fatalf("error saving message: %v", err.Error())
	}

	var result domain.Message
	err = collection.FindOne(ctx, bson.M{"sku": testSku}).Decode(&result)
	if err != nil {
		t.Fatal(err)
		return
	}
	if result.Sku != message.Sku {
		t.Fatalf("error finding message, expected: %v, got: %v", message.Sku, result.Sku)
	}
}

func TestSaveMethodShouldThrowAnError(t *testing.T) {
	ctx := context.Background()
	mongoClient := createMongoClient(ctx)
	collection := createMongoCollection(mongoClient)

	defer cleanUp(mongoClient, collection, ctx)

	repository := persistence.NewMongoRepository(nil, collection, ctx)
	testSku := "test-1111"
	message := domain.NewMessage("1", testSku)

	err := repository.Save(&message)
	if err == nil {
		t.Fatalf("save method should throw an error if there is no client")
	}
}

func TestMessagesShouldBeFound(t *testing.T) {
	ctx := context.Background()
	mongoClient := createMongoClient(ctx)
	collection := createMongoCollection(mongoClient)

	defer cleanUp(mongoClient, collection, ctx)
	id := "1"

	message1 := domain.NewMessage(id, "test-1111")
	message2 := domain.NewMessage(id, "test-2222")
	_, err := collection.InsertOne(ctx, &message1)
	if err != nil {
		t.Fail()
	}
	_, err = collection.InsertOne(ctx, &message2)
	if err != nil {
		t.Fail()
	}

	repository := persistence.NewMongoRepository(mongoClient, collection, ctx)

	messages := repository.FindAll(id)
	if len(messages) != 2 {
		t.Fatalf("expected to find %v messages, got %v", 2, len(messages))
	}

}

func cleanUp(client *mongo.Client, collection *mongo.Collection, ctx context.Context) {

	err := collection.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
