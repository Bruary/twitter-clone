package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UsersCol *mongo.Collection

func GetDBConn() *mongo.Database {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database("twitter")
}

func InsertDocumentToDB(dbCollection *mongo.Collection, dataToStore interface{}) error {

	_, err := dbCollection.InsertOne(context.TODO(), dataToStore)
	if err != nil {
		return err
	}

	return nil
}

func UserAlreadyExists(dbCollection *mongo.Collection, usersEmail string) (bool, error) {
	count, err := dbCollection.CountDocuments(context.TODO(), bson.M{"email": usersEmail})
	if count >= 1 {
		return true, nil
	} else {
		return false, err
	}
}

func DeleteUser(dbCollection *mongo.Collection, usersEmail string) error {
	_, err := dbCollection.DeleteOne(context.TODO(), bson.M{"email": usersEmail})
	if err != nil {
		return err
	}

	return nil
}
