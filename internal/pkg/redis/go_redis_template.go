package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type goRedisTemplate struct {
	ctx       context.Context
	rdbClient *redis.Client
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

func (template *goRedisTemplate) SetEX(key string, value interface{}, expiration time.Duration) {
	err := template.rdbClient.SetEX(template.ctx, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
}

func (template *goRedisTemplate) Del(key string) {
	err := template.rdbClient.Del(template.ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
