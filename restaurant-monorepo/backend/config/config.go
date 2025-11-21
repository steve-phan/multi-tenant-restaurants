package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DBUrl        string `mapstructure:"DATABASE_URL"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("CLIENT_ORIGIN", "http://localhost:3000")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Construct DBUrl if not present but individual vars are
	if AppConfig.DBUrl == "" {
		host := viper.GetString("POSTGRES_HOST")
		port := viper.GetString("POSTGRES_PORT")
		user := viper.GetString("POSTGRES_USER")
		password := viper.GetString("POSTGRES_PASSWORD")
		dbname := viper.GetString("POSTGRES_DATABASE")

		if host != "" && port != "" && user != "" && dbname != "" {
			AppConfig.DBUrl = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
		}
	}
}
