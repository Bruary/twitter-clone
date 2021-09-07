package twitter

import (
	"context"
	"time"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Saves a tweet to the db with all required information
func (*twitter) CreateTweet(c *fiber.Ctx, req models.CreateTweetRequest) *models.BaseResponse {

	// Request validation
	tokenValueEmpty := validate.IsStringEmpty(req.Token)
	if tokenValueEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field token is missing, or empty.",
		}
	}

	tweetValueEmpty := validate.IsStringEmpty(req.Tweet)
	if tweetValueEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field tweet is missing, or empty.",
		}
	}

	// Validate token
	validToken := validate.IsTokenValid(req.Token)
	if !validToken {

		c.Status(fiber.StatusUnauthorized)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "INVALID_TOKEN",
			Msg:          "Invalid token.",
		}
	}

	// Extract the JWT claims
	tokenClaims := validate.GetJWTclaims(req.Token)

	// Check if user exists
	userDoc, err1_5 := db.GetDocFromDBUsingUUID(db.UsersCol, tokenClaims.User_UUID)
	if err1_5 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Failed while finding user in db.",
		}
	}

	var user models.UserInfo

	err1_55 := userDoc.Decode(&user)
	if err1_55 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "USER_DOES_NOT_EXIST",
			Msg:          "Invalid email address.",
		}
	}

	// Insert all the info that are required to be saved with the tweet
	var tweetInfo = models.TweetDB{
		User_UUID:  user.UUID,
		Account_ID: tokenClaims.Account_ID,
		Tweet_UUID: uuid.NewV4().String(),
		Email:      user.Email,
		Tweet:      req.Tweet,
		Metrics: models.TweetMetrics{
			Retweets_count:   0,
			Likes_count:      0,
			Comments_count:   0,
			Characters_count: len(req.Tweet),
		},
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	err2 := db.InsertDocumentToDB(db.TweetsCol, tweetInfo)
	if err2 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Inserting tweet to the db failed.",
		}
	}

	//update the tweet count on the users document in the db
	updateResult := db.UsersCol.FindOneAndUpdate(context.TODO(), bson.M{"uuid": user.UUID}, bson.M{"$set": bson.M{"metrics.total_tweets_count": user.Metrics.Total_tweets_count + 1}})
	if updateResult.Err() == mongo.ErrNoDocuments {
		return &models.BaseResponse{
			Success:      false,
			ResponseType: "TWEET_SAVED",
			Msg:          updateResult.Err().Error(),
		}
	}

	return &models.BaseResponse{
		Success:      true,
		ResponseType: "TWEET_SAVED",
		Msg:          "Tweet saved to db successfully.",
	}
}
