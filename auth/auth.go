package auth

import (
	"log"

	"github.com/spf13/viper"

	"library_server/auth/jwt"
)

var (
	Jwt          *jwt.JWTAuth
	RmOpenApiJwt *jwt.JWTAuth
)

func Init() {
	_, err := jwt.Init(viper.GetString("jwt.prefix"), viper.GetString("jwt.sign_key"))
	if err != nil {
		log.Printf("jwt init error", err)
	}
}
