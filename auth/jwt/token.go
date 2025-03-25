package jwt

import (
	"encoding/json"
)

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	TokenType   string `json:"token_type"`   // 令牌类型
	ExpiresAt   int64  `json:"expires_at"`   // 令牌到期时间
}

func (t *TokenInfo) GetAccessToken() string {
	return t.AccessToken
}

func (t *TokenInfo) GetTokenType() string {
	return t.TokenType
}

func (t *TokenInfo) GetExpiresAt() int64 {
	return t.ExpiresAt
}

func (t *TokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}
