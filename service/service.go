package service

import (
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	CreateUser(*fiber.Ctx, models.CreateUserRequest) *models.BaseResponse
	SignIn(*fiber.Ctx, models.SignInRequest) *models.SignInResponse
	DeleteUser(*fiber.Ctx, models.DeleteUserRequest) *models.BaseResponse
	CreateTweet(*fiber.Ctx, models.CreateTweetRequest) *models.BaseResponse
	GetTweets(*fiber.Ctx, models.BaseRequest) *models.GetTweetsResponse
	Follow(*fiber.Ctx, models.FollowRequest) *models.BaseResponse
	Feed(*fiber.Ctx, models.BaseRequest) *models.FeedResponse
}
