package service

import (
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	CreateUser(*fiber.Ctx, models.CreateUserRequest) models.BaseResponse
	SignIn(*fiber.Ctx, models.SignInRequest) *models.SignInResponse
	DeleteUser(*fiber.Ctx) error
	CreateTweet(*fiber.Ctx) error
	GetTweets(*fiber.Ctx) error
	Follow(*fiber.Ctx) error
	Feed(*fiber.Ctx) error
}
