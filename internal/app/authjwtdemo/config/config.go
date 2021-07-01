package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	FiberAddr   string      `yaml:"fiber_addr"`
	RedisConfig RedisConfig `yaml:"redis_config"`
	MysqlConfig MysqlConfig `yaml:"mysql_config"`
	JwtConfig   JwtConfig   `yaml:"jwt_config"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Password string `yaml:"pasword"`
}

type MysqlConfig struct {
	Dsn           string `yaml:"dsn"`
	TablePrefix   string `yaml:"table_prefix"`
	SingularTable bool   `yaml:"singular_table"`
}

type JwtConfig struct {
	// 计算token的salt的过期时间 单位为小时
	TokenSaltExpiration time.Duration `yaml:"token_salt_expiration"`
	// token的过期时间 单位为小时
	TokenExpiration time.Duration `yaml:"token_expiration"`
	PrivateSecret   string        `yaml:"private_secret"`
	Issuer          string        `yaml:"issuer"`
}

var Config *AppConfig

func Init(configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := &AppConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}

	Config = config
	Config.JwtConfig.TokenSaltExpiration = Config.JwtConfig.TokenSaltExpiration * time.Hour
	Config.JwtConfig.TokenExpiration = Config.JwtConfig.TokenExpiration * time.Hour
}
