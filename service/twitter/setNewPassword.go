package twitter

import (
	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (*twitterClone) SetNewPassword(c *fiber.Ctx, req models.SetNewPasswordRequest) *models.BaseResponse {

	// Validate token
	isTokenValid := validate.IsTokenValid(req.Token)
	if !isTokenValid {
		c.Status(fiber.StatusForbidden)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "INVALID_TOKEN",
			Msg:          "Invalid Token.",
		}
	}

	// Extract claims
	claims := validate.GetJWTclaims(req.Token)

	// // Get user from DB
	// result, err := db.GetDocFromDBUsingUUID(db.UsersCol, claims.User_UUID)
	// if err != nil {

	// 	// if user does not exist
	// 	if err == mongo.ErrNoDocuments {

	// 		c.Status(fiber.StatusNotFound)

	// 		return &models.BaseResponse{
	// 			Success:      true,
	// 			ResponseType: "USER_DOES_NOT_EXIST",
	// 			Msg:          err.Error(),
	// 		}
	// 	}

	// 	c.Status(fiber.StatusBadRequest)

	// 	return &models.BaseResponse{
	// 		Success:      false,
	// 		ResponseType: "UNKNOWN_ERROR",
	// 		Msg:          err.Error(),
	// 	}
	// }

	// var user models.UserInfo

	// // decode mongo doc into the user struct
	// err2 := result.Decode(&user)
	// if err2 != nil {

	// 	c.Status(fiber.StatusBadRequest)

	// 	return &models.BaseResponse{
	// 		Success:      false,
	// 		ResponseType: "UNKNOWN_ERROR",
	// 		Msg:          err2.Error(),
	// 	}
	// }

	passwordHashedAndSalted, err10_5 := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err10_5 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Hashing password failed in SignUp endpoint.",
		}
	}

	err3 := db.UpdateUsersPassword(db.UsersCol, claims.User_UUID, string(passwordHashedAndSalted))
	if err3 != nil {

		c.Status(fiber.StatusBadRequest)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          err3.Error(),
		}
	}

	return &models.BaseResponse{
		Success: true,
	}
}
