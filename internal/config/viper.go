package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config.config") //TODO: <-- PENTING! Ganti path di sini (DEVLOPMENT MODE)
	// config.SetConfigName("vps.config") //TODO: <-- PENTING! Ganti path di sini (PRODUCTION MODE)
	config.SetConfigType("json")
	config.AddConfigPath("../../.") //TODO: <-- PENTING! Ganti path di sini (DEVLOPMENT MODE)
	config.AddConfigPath(".")       //TODO: <-- PENTING! Ganti path di sini (PRODUCTION MODE)
	err := config.ReadInConfig()
	// Jika gagal membaca file, coba baca dari environment variable
	if err != nil {

		// fallback: coba baca dari ENV
		envConfig := os.Getenv("ATTENDANCE_APP")
		if envConfig == "" {
			panic(fmt.Errorf("failed to read config file and APPLICATION_CONFIG not set: %w", err))
		}

		// parse JSON string di env
		var jsonMap map[string]any
		if err := json.Unmarshal([]byte(envConfig), &jsonMap); err != nil {
			panic(fmt.Errorf("failed to unmarshal APPLICATION_CONFIG: %w", err))
		}

		// load JSON map ke Viper
		if err := config.MergeConfigMap(jsonMap); err != nil {
			panic(fmt.Errorf("failed to merge APPLICATION_CONFIG into viper: %w", err))
		}
		fmt.Print("Using config from environment variable APPLICATION_CONFIG\n")
	}
	return config
}
func GetJWTSecret(v *viper.Viper) string {
	return v.GetString("jwt.secret")
}

func MustGetString(v *viper.Viper, key string) string {
	value := v.GetString(key)
	if value == "" {
		panic(fmt.Errorf("missing required config key: %s", key))
	}
	return value
}
