package usermodels

import "github.com/dgrijalva/jwt-go"

type Token struct {
	User_Id string `json:"user_id"`
	jwt.StandardClaims
}
