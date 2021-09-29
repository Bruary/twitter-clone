package twitter

import (
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
)

func (*twitterClone) NewPassword(c *fiber.Ctx) *models.BaseResponse {

	token := c.Query("token")

	if token == "" {
		return &models.BaseResponse{
			Success:      true,
			ResponseType: "INVALID_TOKEN",
			Msg:          "Token is missing.",
		}
	}

	if isTokenValid := validate.IsTokenValid(token); !isTokenValid {
		return &models.BaseResponse{
			Success:      true,
			ResponseType: "INVALID_TOKEN",
			Msg:          "Invalid token.",
		}
	}

	// 2) redirect to the reset password page on frontend
	// 3) send the JWT to the frontend
	// 4) get the new password from the frontend along with the JWT and update the password

	err := c.Redirect("http://localhost:3000/api/v1/auth/resetPassword/newPassword?token:" + token)
	if err != nil {
		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Failed to redirect to a new page.",
		}
	}

	return &models.BaseResponse{
		Success: true,
	}
}
