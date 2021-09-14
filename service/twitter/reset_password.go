package twitter

import (
	"net/smtp"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jordan-wright/email"
	"go.mongodb.org/mongo-driver/mongo"
)

func (*twitterClone) ResetPassword(c *fiber.Ctx, req models.ResetPasswordRequest) *models.BaseResponse {

	// Check if user exists in the db
	result, err := db.GetDocFromDBUsingEmail(db.UsersCol, req.Email)
	if err != nil {

		// if user does not exist
		if err == mongo.ErrNoDocuments {
			return &models.BaseResponse{
				Success:      true,
				ResponseType: "USER_DOES_NOT_EXIST",
				Msg:          err.Error(),
			}
		}

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          err.Error(),
		}
	}

	var user models.UserInfo

	// decode the single result into the user variable
	if err2 := result.Decode(&user); err2 != nil {
		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          err2.Error(),
		}
	}

	// create JWT to be send in the email
	token, err3 := CreateJWT(user.UUID, user.Account_ID, 15)
	if err3 != nil {
		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          err3.Error(),
		}
	}

	// draft the email
	e := email.Email{
		To:      []string{user.Email},
		From:    "Twitter-clone <bruary99@gmail.com>",
		Subject: "Reset Password",
		Text:    []byte("Please click on the below link to reset your password: \n" + "http://localhost:4000/api/v1/auth/resetPassword/newPassword?token=" + token),
	}

	// send the email
	err2 := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "bruary99@gmail.com", "Bakir12345", "smtp.gmail.com"))
	if err2 != nil {
		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          err2.Error(),
		}
	}

	return &models.BaseResponse{
		Success: true,
	}
}
