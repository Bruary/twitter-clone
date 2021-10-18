package twitter

import (
	"os"
	"time"

	"github.com/Bruary/twitter-clone/service/models"
	"github.com/dgrijalva/jwt-go"
)

func CreateJWT(userUUID string, accountID string, validDurationMinutes time.Duration) (string, error) {

	// create the claims that will be used in the JWT token
	claims := &models.Claims{
		User_UUID:  userUUID,
		Account_ID: accountID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(validDurationMinutes * time.Minute).Unix(),
		},
	}

	// declaring the token with the method used for signing along with the claims§
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := os.Getenv("access_secret")

	// Create the JWT string
	tokenString, err4 := token.SignedString(jwtKey)
	if err4 != nil {
		return tokenString, err4
	}

	return tokenString, nil
}
