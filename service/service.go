package service

import "github.com/gofiber/fiber/v2"

type Service interface {
	CreateUser(*fiber.Ctx) error
	DeleteUser(*fiber.Ctx) error
	SignIn(*fiber.Ctx) error
	CreateTweet(*fiber.Ctx) error
	GetTweets(*fiber.Ctx) error
	Follow(*fiber.Ctx) error
	Feed(*fiber.Ctx) error
}
