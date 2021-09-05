package models

import "time"

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

type FollowRequest struct {
	Following_Account_ID string `json:"following_account_id" bson:"following_account_id"`
	Token                string `json:"token"`
}

type Followers struct {
	ID                   string `json:"id"`                   // an ID unique to this document
	Follower_Account_ID  string `json:"follower_account_id"`  // the person who is following
	Following_Account_ID string `json:"following_account_id"` // the person being followed
}
