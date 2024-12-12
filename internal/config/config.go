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
}

func (c *Config) GetDBDsn() string {
	return c.DBDsn
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

	return config
}
