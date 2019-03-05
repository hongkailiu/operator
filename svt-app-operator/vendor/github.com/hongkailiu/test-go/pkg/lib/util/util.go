package util

import (
	"math/rand"
	"os"
	"time"
)

// GetRandomInt return a random integer
func GetRandomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

// PanicIfError panics if error is not nil
func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Getenv returns value of env. var if it is defined; defaultValue otherwise.
func Getenv(key, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}
	return result
}

// HomeDir returns home dir
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
