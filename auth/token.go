package auth

import (
	"fmt"
	"time"

	"my_app/helper"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mihit")

func CreateToken(profile helper.Profile) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   profile.Id,
			"name": profile.Name,
			"role": profile.Role,
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}
