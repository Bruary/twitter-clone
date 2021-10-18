package twitter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func (*twitterClone) GetTweets(c *fiber.Ctx, req models.BaseRequest) *models.GetTweetsResponse {

	// Request validation
	tokenValueEmpty := validate.IsStringEmpty(req.Token)
	if tokenValueEmpty {
		c.Status(fiber.ErrBadRequest.Code)

		return &models.GetTweetsResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "FIELD_MISSING",
				Msg:          "Field token is missing, or empty.",
			},
		}
	}

	// Validate token
	validToken := validate.IsTokenValid(req.Token)
	if !validToken {

		c.Status(fiber.StatusUnauthorized)

		return &models.GetTweetsResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "INVALID_TOKEN",
				Msg:          "Invalid token.",
			},
		}
	}

	// Extract the JWT claims
	tokenClaims := validate.GetJWTclaims(req.Token)

	// create a unique UUID to save it on redis and then save it
	cacheUUID := "GET_TWEETS:" + tokenClaims.User_UUID

	// try to get it from cache
	result, err := db.RedisClient.Get(context.Background(), cacheUUID).Result()
	if err == redis.Nil {
		fmt.Println(cacheUUID + " cache key was not found.")
	} else if err != nil {
		fmt.Println("Cache get failed", err)
	} else {

		var cachedTweets []models.Tweet

		unmarshalErr := json.Unmarshal([]byte(result), &cachedTweets)
		if unmarshalErr != nil {
			fmt.Println("Marshaling failed in getting cahce.")
		}

		fmt.Println("from cache")

		return &models.GetTweetsResponse{
			BaseResponse: models.BaseResponse{
				Success: true,
			},
			Tweets: cachedTweets,
		}
	}

	tweets, err12 := db.GetTweetsUsingUUID(db.TweetsCol, tokenClaims.User_UUID)
	if err12 != nil {

		return &models.GetTweetsResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Failed to get tweets from db.",
			},
		}
	}

	cachedTweets, _ := json.Marshal(tweets)

	// save response on cache
	cacheErr := db.RedisClient.Set(context.Background(), cacheUUID, cachedTweets, 0).Err()
	if cacheErr != nil {
		fmt.Println("Writing cache failed: ", cacheErr)
	}

	return &models.GetTweetsResponse{
		BaseResponse: models.BaseResponse{
			Success: true,
		},
		Tweets: tweets,
	}
}
