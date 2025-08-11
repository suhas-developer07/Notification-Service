package config

import (
	"fmt"
	"os"
)

type Config struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
}

func LoadConfig() Config {
	fmt.Print(os.Getenv("SMTP_HOST"))
	return Config{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
	}
}
