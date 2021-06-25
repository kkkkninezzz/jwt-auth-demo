package redis

import (
	"log"
	"time"
)

type RedisTemplate interface {
	Close()

	Set(key string, value interface{})

	SetEX(key string, value interface{}, expiration time.Duration)

	Get(key string) (val string)
}

// 对外部使用的模板类
var Template RedisTemplate

func newRedisTemplate(host string, port int32) RedisTemplate {
	var template RedisTemplate = newGoRedisTemplate(host, port)
	return template
}

func Connect(host string, port int32) {
	Template = newRedisTemplate(host, port)
	log.Println("Connection Opened to Redis")
}

func Shutdown() {
	if Template != nil {
		Template.Close()
	}
}
