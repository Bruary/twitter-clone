package twitter

import (
	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
)

// Delete user from the user using the email
func (*twitter) DeleteUser(c *fiber.Ctx, req models.DeleteUserRequest) *models.BaseResponse {

	c.Context().SetContentType("application/jsons")

	tokenEmptyValue := validate.IsStringEmpty(req.Token)
	if tokenEmptyValue {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field token is missing, or empty.",
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

	err1 := db.DeleteUser(db.UsersCol, tokenClaims.User_UUID)
	if err1 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Failed to delete user from DB.",
		}

	} else {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "USER_DELETED",
			Msg:          "User has been successfuly deleted from the db.",
		}

	}
}
