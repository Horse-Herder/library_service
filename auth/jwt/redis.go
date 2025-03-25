package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cast"

	"library_server/common"
	"library_server/utils"
)

const (
	USER_REDIS_KEY_MAX_NUM = 1
)

type redisStore struct {
	cli    common.RedisClient
	prefix string
}

// NewStore 创建基于redis存储实例
func NewRedisStore(cli common.RedisClient, keyPrefix string) *redisStore {
	return &redisStore{
		cli:    cli,
		prefix: keyPrefix,
	}
}

// Store redis存储

func (s *redisStore) wrapperKey(uid string, tokenType string) string {
	return fmt.Sprintf("%s:%s:%s", s.prefix, tokenType, uid)
}

func (s *redisStore) wrapperVal(orgId int64, token string, info map[string]interface{}) string {
	info["org_id"] = orgId
	info["token"] = utils.Md5(token)

	jsonByte, _ := json.Marshal(info)
	//return fmt.Sprintf("%v_%v", orgId, utils.Md5(token))
	return string(jsonByte)
}

// Set ...
func (s *redisStore) Set(token string, uid string, tokenType string, orgId int64, info map[string]interface{}, expiration time.Duration) error {
	cmd := s.cli.Set(context.Background(), s.wrapperKey(uid, tokenType), s.wrapperVal(orgId, token, info), expiration)
	return cmd.Err()
}

// Delete ...
func (s *redisStore) Delete(uid string, tokenType string) (bool, error) {
	key := s.wrapperKey(uid, tokenType)
	cmd := s.cli.Del(context.Background(), key)
	if err := cmd.Err(); err != nil {
		fmt.Printf("delete jwt token redis key error %v", cmd.Err())
		return false, err
	}
	return cmd.Val() > 0, nil
}

// Check ...
func (s *redisStore) Check(uid string, token string, tokenType string, orgId int64) (bool, map[string]interface{}, error) {
	info := make(map[string]interface{})
	cmd := s.cli.Get(context.Background(), s.wrapperKey(uid, tokenType))
	if err := cmd.Err(); err != nil {
		return false, info, err
	}

	_ = json.Unmarshal([]byte(cmd.Val()), &info)

	if cast.ToInt64(info["org_id"]) == orgId && cast.ToString(info["token"]) == utils.Md5(token) {
		return true, info, nil
	} else {
		return false, info, nil
	}
	//return cmd.Val() == s.wrapperVal(orgId, token), nil
}

func (s *redisStore) GetInfo(uid string, tokenType string) (map[string]interface{}, error) {
	info := make(map[string]interface{})
	cmd := s.cli.Get(context.Background(), s.wrapperKey(uid, tokenType))
	if err := cmd.Err(); err != nil {
		return info, err
	}
	_ = json.Unmarshal([]byte(cmd.Val()), &info)
	return info, nil
}

// Close ...
func (s *redisStore) Close() error {
	return s.cli.Close()
}
