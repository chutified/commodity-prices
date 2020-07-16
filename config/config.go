package config

type Config struct {
	Host string
	Port int
}

func GetConfig() *Config {

	cfg := &Config{
		Host: "127.0.0.1",
		Port: 10501,
	}
	// TODO from file

	return cfg
}
