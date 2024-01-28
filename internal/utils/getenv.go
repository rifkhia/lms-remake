package utils

import (
	"os"
	"strconv"
)

func GetStrEnv(key string) string {
	return os.Getenv(key)
}

func GetIntEnv(key string) int {
	env, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return env
}
