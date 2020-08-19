package global

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var Rds *redis.Client

func SetRedis() {
	sessConf := Config.Section("redis")
	addr := sessConf.Key("ip").String()
	port := sessConf.Key("port").String()
	pass := sessConf.Key("password").String()
	db := sessConf.Key("db").MustInt()

	Rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", addr, port),
		Password: pass, // no password set
		DB:       db,   // use default DB
	})

	_, err := Rds.Ping().Result()
	if err != nil {
		log.Println(err)
		return
	}
}
