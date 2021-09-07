package twitter

import "github.com/Bruary/twitter-clone/service"

type twitter struct {
}

func NewTwitter() service.Service {
	return &twitter{}
}
