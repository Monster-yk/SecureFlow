package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct{
	DatabaseURL string
	JwtSecret string 
	Port string
	RedisURL string
	KafkaBrokers string 
	ClaudeAPIKey string
}

func getEnv(key string)(string,error){
	value:=os.Getenv(key)

	if value==""{
		return "",fmt.Errorf("%s is required",key)
	}

	return value,nil
}

func laod()(*Config,error){
	err:=godotenv.load()
	if err!=nil{
		return nil,fmt.Errorf("error loading .env file")
	}
	databaseURL, err := getEnv("DATABASE_URL")
    if err != nil {
        return nil, err
    }

    jwtSecret, err := getEnv("JWT_SECRET")
    if err != nil {
        return nil, err
    }

    port, err := getEnv("PORT")
    if err != nil {
        return nil, err
    }

    redisURL, err := getEnv("REDIS_URL")
    if err != nil {
        return nil, err
    }

    kafkaBrokers, err := getEnv("KAFKA_BROKERS")
    if err != nil {
        return nil, err
    }

    claudeAPIKey, err := getEnv("CLAUDE_API_KEY")
    if err != nil {
        return nil, err
    }

	cfg:=&Config{
		DatabaseURL: databaseURL,
		JwtSecret:    jwtSecret,
        Port:         port,
        RedisURL:     redisURL,
        KafkaBrokers: kafkaBrokers,
        ClaudeAPIKey: claudeAPIKey,
	}

	return cfg,nil
}