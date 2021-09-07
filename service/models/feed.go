package models

type FeedResponse struct {
	BaseResponse
	Tweets []Tweet `json:"tweets,omitempty"`
}
