package authentication

import (
	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	jwt.StandardClaims

	Name   string `json:"name"`
	UserId uint   `json:"user_id"`
}

type RefreshClaims struct {
	jwt.StandardClaims

	UserId uint `json:"user_id"`
}
