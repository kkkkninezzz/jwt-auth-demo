package redis

import "testing"

func TestTemplate(t *testing.T) {
	Connect("127.0.0.1", 6379)
	redisTemplate := Template
	defer redisTemplate.Close()

	redisTemplate.Set("testRRRR", "1231231")
	if redisTemplate.Get("testRRRR") != "1231231" {
		t.Error("输入不等")
	}

	redisTemplate.Del("testRRRR")
	redisTemplate.Del("testRRRR")
}
