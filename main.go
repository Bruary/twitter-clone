// main functionality:
// 1) Create new user
// 2) Delete account
// 3) Post a tweet
// 4) Make account/tweets private or public
// 5) News feed

package main

import (
	"encoding/json"
	"time"

	"github.com/Bruary/twitter-clone/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseResponse struct {
	ResponseType string
	Success      bool
	Msg          string
}

type CreateUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

type DeleteUserRequest struct {
	Email string `json:"email"`
}

type MakeATweetRequest struct {
	Email string `json:"email"`
	Tweet string `json:"tweet"`
}

type TweetRequiredFields struct {
	Email         string
	Tweet         string
	Metrics       TweetMetrics
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

type TweetMetrics struct {
	Retweets_count   int
	Likes_count      int
	Comments_count   int
	Characters_count int
}

var UsersCol *mongo.Collection
var TweetsCol *mongo.Collection

func main() {

	dbConn := db.GetDBConn()

	UsersCol = dbConn.Collection("Users")
	TweetsCol = dbConn.Collection("Tweets")

	app := fiber.New()

	app.Post("/createUser", CreateUser)
	app.Delete("/deleteUser", DeleteUser)
	app.Post("/makeATweet", MakeATweet)

	app.Listen(":3000")
}

// Creates a new user and adds it to the db
func CreateUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/json")

	var userInfo CreateUserRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &userInfo)
	if err != nil {
		c.SendString("Unmashaling failed in CreateUser endpoint.")
		return err
	}

	// check is user already exist in db
	doesUserExist, err10 := db.UserAlreadyExists(UsersCol, userInfo.Email)
	if err10 == nil && doesUserExist {
		resp.ResponseType = "USER_ALREADY_EXISTS"
		resp.Success = false
		resp.Msg = "User's email already exists."

		output01, err20 := json.Marshal(&resp)
		if err20 != nil {
			c.SendString("Mashaling failed in CreateUser endpoint.")
			return err20
		}

		c.Context().Response.SetBody(output01)

		return nil
	}

	// call the InsertDocumentToDB func to add a new user to the db collection 'Users'
	err1 := db.InsertDocumentToDB(UsersCol, userInfo)
	if err1 != nil {
		c.SendString("Inserting new user to the db failed.")
		return err1
	}

	resp.ResponseType = "NEW_USER_CREATED"
	resp.Success = true
	resp.Msg = "New user was added and saved to the db."

	output, err2 := json.Marshal(&resp)
	if err2 != nil {
		c.SendString("Mashaling failed in CreateUser endpoint.")
		return err
	}

	c.Context().Response.SetBody(output)
	return nil
}

// Delete user from the user using the email
func DeleteUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/jsons")

	var userEmail DeleteUserRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &userEmail)
	if err != nil {
		c.SendString("Unmarshaling failed in DeleteUser endpoint.")
		return err
	}

	err1 := db.DeleteUser(UsersCol, userEmail.Email)
	if err1 != nil {

		resp.ResponseType = "UNKNOWN_ERROR"
		resp.Success = false
		resp.Msg = "Failed to delete user from DB."

	} else {

		resp.ResponseType = "USER_DELETED"
		resp.Success = true
		resp.Msg = "User has been successfuly deleted from the db."

	}

	result, err2 := json.Marshal(resp)
	if err2 != nil {
		c.SendString("Marshaling failed in DeleteUser endpoint.")
		return err2
	}

	c.Context().SetBody(result)
	return nil
}

// Saves a tweet to the db with all required information
func MakeATweet(c *fiber.Ctx) error {

	c.Context().SetContentType("application/jsons")

	var tweetRequest MakeATweetRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &tweetRequest)
	if err != nil {
		c.SendString("Marshaling failed in MakeATweet endpoint.")
		return err
	}

	// Insert all the info that are required to be saved with the tweet
	var tweetInfo = TweetRequiredFields{
		Email: tweetRequest.Email,
		Tweet: tweetRequest.Tweet,
		Metrics: TweetMetrics{
			Retweets_count:   0,
			Likes_count:      0,
			Comments_count:   0,
			Characters_count: len(tweetRequest.Tweet),
		},
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}

	err2 := db.InsertDocumentToDB(TweetsCol, tweetInfo)
	if err2 != nil {
		c.SendString("Inserting tweet to the db failed.")
		return err2
	}

	resp.Success = true
	resp.ResponseType = "TWEET_SAVED"
	resp.Msg = "Tweet saved to db successfully."

	result, err3 := json.Marshal(resp)
	if err3 != nil {
		c.SendString("Marshaling failed in MakeATweet endpoint.")
		return err3
	}

	c.Context().SetBody(result)

	return nil
}
