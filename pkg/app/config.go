package app

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	Port     string
	Host     string
	HostName string
	Dns      string
	SmsToken string
	Email    string
}

func CreateConfig() Config {
	var c Config
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Error loading .env file. Loading default...")
		return Config{
			Port:     "8080",
			Host:     "0.0.0.0",
			HostName: "localhost",
			Dns:      "host=localhost user=nikolai password=nikolai dbname=persons",
			SmsToken: "bsineEFPwM-HN6fn4xSdj3EIiPbkhvGG",
			Email:    "nikolay.sychev.1999@gmail.com",
		}
	}
	c.Port = os.Getenv("PORT")
	if c.Port == "" {
		c.Port = "8080"
	}

	c.Host = os.Getenv("HOST")
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}

	c.HostName = os.Getenv("HOSTNAME")
	if c.HostName == "" {
		c.HostName = "localhost"
	}

	c.Dns = os.Getenv("DNS")
	if c.Dns == "" {
		c.Dns = "host=localhost user=nikolai password=nikolai dbname=dating"
	}

	c.SmsToken = os.Getenv("SMS_TOKEN")

	c.Email = os.Getenv("EMAIL")
	return c
}
