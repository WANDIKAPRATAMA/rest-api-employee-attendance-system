package main

import "employee-attendance-system/internal/config"

func main() {
	viper := config.NewViper()
	log := config.NewLogger(viper)
	database := config.NewDatabase(viper, log)
	validator := config.NewValidator(viper)
	fiber := config.NewFiber(viper)

	cfg := &config.AppConfig{
		DB:       database,
		App:      fiber,
		Log:      log,
		Validate: validator,
		Viper:    viper,
	}

	config.NewAppConfig(cfg)

}
