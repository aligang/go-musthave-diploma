package config

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"  envDefault:""`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:""`
	AccuralSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
}

func New() *Config {
	return &Config{}
}

func Init() *Config {
	config := New()
	ReadCLIParams(config)
	ReadENVParams(config)
	return config
}
