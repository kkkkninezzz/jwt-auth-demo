package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisTemplate interface {
	Close()

	Set(key string, value interface{})

	Get(key string) (val string)
}

type goRedisTemplate struct {
	ctx       context.Context
	rdbClient *redis.Client
}

func NewRedisTemplate(host string, port int32) RedisTemplate {
	var template RedisTemplate = newGoRedisTemplate(host, port)
	return template
}

func newGoRedisTemplate(host string, port int32) *goRedisTemplate {
	template := &goRedisTemplate{}
	template.ctx = context.Background()

	addr := host + ":" + strconv.Itoa(int(port))
	template.rdbClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return template
}

func (template *goRedisTemplate) Close() {
	template.rdbClient.Close()
}

func (template *goRedisTemplate) Set(key string, value interface{}) {
	err := template.rdbClient.Set(template.ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (template *goRedisTemplate) Get(key string) (val string) {
	val, err := template.rdbClient.Get(template.ctx, key).Result()
	if err == redis.Nil {
		val = ""
	} else if err != nil {
		panic(err)
	}
	return
}
