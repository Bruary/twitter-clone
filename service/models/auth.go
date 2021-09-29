package models

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type SetNewPasswordRequest struct {
	BaseRequest
	Password string `json:"password"`
}
