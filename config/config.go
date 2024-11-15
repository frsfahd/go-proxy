package config

type Config struct {
	Redis struct {
		Host     string
		Port     string
		Username string
		Password string
		Index    int
	}
}
