package verify

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func VerifyToken(tokenStr string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Unauthorized")
		return false
	}
	return true

}
