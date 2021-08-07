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
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

type DeleteUserRequest struct {
	Email string `json:"email"`
}

type BaseResponse struct {
	ResponseType string
	Success      bool
	Msg          string
}

var UsersCol *mongo.Collection

func main() {

	dbConn := db.GetDBConn()

	UsersCol = dbConn.Collection("Users")

	app := fiber.New()

	app.Post("/createUser", CreateUser)
	app.Delete("/deleteUser", DeleteUser)

	app.Listen(":3000")
}

// Creates a new user and adds it to the db
func CreateUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/json")

	var userInfo CreateUserRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &userInfo)
	if err != nil {
		c.SendString("Unmashaling failed in CreateUser endpoint.")
		return err
	}

	// check is user already exist in db
	doesUserExist, err10 := db.UserAlreadyExists(UsersCol, userInfo.Email)
	if err10 == nil && doesUserExist {
		resp.ResponseType = "USER_ALREADY_EXISTS"
		resp.Success = false
		resp.Msg = "User's email already exists."

		output01, err20 := json.Marshal(&resp)
		if err20 != nil {
			c.SendString("Mashaling failed in CreateUser endpoint.")
			return err20
		}

		c.Context().Response.SetBody(output01)

		return nil
	}

	// call the CreateUser func to add a new user to the db collection 'Users'
	err1 := db.CreateUser(UsersCol, userInfo)
	if err1 != nil {
		c.SendString("Inserting new user to the db failed.")
		return err1
	}

	resp.ResponseType = "NEW_USER_CREATED"
	resp.Success = true
	resp.Msg = "New user was added and saved to the db."

	output, err2 := json.Marshal(&resp)
	if err2 != nil {
		c.SendString("Mashaling failed in CreateUser endpoint.")
		return err
	}

	c.Context().Response.SetBody(output)
	return nil
}

// Delete user from the user using the email
func DeleteUser(c *fiber.Ctx) error {

	c.Context().SetContentType("application/jsons")

	var userEmail DeleteUserRequest
	var resp BaseResponse

	err := json.Unmarshal(c.Body(), &userEmail)
	if err != nil {
		c.SendString("Unmarshaling failed in DeleteUser endpoint.")
		return err
	}

	err1 := db.DeleteUser(UsersCol, userEmail.Email)
	if err1 != nil {

		resp.ResponseType = "UNKNOWN_ERROR"
		resp.Success = false
		resp.Msg = "Failed to delete user from DB."

	} else {

		resp.ResponseType = "USER_DELETED"
		resp.Success = true
		resp.Msg = "User has been successfuly deleted from the db."

	}

	result, err2 := json.Marshal(resp)
	if err2 != nil {
		c.SendString("Marshaling failed in DeleteUser endpoint.")
		return err2
	}

	c.Context().SetBody(result)
	return nil
}
