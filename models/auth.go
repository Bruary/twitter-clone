package models

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Success bool
	Token   string `json:"token"`
}
