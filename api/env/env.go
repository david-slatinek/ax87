package env

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

// Load environment variables from the '.env' file.
func Load(fileName string) error {
	return godotenv.Load(fileName)
}

// Validate loaded env variables.
func Validate() error {
	if os.Getenv("INFLUXDB_URL") == "" {
		return errors.New("empty INFLUXDB_URL")
	}

	if os.Getenv("INFLUXDB_TOKEN") == "" {
		return errors.New("empty INFLUXDB_TOKEN")
	}

	if os.Getenv("ORGANIZATION") == "" {
		return errors.New("empty ORGANIZATION")
	}

	if os.Getenv("BUCKET") == "" {
		return errors.New("empty BUCKET")
	}

	if os.Getenv("REDIS_URL") == "" {
		return errors.New("empty REDIS_URL")
	}

	if os.Getenv("GO_ENV") == "" {
		return errors.New("empty GO_ENV")
	}

	return nil
}
