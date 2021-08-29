package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type BaseRequest struct {
	Token string `json:"token"`
}

type BaseResponse struct {
	Success      bool   `json:"success"`
	ResponseType string `json:"response_type"`
	Msg          string `json:"Msg,omitempty"`
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
	UUID       string      `json:"uuid"`
	Account_ID string      `json:"account_id" bson:"account_id"`
	FirstName  string      `json:"firstname"`
	LastName   string      `json:"lastname"`
	Age        int         `json:"age"`
	Email      string      `json:"email"`
	Password   string      `json:"password"`
	Metrics    UserMetrics `json:"-"`
	Created_At time.Time   `json:"created_at" bson:"created_at"`
	Updated_At time.Time   `json:"updates_at" bson:"updated_at"`
}

type FollowRequest struct {
	Following_Account_ID string `json:"following_account_id" bson:"following_account_id"`
	Token                string `json:"token"`
}

type Followers struct {
	ID                   string `json:"id"`                   // an ID unique to this document
	Follower_Account_ID  string `json:"follower_account_id"`  // the person who is following
	Following_Account_ID string `json:"following_account_id"` // the person being followed
}

type UserMetrics struct {
	Followers_count      int
	Following_count      int
	Total_tweets_count   int
	Total_retweets_count int
	Total_likes_count    int
}

type DeleteUserRequest struct {
	Token string `json:"token"`
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
	Tweet string `json:"tweet"`
	Token string `json:"token"`
}

// Tweets info to be saved in the db
type TweetDB struct {
	User_UUID       string `json:"user_uuid"`
	User_Account_ID string `json:"user_account_id"`
	Tweet_UUID      string `json:"tweet_uuid"`
	Email           string `json:"email"`
	Tweet           string `json:"tweet"`
	Metrics         TweetMetrics
	Created_At      time.Time `json:"created_at"`
	Updated_At      time.Time `json:"updated_at"`
}

type Tweet struct {
	Tweet_UUID string `json:"tweet_uuid"`
	Tweet      string `json:"tweet"`
	Metrics    TweetMetrics
}

type TweetMetrics struct {
	Retweets_count   int
	Likes_count      int
	Comments_count   int
	Characters_count int
}

type Claims struct {
	User_UUID  string
	Account_ID string
	jwt.StandardClaims
}

type GetTweetsRequest struct {
	Token string `json:"token"`
}

type GetTweetsResponse struct {
	Success bool    `json:"success"`
	Tweets  []Tweet `json:"tweets"`
}
