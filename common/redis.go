package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"os"
	"time"
)

// ParametersRedis Parameters 读取json配置数据的结构体
type ParametersRedis struct {
	RedisIp       string `json:"redis_ip"`
	RedisPort     int    `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`
}

var REDISClient *redis.Client
var ctxRedis = context.Background()

// InitRedis REDIS InitRedis 初始化
func InitRedis() error {
	//读取config.json配置文件
	jsonFile, err := os.Open("./common/config.json")
	if err != nil {
		fmt.Println("open config.json failed!")
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error reading json file")
	}
	var parameters ParametersRedis
	json.Unmarshal(jsonData, &parameters)

	//建立连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", parameters.RedisIp, parameters.RedisPort),
		Password: parameters.RedisPassword,
		DB:       parameters.RedisDB,
	})
	ctx, cancel := context.WithTimeout(ctxRedis, 5*time.Second)
	defer cancel()
	_, err = redisClient.Ping(ctx).Result()
	REDISClient = redisClient
	return err

}

func GetRedis() *redis.Client {
	return REDISClient
}
