package authenticate

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type Response struct {
	Token string `json:"token"`
}
