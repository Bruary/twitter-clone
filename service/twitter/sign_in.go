package twitter

import (
	"time"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("White_Yasmin")
var UsersCol *mongo.Collection

// validate user and then sign in if creds are correct (send back token)
func SignIn(c *fiber.Ctx, req models.SignInRequest) *models.SignInResponse {

	c.Context().SetContentType("application/jsons")

	// Validate request
	emailValueEmpty := validate.IsStringEmpty(req.Email)
	if emailValueEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "FIELD_MISSING",
				Msg:          "Field email is missing, or empty.",
			},
		}
	}

	passwordEmpty := validate.IsStringEmpty(req.Password)
	if passwordEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "FIELD_MISSING",
				Msg:          "Field password is missing, or empty.",
			},
		}
	}

	// check if email exists in the db
	doesUserExist, err2 := db.UserAlreadyExists(UsersCol, req.Email)
	if err2 != nil {

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Failed while searching for the User in SignIn endpoint.",
			},
		}
	}

	if !doesUserExist {

		c.Status(fiber.StatusUnauthorized)

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "USER_DOES_NOT_EXIST",
				Msg:          "User does not exist.",
			},
		}
	}

	// Get the user document from the db to check the password later on
	userDocument, err2_5 := db.GetDocFromDBUsingEmail(UsersCol, req.Email)
	if err2_5 != nil {

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Finding the user document in the DB failed in SignIn endpoint.",
			},
		}
	}

	// Check if the request password matches the one stored in the DB
	var userDocumentDecoded models.UserInfo

	err2_6 := userDocument.Decode(&userDocumentDecoded)
	if err2_6 != nil {

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Decoding failed in SignIn endpoint.",
			},
		}
	}

	// Check if the received password with hashing matches the one saved in the db
	isPasswordCorrect := DoPasswordsMatch([]byte(req.Password), userDocumentDecoded.Password)
	if !isPasswordCorrect {

		c.Status(fiber.StatusUnauthorized)

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "INVALID_CREDENTIALS",
				Msg:          "Invalid email or password.",
			},
		}
	}

	// If password matches then do the following

	// create the claims that will be used in the JWT token
	claims := &models.Claims{
		User_UUID:  userDocumentDecoded.UUID,
		Account_ID: userDocumentDecoded.Account_ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		},
	}

	// declaring the token with the method used for signing along with the claims§
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err4 := token.SignedString(jwtKey)
	if err4 != nil {

		return &models.SignInResponse{
			BaseResponse: models.BaseResponse{
				Success:      false,
				ResponseType: "UNKNOWN_ERROR",
				Msg:          "Creating the token failed in SignIn endpoint.",
			},
		}
	}

	// filing in the final response with the generated token string
	return &models.SignInResponse{
		BaseResponse: models.BaseResponse{
			Success: true,
		},
		Token: tokenString,
	}
}

func DoPasswordsMatch(password []byte, savedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), password)
	return err == nil
}
