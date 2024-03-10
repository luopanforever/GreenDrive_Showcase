package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Client

func InItDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Ping the primary
	// if err := client.Ping(ctx, nil); err != nil {
	// 	log.Fatal(err)
	// }
	print("连接成功\n")
	mongoDB = client
	return client
}

func GetDB() *mongo.Client {
	return mongoDB
}

func Close(db *mongo.Client) {
	if err := db.Disconnect(context.Background()); err != nil {
		log.Fatalf("Error on disconnecting from MongoDB: %v", err)
	}
}
