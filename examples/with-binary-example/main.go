package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	ParseEnv(".env")

	DebugEnv("Database Host", "DATABASE_HOST")
	DebugEnv("Database Port", "DATABASE_PORT")
	DebugEnv("Database Name", "DATABASE_NAME")
	DebugEnv("Database Password", "DATABASE_PASSWORD")
}

func DebugEnv(description string, envName string) {
	fmt.Printf("%s: %s\n", description, os.Getenv(envName))
}

// ParseEnv is a custom env parser, not respecting any specs just for this example
func ParseEnv(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") {
			continue
		}
		keyAndValue := strings.Split(line, "=")
		envKey, envValue := keyAndValue[0], keyAndValue[1]
		if errSetEnv := os.Setenv(envKey, envValue); errSetEnv != nil {
			fmt.Printf("could not set %s as %s: %s", envKey, envValue, errSetEnv.Error())
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Environment Used:", fileName)

}
