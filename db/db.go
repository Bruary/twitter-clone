package db

import (
	"context"
	"fmt"
	"log"

	models "github.com/Bruary/twitter-clone/service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbConn *mongo.Database
var UsersCol *mongo.Collection
var TweetsCol *mongo.Collection
var FollowersCol *mongo.Collection

func SetUpDBConnection() {

	// establish a connection to the db
	dbConn = GetDBConn()

	// connect to the following connections after establishing a connection to the db
	UsersCol = ConnectToUsersCol()
	TweetsCol = ConnectToTweetsCol()
	FollowersCol = ConnectToFollowersCol()
}

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

func ConnectToUsersCol() *mongo.Collection {
	return dbConn.Collection("Users")
}

func ConnectToTweetsCol() *mongo.Collection {
	return dbConn.Collection("Tweets")
}

func ConnectToFollowersCol() *mongo.Collection {
	return dbConn.Collection("Followers")
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

func GetDocFromDBUsingEmail(dbCollection *mongo.Collection, usersEmail string) (*mongo.SingleResult, error) {
	result := dbCollection.FindOne(context.TODO(), bson.M{"email": usersEmail})
	if result.Err() == mongo.ErrNoDocuments {
		return result, result.Err()
	}

	return result, nil

}

func GetDocFromDBUsingUUID(dbCollection *mongo.Collection, usersUUID string) (*mongo.SingleResult, error) {
	result := dbCollection.FindOne(context.TODO(), bson.M{"uuid": usersUUID})
	if result.Err() == mongo.ErrNoDocuments {
		return result, result.Err()
	}

	return result, nil

}

func DeleteUser(dbCollection *mongo.Collection, usersUUID string) error {
	_, err := dbCollection.DeleteOne(context.TODO(), bson.M{"uuid": usersUUID})
	if err != nil {
		return err
	}

	return nil
}

func GetAllMatchingDocuments(dbCollection *mongo.Collection, userUUID string) ([]bson.M, error) {
	cursor, err := dbCollection.Find(context.TODO(), bson.M{"user_uuid": userUUID})
	if err != nil {
		return nil, err
	}

	var result []bson.M

	err2 := cursor.All(context.TODO(), &result)
	if err2 != nil {
		return nil, err2
	}

	return result, nil
}

func GetTweets(dbCollection *mongo.Collection, userUUID string) ([]models.Tweet, error) {
	var tweets []models.Tweet

	// ignore the following fields and return the rest from the db
	var projection = bson.M{
		"user_uuid":  0,
		"email":      0,
		"created_at": 0,
		"updated_at": 0,
	}

	cursor, err := dbCollection.Find(context.TODO(), bson.M{"user_uuid": userUUID}, options.Find().SetProjection(projection))
	if err != nil {
		return tweets, err
	}

	// If there are no tweets then return an empty string
	if cursor.RemainingBatchLength() == 0 {
		return []models.Tweet{}, nil
	}

	// Decode all the tweets from the db to the Tweets struct
	err2 := cursor.All(context.TODO(), &tweets)
	if err2 != nil {
		return tweets, err2
	}

	return tweets, nil
}

func FollowerFollowingCombinationExists(dbCollection *mongo.Collection, followerAccountID string, followingAccountID string) bool {

	result := dbCollection.FindOne(context.TODO(), bson.M{
		"follower_account_id":  followerAccountID,
		"following_account_id": followingAccountID})

	return result.Err() != mongo.ErrNoDocuments
}

func UpdateFollowingCount(dbCollection *mongo.Collection, accountID string) error {

	result := dbCollection.FindOneAndUpdate(context.TODO(),
		bson.M{"account_id": accountID},
		bson.M{"$inc": bson.M{"metrics.following_count": 1}})

	if result.Err() == mongo.ErrNoDocuments {
		return result.Err()
	}

	return nil
}

func UpdateFollowersCount(dbCollection *mongo.Collection, accountID string) error {

	result := dbCollection.FindOneAndUpdate(context.TODO(),
		bson.M{"account_id": accountID},
		bson.M{"$inc": bson.M{"metrics.followers_count": 1}})

	if result.Err() == mongo.ErrNoDocuments {
		return result.Err()
	}

	return nil
}
