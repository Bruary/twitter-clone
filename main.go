// main functionality:
// 1) Create new user
// 2) Delete account
// 3) Post a tweet
// 4) Make account/tweets private or public
// 5) News feed

package main

import (
	"encoding/json"

	"github.com/Bruary/twitter-clone/db"
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/Bruary/twitter-clone/service/twitter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// Connect to the db
	db.SetUpDBConnection()

	svc := twitter.NewTwitter()

	app := fiber.New()

	app.Use(cors.New())

	api := app.Group("/api") // api/

	v1 := api.Group("/v1") // api/v1/

	auth := v1.Group("/auth") // api/v1/auth/

	auth.Post("/signin", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.SignInRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the sign in logic
		resp := svc.SignIn(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	auth.Post("/signup", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.CreateUserRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the create user logic
		resp := svc.CreateUser(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	auth.Post("/resetPassword", func(c *fiber.Ctx) error {
		c.Context().SetContentType("application/jsons")

		req := models.ResetPasswordRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the reset password logic
		resp := svc.ResetPassword(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	auth.Get("/resetPassword/newPassword", func(c *fiber.Ctx) error {
		c.Context().SetContentType("application/jsons")

		// run the reset password logic
		resp := svc.NewPassword(c)

		if err := MarshalResponseAndSetBody(resp, c); err != nil {
			return err
		}

		return nil
	})

	user := v1.Group("/user") // api/v1/user/
	user.Delete("/delete", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.DeleteUserRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the delete user logic
		resp := svc.DeleteUser(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	tweet := v1.Group("/tweets") // api/v1/tweet/

	tweet.Post("/create", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.CreateTweetRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the create tweet logic
		resp := svc.CreateTweet(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	tweet.Post("/get", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.BaseRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the get tweets logic
		resp := svc.GetTweets(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	v1.Post("/follow", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.FollowRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the follow logic
		resp := svc.Follow(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	v1.Post("/feed", func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/jsons")

		req := models.BaseRequest{}
		if err := UnmarshalRequest(&req, c); err != nil {
			return err
		}

		// run the follow logic
		resp := svc.Feed(c, req)

		if err2 := MarshalResponseAndSetBody(resp, c); err2 != nil {
			return err2
		}

		return nil
	})

	app.Listen(":4000")
}

func MarshalResponseAndSetBody(resp interface{}, c *fiber.Ctx) error {
	result, err := json.Marshal(resp)
	if err != nil {
		c.SendString("Marshaling failed.")
		return err
	}

	c.Context().SetBody(result)
	return nil
}

func UnmarshalRequest(reqStruct interface{}, c *fiber.Ctx) error {

	err := json.Unmarshal(c.Body(), reqStruct)
	if err != nil {
		c.SendString("Unmarshaling failed.")
		return err
	}

	return nil
}
