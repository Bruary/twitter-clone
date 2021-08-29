// main functionality:
// 1) Create new user
// 2) Delete account
// 3) Post a tweet
// 4) Make account/tweets private or public
// 5) News feed

package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"crypto/rand"

	"encoding/hex"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/models"
	"github.com/Bruary/twitter-clone/validate"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	bcrypt "golang.org/x/crypto/bcrypt"
)

var UsersCol *mongo.Collection
var TweetsCol *mongo.Collection
var FollowersCol *mongo.Collection

var jwtKey = []byte("White_Yasmin")

func main() {

	dbConn := db.GetDBConn()

	UsersCol = dbConn.Collection("Users")
	TweetsCol = dbConn.Collection("Tweets")
	FollowersCol = dbConn.Collection("Followers")

	app := fiber.New()

	app.Use(cors.New())

	app.Post("/signin", SignIn)
	app.Post("/createUser", CreateUser)
	app.Delete("/deleteUser", DeleteUser)
	app.Post("/makeATweet", MakeATweet)
	app.Post("/getTweets", GetTweets)
	app.Post("/follow", Follow)

	app.Listen(":4000")
}

// Creates a new user and adds it to the db
func CreateUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/json")

	var userReceivedInfo models.CreateUserRequest
	var resp models.BaseResponse

	err := UnmarshalRequest(&userReceivedInfo, c)
	if err != nil {
		return err
	}

	// Validate request
	firstNameEmpty := validate.IsStringEmpty(userReceivedInfo.FirstName)
	if firstNameEmpty {
		resp = SetMissingFieldResponse("firstname")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	lastNameEmpty := validate.IsStringEmpty(userReceivedInfo.LastName)
	if lastNameEmpty {
		resp = SetMissingFieldResponse("lastname")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	emailEmpty := validate.IsStringEmpty(userReceivedInfo.Email)
	if emailEmpty {
		resp = SetMissingFieldResponse("email")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	passwordValid := validate.IsPasswordLengthCorrect(userReceivedInfo.Password)
	if !passwordValid {
		resp = SetCriteriaErrorResponse("password")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	ageValid := validate.IsAge12AndAbove(userReceivedInfo.Age)
	if !ageValid {
		resp = SetCriteriaErrorResponse("age")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	// check is user already exist in db
	doesUserExist, err10 := db.UserAlreadyExists(UsersCol, userReceivedInfo.Email)
	if err10 == nil && doesUserExist {
		resp.ResponseType = "USER_ALREADY_EXISTS"
		resp.Success = false
		resp.Msg = "User's email already exists."

		err20 := MarshalResponseAndSetBody(resp, c)
		if err20 != nil {
			return err20
		}

		return nil
	}

	passwordHashedAndSalted, err10_5 := bcrypt.GenerateFromPassword([]byte(userReceivedInfo.Password), bcrypt.MinCost)
	if err10_5 != nil {
		c.SendString("Hashing password failed in SignUp endpoint.")
		return err10_5
	}

	// Create a random account ID
	a := make([]byte, 5)

	_, err20_1 := rand.Read(a)
	if err20_1 != nil {
		c.SendString("Creating account ID failed.")
		return err20_1
	}

	accountID := strings.ToUpper("#" + hex.EncodeToString(a))

	// Fill in the user details
	userInfo := &models.UserInfo{
		UUID:       uuid.NewV4().String(),
		Account_ID: accountID,
		FirstName:  userReceivedInfo.FirstName,
		LastName:   userReceivedInfo.LastName,
		Age:        userReceivedInfo.Age,
		Email:      userReceivedInfo.Email,
		Password:   string(passwordHashedAndSalted),
		Metrics: models.UserMetrics{
			Followers_count:      0,
			Total_tweets_count:   0,
			Total_retweets_count: 0,
			Total_likes_count:    0,
		},
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	// call the InsertDocumentToDB func to add a new user to the db collection 'Users'
	err1 := db.InsertDocumentToDB(UsersCol, userInfo)
	if err1 != nil {
		c.SendString("Inserting new user to the db failed.")
		return err1
	}

	resp.ResponseType = "NEW_USER_CREATED"
	resp.Success = true
	resp.Msg = "New user was added and saved to the db."

	err2 := MarshalResponseAndSetBody(resp, c)
	if err2 != nil {
		return err2
	}

	return nil
}

// Delete user from the user using the email
func DeleteUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/jsons")

	var req models.DeleteUserRequest
	var resp models.BaseResponse

	err := UnmarshalRequest(&req, c)
	if err != nil {
		return err
	}

	tokenEmptyValue := validate.IsStringEmpty(req.Token)
	if tokenEmptyValue {
		resp = SetMissingFieldResponse("token")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	// Validate token
	validToken := IsTokenValid(req.Token)
	if !validToken {
		resp.Success = false
		resp.ResponseType = "INVALID_TOKEN"
		resp.Msg = "Invalid token."

		c.Status(fiber.StatusUnauthorized)
		err2 := MarshalResponseAndSetBody(resp, c)
		if err2 != nil {
			return err2
		}

		return nil
	}

	// Extract the JWT claims
	tokenClaims := GetJWTclaims(req.Token)

	err1 := db.DeleteUser(UsersCol, tokenClaims.User_UUID)
	if err1 != nil {

		resp.ResponseType = "UNKNOWN_ERROR"
		resp.Success = false
		resp.Msg = "Failed to delete user from DB."

	} else {

		resp.ResponseType = "USER_DELETED"
		resp.Success = true
		resp.Msg = "User has been successfuly deleted from the db."

	}

	err2 := MarshalResponseAndSetBody(resp, c)
	if err2 != nil {
		return err2
	}

	return nil
}

// validate user and then sign in if creds are correct (send back token)
func SignIn(c *fiber.Ctx) error {

	c.Context().SetContentType("application/jsons")

	var req models.SignInRequest
	var baseResp models.BaseResponse

	err := UnmarshalRequest(&req, c)
	if err != nil {
		return err
	}

	// Validate request
	emailValueEmpty := validate.IsStringEmpty(req.Email)
	if emailValueEmpty {
		baseResp = SetMissingFieldResponse("email")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(baseResp, c)
		if err != nil {
			return err
		}
		return nil
	}

	passwordEmpty := validate.IsStringEmpty(req.Password)
	if passwordEmpty {
		baseResp = SetMissingFieldResponse("password")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(baseResp, c)
		if err != nil {
			return err
		}
		return nil
	}

	// check if email exists in the db
	doesUserExist, err2 := db.UserAlreadyExists(UsersCol, req.Email)
	if err2 != nil {
		c.SendString("Failed while searching for the User in SignIn endpoint.")
		return err2
	}

	if !doesUserExist {
		baseResp.Success = false
		baseResp.ResponseType = "USER_DOES_NOT_EXIST"
		baseResp.Msg = "Invalid email or password."

		err3 := MarshalResponseAndSetBody(baseResp, c)
		if err3 != nil {
			return err3
		}

		return nil
	}

	// Get the user document from the db to check the password later on
	userDocument, err2_5 := db.GetDocFromDBUsingEmail(UsersCol, req.Email)
	if err2_5 != nil {
		c.SendString("Finding the user document in the DB failed in SignIn endpoint.")
		return err2_5
	}

	// Check if the request password matches the one stored in the DB
	var userDocumentDecoded models.UserInfo

	err2_6 := userDocument.Decode(&userDocumentDecoded)
	if err2_6 != nil {
		c.SendString("Decoding failed in SignIn endpoint.")
		return err2_6
	}

	// Check if the received password with hashing matches the one saved in the db
	isPasswordCorrect := DoPasswordsMatch([]byte(req.Password), userDocumentDecoded.Password)

	if !isPasswordCorrect {
		baseResp.Success = false
		baseResp.ResponseType = "INVALID_CREDENTIALS"
		baseResp.Msg = "Invalid email or password."

		err3_2 := MarshalResponseAndSetBody(baseResp, c)
		if err3_2 != nil {
			return err3_2
		}

		return nil
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
		c.SendString("Creating the token failed in SignIn endpoint.")
		return err4
	}

	// filing in the final response with the generated token string
	resp := &models.SignInResponse{
		Success: true,
		Token:   tokenString,
	}

	err5 := MarshalResponseAndSetBody(resp, c)
	if err5 != nil {
		return err5
	}

	return nil
}

// Saves a tweet to the db with all required information
func MakeATweet(c *fiber.Ctx) error {

	c.Context().SetContentType("application/jsons")

	var tweetRequest models.MakeATweetRequest
	var resp models.BaseResponse

	err := UnmarshalRequest(&tweetRequest, c)
	if err != nil {
		return err
	}

	// Request validation
	tokenValueEmpty := validate.IsStringEmpty(tweetRequest.Token)
	if tokenValueEmpty {
		resp = SetMissingFieldResponse("token")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	tweetValueEmpty := validate.IsStringEmpty(tweetRequest.Tweet)
	if tweetValueEmpty {
		resp = SetMissingFieldResponse("tweet")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	// Validate token
	validToken := IsTokenValid(tweetRequest.Token)
	if !validToken {
		resp.Success = false
		resp.ResponseType = "INVALID_TOKEN"
		resp.Msg = "Invalid token."

		c.Status(fiber.StatusUnauthorized)

		err1_4 := MarshalResponseAndSetBody(resp, c)
		if err1_4 != nil {
			return err1_4
		}

		return nil
	}

	// Extract the JWT claims
	tokenClaims := GetJWTclaims(tweetRequest.Token)

	// Check if user exists
	userDoc, err1_5 := db.GetDocFromDBUsingUUID(UsersCol, tokenClaims.User_UUID)
	if err1_5 != nil {
		c.SendString("Failed while finding user in db.")
		return err1_5
	}

	var user models.UserInfo

	err1_55 := userDoc.Decode(&user)
	if err1_55 != nil {
		resp.Success = false
		resp.ResponseType = "USER_DOES_NOT_EXIST"
		resp.Msg = "Invalid email address."

		err1_6 := MarshalResponseAndSetBody(resp, c)
		if err1_6 != nil {
			return err1_6
		}

		return nil
	}

	// Insert all the info that are required to be saved with the tweet
	var tweetInfo = models.TweetDB{
		User_UUID:  user.UUID,
		Tweet_UUID: uuid.NewV4().String(),
		Email:      user.Email,
		Tweet:      tweetRequest.Tweet,
		Metrics: models.TweetMetrics{
			Retweets_count:   0,
			Likes_count:      0,
			Comments_count:   0,
			Characters_count: len(tweetRequest.Tweet),
		},
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	err2 := db.InsertDocumentToDB(TweetsCol, tweetInfo)
	if err2 != nil {
		c.SendString("Inserting tweet to the db failed.")
		return err2
	}

	resp.Success = true
	resp.ResponseType = "TWEET_SAVED"
	resp.Msg = "Tweet saved to db successfully."

	err3 := MarshalResponseAndSetBody(resp, c)
	if err3 != nil {
		return err3
	}

	//update the tweet count on the users document in the db
	updateResult := UsersCol.FindOneAndUpdate(context.TODO(), bson.M{"uuid": user.UUID}, bson.M{"$set": bson.M{"metrics.total_tweets_count": user.Metrics.Total_tweets_count + 1}})
	if updateResult.Err() == mongo.ErrNoDocuments {
		c.SendString("Failed to update number of tweets for a user in MakeATweet endpoint.")
		return updateResult.Err()
	}

	return nil
}

func GetTweets(c *fiber.Ctx) error {

	c.Context().SetContentType("applications/json")

	var req models.GetTweetsRequest
	var baseResp models.BaseResponse

	err := UnmarshalRequest(&req, c)
	if err != nil {
		return err
	}

	// Request validation
	tokenValueEmpty := validate.IsStringEmpty(req.Token)
	if tokenValueEmpty {
		baseResp = SetMissingFieldResponse("token")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(baseResp, c)
		if err != nil {
			return err
		}
		return nil
	}

	// Validate token
	validToken := IsTokenValid(req.Token)
	if !validToken {
		baseResp.Success = false
		baseResp.ResponseType = "INVALID_TOKEN"
		baseResp.Msg = "Invalid token."

		c.Status(fiber.StatusUnauthorized)
		err2 := MarshalResponseAndSetBody(baseResp, c)
		if err2 != nil {
			return err2
		}

		return nil
	}

	// Extract the JWT claims
	tokenClaims := GetJWTclaims(req.Token)

	tweets, err12 := db.GetTweets(TweetsCol, tokenClaims.User_UUID)
	if err12 != nil {
		return err12
	}

	resp := &models.GetTweetsResponse{
		Success: true,
		Tweets:  tweets,
	}

	err122 := MarshalResponseAndSetBody(resp, c)
	if err122 != nil {
		return err122
	}

	return nil
}

func Follow(c *fiber.Ctx) error {

	c.Context().SetContentType("applications/json")

	var req models.FollowRequest
	var resp models.BaseResponse

	err := UnmarshalRequest(&req, c)
	if err != nil {
		return err
	}

	// Request validation
	tokenValueEmpty := validate.IsStringEmpty(req.Token)
	if tokenValueEmpty {
		resp = SetMissingFieldResponse("token")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	followingAccountIDEmpty := validate.IsStringEmpty(req.Following_Account_ID)
	if followingAccountIDEmpty {
		resp = SetMissingFieldResponse("following_account_id")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	// Validate token
	validToken := IsTokenValid(req.Token)
	if !validToken {
		resp.Success = false
		resp.ResponseType = "INVALID_TOKEN"
		resp.Msg = "Invalid token."

		c.Status(fiber.StatusUnauthorized)
		err2 := MarshalResponseAndSetBody(resp, c)
		if err2 != nil {
			return err2
		}

		return nil
	}

	// Extract the JWT claims
	tokenClaims := GetJWTclaims(req.Token)

	followerData := &models.Followers{
		ID:                   uuid.NewV4().String(),
		Follower_Account_ID:  tokenClaims.Account_ID,
		Following_Account_ID: req.Following_Account_ID,
	}

	err2 := db.InsertDocumentToDB(FollowersCol, followerData)
	if err2 != nil {
		c.SendString("Inserting tweet to the db failed.")
		return err2
	}

	resp.Success = true
	resp.ResponseType = "ACCOUNT_FOLLOWED"

	err3 := MarshalResponseAndSetBody(resp, c)
	if err3 != nil {
		return err3
	}

	return nil

}

func DoPasswordsMatch(password []byte, savedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), password)
	return err == nil
}

func IsTokenValid(tokenString string) bool {
	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token.Valid
}

func GetJWTclaims(tokenString string) *models.Claims {

	claims := &models.Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return claims
}

func MarshalResponseAndSetBody(resp interface{}, c *fiber.Ctx) error {
	result, err := json.Marshal(resp)
	if err != nil {
		c.SendString("Marshaling failed.")
		return err
	}

	c.Context().SetBody(result)
	return nil
}

func UnmarshalRequest(reqStruct interface{}, c *fiber.Ctx) error {

	err := json.Unmarshal(c.Body(), reqStruct)
	if err != nil {
		c.SendString("Unmarshaling failed.")
		return err
	}

	return nil
}

func SetMissingFieldResponse(fieldName string) models.BaseResponse {
	return models.BaseResponse{
		Success:      false,
		ResponseType: "FIELD_MISSING",
		Msg:          "Field " + "'" + fieldName + "'" + " is missing, or empty.",
	}
}

func SetCriteriaErrorResponse(fieldName string) models.BaseResponse {

	if fieldName == "age" {

		return models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_ERROR",
			Msg:          "Age should be 12 and above years old to create an account.",
		}

	} else if fieldName == "password" {

		return models.BaseResponse{
			Success:      false,
			ResponseType: "FIELD_ERROR",
			Msg:          "Password should atleast have 8 characters.",
		}

	}

	return models.BaseResponse{
		Success:      false,
		ResponseType: "UNKNOWN_ERROR",
		Msg:          "Can't find error type, SetCriteriaErrorResponse.",
	}
}
