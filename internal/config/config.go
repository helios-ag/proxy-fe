package config

import "time"

const (
	CacheAuthorArticles   = 2 * time.Minute
	DetailedArticledCache = 5 * time.Minute
	CacheAuthorsList      = 5 * time.Minute
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
