package model

import "github.com/golang-jwt/jwt/v5"

type GoShopClaims struct {
	jwt.RegisteredClaims
	BaseClaims
}
type BaseClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
