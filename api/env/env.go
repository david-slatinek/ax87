package env

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

// Load environment variables from '.env' file.
func Load(fileName string) error {
	return godotenv.Load(fileName)
}

// Validate loaded env variables.
func Validate() error {
	if os.Getenv("INFLUXDB_TOKEN") == "" {
		return errors.New("empty INFLUXDB_TOKEN")
	}

	if os.Getenv("INFLUXDB_URL") == "" {
		return errors.New("empty INFLUXDB_URL")
	}

	if os.Getenv("ORGANIZATION") == "" {
		return errors.New("empty ORGANIZATION")
	}

	if os.Getenv("BUCKET") == "" {
		return errors.New("empty BUCKET")
	}

	return nil
}
