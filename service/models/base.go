package models

import (
	"github.com/dgrijalva/jwt-go"
)

type BaseRequest struct {
	Token string `json:"token"`
}

type BaseResponse struct {
	Success      bool   `json:"success"`
	ResponseType string `json:"response_type,omitempty"`
	Msg          string `json:"Msg,omitempty"`
}

type Claims struct {
	User_UUID  string
	Account_ID string
	jwt.StandardClaims
}
