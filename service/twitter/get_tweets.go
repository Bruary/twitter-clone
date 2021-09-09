package twitter

import (
	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
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

	return &models.GetTweetsResponse{
		BaseResponse: models.BaseResponse{
			Success: true,
		},
		Tweets: tweets,
	}
}
