package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Client

func InItDB() *mongo.Client {
	// 使用云mongodb
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI("mongodb+srv://luopan:luopan@tdcars.spuljs2.mongodb.net/?retryWrites=true&w=majority&appName=tdCars").SetServerAPIOptions(serverAPI)
	// client, err := mongo.Connect(context.TODO(), opts)
	// if err != nil {
	// 	panic(err)
	// }

	// 使用本地mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

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
