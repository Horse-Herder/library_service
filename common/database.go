package common

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"library_server/model"
)

func InitDB() {
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	charset := viper.GetString("datasource.charset")
	parseTime := viper.GetString("datasource.parseTime")
	loc := viper.GetString("datasource.loc")

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		username, password, host, port, database, charset, parseTime, url.QueryEscape(loc))

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("fail to init database: %s\n", err))
	}
	// 数据库自动迁移
	DB = db
	DB.AutoMigrate(&model.Admin{})
	DB.AutoMigrate(&model.Book{})
	DB.AutoMigrate(&model.Reader{})
	DB.AutoMigrate(&model.Reserve{})
	DB.AutoMigrate(&model.Comment{})
	DB.AutoMigrate(&model.Borrow{})
	DB.AutoMigrate(&model.Report{})

}

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func GetRedis() RedisClient {

	redis, err := NewRedis(&RedisConfig{
		Switch:         true,
		ConnectionMode: viper.GetInt("redis.connection_mode"),
		Address:        viper.GetStringSlice("redis.addr"),
		Password:       viper.GetString("redis.password"),
		DB:             viper.GetInt("redis.db"),
	})
	if err != nil {
		panic(fmt.Sprintf("fail to init redis: %s\n", err))

	}
	return redis
}
