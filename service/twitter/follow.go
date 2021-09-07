package twitter

import (
	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func (*twitter) Follow(c *fiber.Ctx, req models.FollowRequest) *models.BaseResponse {
	c.Context().SetContentType("applications/json")

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

	followingAccountIDEmpty := validate.IsStringEmpty(req.Following_Account_ID)
	if followingAccountIDEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field following_account_id is missing, or empty.",
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

	// Check if this follower-follower relationship exist in the db
	followerFollowingCombExits := db.FollowerFollowingCombinationExists(db.FollowersCol, tokenClaims.Account_ID, req.Following_Account_ID)
	if !followerFollowingCombExits {

		followerData := &models.Followers{
			ID:                   uuid.NewV4().String(),
			Follower_Account_ID:  tokenClaims.Account_ID,
			Following_Account_ID: req.Following_Account_ID,
		}

		err2 := db.InsertDocumentToDB(db.FollowersCol, followerData)
		if err2 != nil {

			return &models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Inserting tweet to the db failed.",
			}
		}

		// increment the "following" field
		err3 := db.UpdateFollowingCount(db.UsersCol, tokenClaims.Account_ID)
		if err3 != nil {

			return &models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Updating following count failed.",
			}
		}

		// increment the "followers" field
		err4 := db.UpdateFollowersCount(db.UsersCol, req.Following_Account_ID)
		if err4 != nil {

			return &models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Updating followers count failed.",
			}
		}

	}

	return &models.BaseResponse{
		Success: true,
	}
}
