package jwt

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type ClaimsElement struct {
	UserId   string
	IsAdmin  bool
	UserName string
}

type Claims struct {
	UserId   string
	IsAdmin  bool
	UserName string
	jwt.RegisteredClaims
}

func (c *Claims) GetExpiresAt() int64 {
	return c.ExpiresAt.Unix()
}

func (c *Claims) GetUserId() string {
	return c.UserId
}

func (c *Claims) GetUserName() string {
	return c.UserName
}

func (c *Claims) GetIsAdmin() bool {
	return c.IsAdmin
}
