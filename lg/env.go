package lg

import (
	"os"
	"strconv"
)

func GetenvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	return IfeI(err == nil, value, defaultValue)
}

func GetenvStr(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	return IfeS(exists, value, defaultValue)
}
