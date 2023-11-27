package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func StringToInt(str string) (int, error) {

	if str == "" {
		return 0, errors.New("provided string is empty")
	}

	intValue, err := strconv.Atoi(str)

	if err != nil {
		return 0, err
	}

	return intValue, nil
}

func LoadEnvironmentVariables() error {

	requiredVariables := []string{"POSTGRES_DB", "POSTGRES_PASSWORD", "POSTGRES_USER", "JWT_SECRET", "REDIS_URL"}
	for _, variable := range requiredVariables {
		if os.Getenv(variable) == "" {
			return fmt.Errorf("missing required environment variable: %s", variable)
		}
		if os.Getenv("PORT") == "" {
			os.Setenv("PORT", "8080")
		}
	}

	return nil
}
