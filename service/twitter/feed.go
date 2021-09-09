package twitter

import (
	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
)

func (*twitterClone) Feed(c *fiber.Ctx, req models.BaseRequest) *models.FeedResponse {

	// STEPS
	tokenValueEmpty := validate.IsStringEmpty(req.Token)
	if tokenValueEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.FeedResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "FIELD_MISSING",
				Msg:          "Field token is missing, or empty.",
			}}
	}

	// Validate token
	validToken := validate.IsTokenValid(req.Token)
	if !validToken {

		c.Status(fiber.StatusUnauthorized)

		return &models.FeedResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "INVALID_TOKEN",
				Msg:          "Invalid token.",
			}}
	}

	// Extract the JWT claims
	tokenClaims := validate.GetJWTclaims(req.Token)

	// get all following_account_ids
	followingAccountIDs, err := db.GetAllFollowingAccountIDs(db.FollowersCol, tokenClaims.Account_ID)
	if err != nil {
		return &models.FeedResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Getting following IDs from db failed.",
			}}
	}

	// Get all tweets for each of the following accounts
	tweets, err2 := db.GetTweetsForAListOfAccounts(db.TweetsCol, followingAccountIDs)
	if err2 != nil {
		return &models.FeedResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Failed to get tweets for feed.",
			}}
	}

	return &models.FeedResponse{
		BaseResponse: models.BaseResponse{
			Success: true,
		},
		Tweets: tweets,
	}

}
