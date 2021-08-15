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
	"time"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/validate"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type BaseResponse struct {
	ResponseType string
	Success      bool
	Msg          string
}

type CreateUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// User info to be saved in the db
type UserRequiredFields struct {
	UUID          string
	FirstName     string
	LastName      string
	Age           int
	Email         string
	Password      string
	Metrics       UserMetrics
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

type UserMetrics struct {
	Followers_count      int
	Total_tweets_count   int
	Total_retweets_count int
	Total_likes_count    int
}

type DeleteUserRequest struct {
	Email string `json:"email"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Success bool
	Token   string `json:"token"`
}

type MakeATweetRequest struct {
	Email string `json:"email"`
	Tweet string `json:"tweet"`
	Token string `json:"token"`
}

// Tweets info to be saved in the db
type TweetRequiredFields struct {
	UserUUID      string
	TweetUUID     string
	Email         string
	Tweet         string
	Metrics       TweetMetrics
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

type TweetMetrics struct {
	Retweets_count   int
	Likes_count      int
	Comments_count   int
	Characters_count int
}

type Claims struct {
	Email string
	jwt.StandardClaims
}

var UsersCol *mongo.Collection
var TweetsCol *mongo.Collection

var jwtKey = []byte("White_Yasmin")

func main() {

	dbConn := db.GetDBConn()

	UsersCol = dbConn.Collection("Users")
	TweetsCol = dbConn.Collection("Tweets")

	app := fiber.New()

	app.Post("/createUser", CreateUser)
	app.Delete("/deleteUser", DeleteUser)
	app.Post("/makeATweet", MakeATweet)
	app.Post("/signin", SignIn)

	app.Listen(":3000")
}

// Creates a new user and adds it to the db
func CreateUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/json")

	var userReceivedInfo CreateUserRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &userReceivedInfo)
	if err != nil {
		c.SendString("Unmashaling failed in CreateUser endpoint.")
		return err
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

	// Fill in the user details
	userInfo := &UserRequiredFields{
		UUID:      uuid.NewV4().String(),
		FirstName: userReceivedInfo.FirstName,
		LastName:  userReceivedInfo.LastName,
		Age:       userReceivedInfo.Age,
		Email:     userReceivedInfo.Email,
		Password:  string(passwordHashedAndSalted),
		Metrics: UserMetrics{
			Followers_count:      0,
			Total_tweets_count:   0,
			Total_retweets_count: 0,
			Total_likes_count:    0,
		},
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
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

	var userEmail DeleteUserRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &userEmail)
	if err != nil {
		c.SendString("Unmarshaling failed in DeleteUser endpoint.")
		return err
	}

	err1 := db.DeleteUser(UsersCol, userEmail.Email)
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

	var req SignInRequest
	var baseResp BaseResponse

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		c.SendString("Unmarshaling failed in SignIn endpoint.")
		return err
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
	userDocument, err2_5 := db.GetDocumentFromDB(UsersCol, req.Email)
	if err2_5 != nil {
		c.SendString("Finding the user document in the DB failed in SignIn endpoint.")
		return err2_5
	}

	// Check if the request password matches the one stored in the DB
	var userDocumentDecoded UserRequiredFields

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
	claims := &Claims{
		Email: req.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
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
	resp := &SignInResponse{
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

	var tweetRequest MakeATweetRequest
	var resp BaseResponse

	err := UnmarshalRequest(&tweetRequest, c)
	if err != nil {
		return err
	}

	// Request validation
	tokenValueEmpty := validate.IsStringEmpty(tweetRequest.Token)
	if tokenValueEmpty {
		resp = SetRequestErrorResponse("token")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	emailValueEmpty := validate.IsStringEmpty(tweetRequest.Email)
	if emailValueEmpty {
		resp = SetRequestErrorResponse("email")
		c.Status(fiber.ErrBadRequest.Code)

		err := MarshalResponseAndSetBody(resp, c)
		if err != nil {
			return err
		}
		return nil
	}

	tweetValueEmpty := validate.IsStringEmpty(tweetRequest.Tweet)
	if tweetValueEmpty {
		resp = SetRequestErrorResponse("tweet")
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

		err1_4 := MarshalResponseAndSetBody(resp, c)
		if err1_4 != nil {
			return err1_4
		}

		return nil
	}

	// Check if user exists
	userDoc, err1_5 := db.GetDocumentFromDB(UsersCol, tweetRequest.Email)
	if err1_5 != nil {
		c.SendString("Failed while finding user in db.")
		return err1_5
	}

	var user UserRequiredFields

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
	var tweetInfo = TweetRequiredFields{
		UserUUID:  user.UUID,
		TweetUUID: uuid.NewV4().String(),
		Email:     tweetRequest.Email,
		Tweet:     tweetRequest.Tweet,
		Metrics: TweetMetrics{
			Retweets_count:   0,
			Likes_count:      0,
			Comments_count:   0,
			Characters_count: len(tweetRequest.Tweet),
		},
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
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

func DoPasswordsMatch(password []byte, savedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), password)
	return err == nil
}

func IsTokenValid(tokenString string) bool {
	claims := &Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token.Valid
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

func SetRequestErrorResponse(fieldName string) BaseResponse {
	return BaseResponse{
		Success:      false,
		ResponseType: "FIELD_MISSING",
		Msg:          "Field " + "'" + fieldName + "'" + " is missing, or empty",
	}
}
