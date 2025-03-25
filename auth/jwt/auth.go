package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultKey = "gojwt"
const ACCESS_TOKEN = "access"
const OPENAPI_TOKEN = "openapi"
const REFRESH_TOKEN = "refresh"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

// Option 定义参数项
type Option func(*options)

// SetSigningMethod 设定签名方式
func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// SetSigningKey 设定签名key
func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// SetKeyfunc 设定验证key的回调函数
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// SetExpired 设定令牌过期时长(单位秒，默认7200)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// New 创建认证实例
func newJwt(store *redisStore, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

// JWTAuth jwt认证
type JWTAuth struct {
	opts  *options
	store *redisStore
}

// Generate 生成令牌
func (j *JWTAuth) Generate(claimsElement *ClaimsElement, info map[string]interface{}, expiredTime int64) (*TokenInfo, error) {
	now := time.Now()

	expiresAt := now.Add(time.Duration(86400) * time.Second)
	claims := &Claims{
		claimsElement.UserId,
		claimsElement.IsAdmin,
		claimsElement.UserName,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   claimsElement.UserId,
		},
	}

	token := jwt.NewWithClaims(defaultOptions.signingMethod, claims)
	tokenString, err := token.SignedString(defaultOptions.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &TokenInfo{
		ExpiresAt:   expiresAt.Unix(),
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

// Parse 解析令牌
func (j *JWTAuth) Parse(tokenString string) (*Claims, map[string]interface{}, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, defaultOptions.keyfunc)
	info := make(map[string]interface{})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Printf("jwt error 1:%v %v", err, tokenString)
			return nil, info, ErrExpiredToken
		} else {
			log.Printf("jwt error 2:%v %v", err, tokenString)
			return nil, info, ErrInvalidToken
		}
	}

	claims, _ := token.Claims.(*Claims)

	return claims, info, nil
}

func (j *JWTAuth) GetInfo(userId string, thirdId string) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	var err error

	err = j.callStore(func(store *redisStore) error {
		if info, err = store.GetInfo(userId+thirdId, ACCESS_TOKEN); err != nil {
			log.Printf("jwt get redis info error :%v", err.Error())
			return err
		}
		return nil
	})

	return info, nil
}

// Destroy 销毁令牌
func (j *JWTAuth) Destroy(uid string, tokenType string) error {
	//销毁redis中存的token
	return j.callStore(func(store *redisStore) error {
		_, err := store.Delete(uid, tokenType)
		return err
	})
}

// Release 释放资源
func (j *JWTAuth) Release() error {
	return j.callStore(func(store *redisStore) error {
		return store.Close()
	})
}

func (j *JWTAuth) callStore(fn func(*redisStore) error) error {
	if store := j.store; store != nil {
		return fn(store)
	}
	return nil
}
