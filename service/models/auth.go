package models

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}
