package twitter

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/validate"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Creates a new user and adds it to the db
func (s *twitter) CreateUser(c *fiber.Ctx, req models.CreateUserRequest) *models.BaseResponse {

	// Validate request
	firstNameEmpty := validate.IsStringEmpty(req.FirstName)
	if firstNameEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field firstname is missing, or empty.",
		}
	}

	lastNameEmpty := validate.IsStringEmpty(req.LastName)
	if lastNameEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field lastname is missing, or empty.",
		}
	}

	emailEmpty := validate.IsStringEmpty(req.Email)
	if emailEmpty {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_MISSING",
			Msg:          "Field email is missing, or empty.",
		}
	}

	passwordValid := validate.IsPasswordLengthCorrect(req.Password)
	if !passwordValid {
		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_ERROR",
			Msg:          "Password should atleast have 8 characters.",
		}
	}

	ageValid := validate.IsAge12AndAbove(req.Age)
	if !ageValid {

		c.Status(fiber.ErrBadRequest.Code)

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_ERROR",
			Msg:          "Age should be 12 and above years old to create an account.",
		}
	}

	// check is user already exist in db
	doesUserExist, err10 := db.UserAlreadyExists(db.UsersCol, req.Email)
	if err10 == nil && doesUserExist {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "USER_ALREADY_EXISTS",
			Msg:          "User's email already exists.",
		}
	}

	passwordHashedAndSalted, err10_5 := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err10_5 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Hashing password failed in SignUp endpoint.",
		}
	}

	// Create a random account ID
	a := make([]byte, 5)

	_, err20_1 := rand.Read(a)
	if err20_1 != nil {

		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Creating account ID failed.",
		}
	}

	accountID := strings.ToUpper("#" + hex.EncodeToString(a))

	// Fill in the user details
	userInfo := &models.UserInfo{
		UUID:       uuid.NewV4().String(),
		Account_ID: accountID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Age:        req.Age,
		Email:      req.Email,
		Password:   string(passwordHashedAndSalted),
		Metrics: models.UserMetrics{
			Followers_count:      0,
			Following_count:      0,
			Total_tweets_count:   0,
			Total_retweets_count: 0,
			Total_likes_count:    0,
		},
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	// call the InsertDocumentToDB func to add a new user to the db collection 'Users'
	err1 := db.InsertDocumentToDB(db.UsersCol, userInfo)
	if err1 != nil {
		return &models.BaseResponse{
			Success:      false,
			ResponseType: "UNKNOWN_ERROR",
			Msg:          "Inserting new user to the db failed.",
		}
	}

	return &models.BaseResponse{
		Success:      true,
		ResponseType: "NEW_USER_CREATED",
		Msg:          "New user was added and saved to the db.",
	}
}
