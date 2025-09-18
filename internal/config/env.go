package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// Coba load dari ../../.env (DEV mode)
	if err := godotenv.Load("../../.env"); err != nil {
		// Jika gagal, coba dari root .env (PROD mode)
		if err2 := godotenv.Load(".env"); err2 != nil {
			log.Printf("Failed to load .env file from both paths: %v | fallback error: %v", err, err2)
		}
	}
}
