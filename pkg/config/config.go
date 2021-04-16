package config

type ApplicationConfig struct {
	Test  string
	A     map[string]string
	Redis RedisConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
