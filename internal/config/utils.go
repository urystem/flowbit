package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func mustGetEnvString(key string) string {
	str := os.Getenv(key)
	if str == "" {
		log.Fatalln(key, "not found in env")
	}
	return str
}

func mustGetEnvInt(key string) int {
	str := mustGetEnvString(key)
	n, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalln(key, "invalid", str)
	}
	return n
}

func mustGetBoolean(key string) bool {
	switch mustGetEnvString(key) {
	case "true":
		return true
	case "false":
		return false
	default:
		log.Fatalln(key, "invalid")
		return false
	}
}

func mustGetArrayStr(key string) []string {
	str := mustGetEnvString(key)
	return strings.Split(str, ",") // []string{"40101", "40102", "40103"}
}
