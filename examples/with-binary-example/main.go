package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/afero"
	"log"
	"os"
	"strings"
)

func main() {
	if errParse := ParseEnv(afero.NewOsFs(), ".env"); errParse != nil {
		log.Fatal(errParse)
	}
	DebugEnv("Database Host", "DATABASE_HOST")
	DebugEnv("Database Port", "DATABASE_PORT")
	DebugEnv("Database Name", "DATABASE_NAME")
	DebugEnv("Database Password", "DATABASE_PASSWORD")
}

func DebugEnv(description string, envName string) {
	fmt.Printf("%s: %s\n", description, os.Getenv(envName))
}

// ParseEnv is a custom env parser, not respecting any specs just for this example
func ParseEnv(fs afero.Fs, fileName string) error {
	file, err := fs.Open(fileName)
	if err != nil {
		return err
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
			return fmt.Errorf("could not set %s as %s: %s", envKey, envValue, errSetEnv.Error())

		}
	}

	if errScan := scanner.Err(); errScan != nil {
		return errScan
	}

	fmt.Println("Environment Used:", fileName)

	return nil

}
