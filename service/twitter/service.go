package twitter

import "github.com/Bruary/twitter-clone/service"

type twitter struct {
}

// NewTwitter: fill the interface with the following struct
func NewTwitter() service.Service {
	return &twitter{}
}
