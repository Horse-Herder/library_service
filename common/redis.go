package common

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"library_server/utils"
)

type RedisClient interface {
	redis.Cmdable
	Close() error
	Ping(ctx context.Context) *redis.StatusCmd
	Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error
}

type RedisConfig struct {
	Switch         bool     `mapstructure:"switch"`
	ConnectionMode int      `mapstructure:"connection_mode"`
	Address        []string `mapstructure:"addr"`
	Password       string   `mapstructure:"password"`
	DB             int      `mapstructure:"db"`
}

var (
	redisClients = make(map[string]RedisClient)
)

func NewRedis(config *RedisConfig) (RedisClient, error) {

	//config := RedisConfig{}

	mode := ""

	if config.ConnectionMode == 0 {
		mode = "单机模式"
	} else if config.ConnectionMode == 1 {
		mode = "集群模式"
	}

	if client, ok := redisClients[getRedisClientKey(config.ConnectionMode, config.Address)]; ok {
		fmt.Printf("\nredis config:%+v 连接模式：%s\n", config, mode)

		return client, nil
	}
	var client RedisClient

	if config.ConnectionMode == 0 {
		client = redis.NewClient(&redis.Options{
			Addr:     config.Address[0],
			Password: config.Password,
			DB:       config.DB},
		)
	} else if config.ConnectionMode == 1 {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           config.Address,
			Password:        config.Password,
			PoolSize:        40,
			MinIdleConns:    10,
			DialTimeout:     5 * time.Second,
			ConnMaxIdleTime: 5 * time.Second,
		})
	}

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		return nil, fmt.Errorf("connect to redis failed, %v", err)
	}

	fmt.Printf("\nredis config:%+v 连接模式：%s\n", config, mode)

	redisClients[getRedisClientKey(config.ConnectionMode, config.Address)] = client

	return client, nil
}

func getRedisClientKey(connectionMode int, addr []string) string {
	return utils.Md5(fmt.Sprintf("%d:%s", connectionMode, strings.Join(addr, ",")))
}
