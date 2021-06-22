package redis

import "testing"

func TestTemplate(t *testing.T) {
	redisTemplate := NewRedisTemplate("127.0.0.1", 6379)

	defer redisTemplate.Close()

	redisTemplate.Set("testRRRR", "1231231")
	if redisTemplate.Get("testRRRR") != "1231231" {
		t.Error("输入不等")
	}
}
