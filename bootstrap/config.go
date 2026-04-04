package bootstrap

import (
	"github.com/joho/godotenv"
)

var Env map[string]string

func SetupConfig(path string) error {
	var err error
	Env, err = godotenv.Read(path)
	return err
}

func GetEnv(key string, defaultValue string) string {
	result := Env[key]
	if result == "" {
		return defaultValue
	}
	return result
}
