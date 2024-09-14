package config

import "time"

const (
	CacheAuthorArticles   = 2 * time.Minute
	DetailedArticledCache = 5 * time.Minute
	CacheAuthorsList      = 5 * time.Minute
	UsersURL              = "https://jsonplaceholder.typicode.com/users"
	PostsURL              = "https://jsonplaceholder.typicode.com/posts"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host string
	Port int
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host: "proxy-service-db.local",
			Port: 6379,
		},
	}
}
