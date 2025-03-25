package jwt

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"library_server/common"
)

// 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
	ErrDeletedToken = errors.New("deleted token")
)

func Init(redisPrefix, signKey string) (*JWTAuth, error) {
	if viper.GetString("jwt.sign_key") == "" {
		log.Fatalf("jwt sign error")
	}

	if viper.GetInt("jwt.expired_time") <= 0 {
		log.Fatalf("jwt expired_time error")
	}

	if viper.GetString("jwt.sign_key") == "" {
		log.Fatalf("jwt prefix error")
	}

	var opts []Option

	opts = append(opts, SetExpired(viper.GetInt("jwt.expired_time")))
	opts = append(opts, SetSigningKey([]byte(signKey)))
	opts = append(opts, SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(signKey), nil
	}))

	method := jwt.SigningMethodHS256
	opts = append(opts, SetSigningMethod(method))

	auth := newJwt(NewRedisStore(common.GetRedis(), redisPrefix), opts...)

	return auth, nil
}
