package c9bot

import (
	"log"
	"os"
)

// MustEnv gets an environment variable or panics
func MustEnv(env string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		log.Panicf("Environment variable $%v not set", env)
	}
	return v
}
