package models

import "time"

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

type CreateTweetRequest struct {
	Tweet string `json:"tweet"`
	Token string `json:"token"`
}

type GetTweetsResponse struct {
	Success bool    `json:"success"`
	Tweets  []Tweet `json:"tweets"`
}
