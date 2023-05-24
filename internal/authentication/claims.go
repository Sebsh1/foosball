package authentication

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims

	Name           string `json:"name"`
	UserID         uint   `json:"user_id"`
	OrganizationID uint   `json:"organization_id"`
	Role           string `json:"role"`
}
