package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
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
type UserInfo struct {
	UUID      string      `json:"uuid"`
	FirstName string      `json:"firstname"`
	LastName  string      `json:"lastname"`
	Age       int         `json:"age"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Metrics   UserMetrics `json:"-"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updates_at"`
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
type Tweet struct {
	UserUUID  string `json:"user_uuid"`
	TweetUUID string `json:"tweet_uuid"`
	Email     string `json:"email"`
	Tweet     string `json:"tweet"`
	Metrics   TweetMetrics
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TweetMetrics struct {
	Retweets_count   int
	Likes_count      int
	Comments_count   int
	Characters_count int
}

type Claims struct {
	UserUUID string
	jwt.StandardClaims
}

type GetTweetsRequest struct {
	UserUUID string `json:"user_uuid"`
	Token    string `json:"token"`
}

type GetTweetsResponse struct {
	Success bool     `json:"success"`
	Tweets  []bson.M `json:"tweets"`
}
