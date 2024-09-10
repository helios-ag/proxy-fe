package config

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
			Host: "127.0.0.1",
			Port: 6379,
		},
	}
}
