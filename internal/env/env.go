package env

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	godotenv.Load()
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		println("env.GetString err")
		return fallback
	}
	return val
}
