package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	godotenv.Load()
}

// Helper function to lookup the environment variable, or use a fallback
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		println("env.GetString err")
		return fallback
	}
	return val
}

func GetCallbackUrlHost() string {
	val, _ := os.LookupEnv("ENV")
	var res string

	if val == "PRD" {
		res, _ = os.LookupEnv("PRD_CALLBACK_URL_HOST")
	} else {
		res, _ = os.LookupEnv("DEV_CALLBACK_URL_HOST")
	}
	log.Printf("GetCallbackUrlHost() ENV: %v, using %v", val, res)
	return res
}

func GetDBAddr() string {
	val, _ := os.LookupEnv("ENV")
	var res string

	if val == "PRD" {
		res, _ = os.LookupEnv("PRD_DB_ADDR")
	} else {
		res, _ = os.LookupEnv("DEV_DB_ADDR")
	}
	return res
}
