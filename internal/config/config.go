package config

import (
	"flag"
	"os"
)

type Config struct {
	DBDsn                string
	FlagRunAddr          string
	AccrualSystemAddress string
	JwtSecretKey         string
	RabbitURL            string
	RabbitQueue          string
}

func (c *Config) GetDBDsn() string {
	return c.DBDsn
}

func (c *Config) GetRabbitURL() string {
	return c.RabbitURL
}

func (c *Config) GetRabbitQueue() string {
	return c.RabbitQueue
}

func (c *Config) GetFlagRunAddr() string {
	return c.FlagRunAddr
}

func NewConfigCommand() (cf *Config) {
	config := new(Config)

	flag.StringVar(&config.DBDsn, "d", "", "db dsn")
	flag.StringVar(&config.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:3004", "AccrualSystemAddress")
	flag.StringVar(&config.JwtSecretKey, "k", "yoursecretkey", "jwtSecretKey")
	flag.StringVar(&config.RabbitURL, "u", "amqp://guest:guest@localhost:5672/", "RabbitURL")
	flag.StringVar(&config.RabbitQueue, "q", "processing", "RabbitQueue")

	flag.Parse()

	if Dsn := os.Getenv("DATABASE_URI"); Dsn != "" {
		config.DBDsn = Dsn
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		config.FlagRunAddr = envRunAddr
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.AccrualSystemAddress = envAccrualSystemAddress
	}

	if jwtSecretKey := os.Getenv("JWT_SECRET_KEY"); jwtSecretKey != "" {
		config.JwtSecretKey = jwtSecretKey
	}

	if rabbitURL := os.Getenv("RABBIT_URL"); rabbitURL != "" {
		config.RabbitURL = rabbitURL
	}

	if rabbitQueue := os.Getenv("RABBIT_QUEUE"); rabbitQueue != "" {
		config.RabbitQueue = rabbitQueue
	}

	return config
}
