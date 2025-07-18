package utils

import (
	"fmt"
	"os"
)

func CheckEnvKey(keys []string) error {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			return fmt.Errorf("env with key '%s' is not fill", key)
		}
	}
	return nil
}

func MessageError(header string, err error) error {
	return fmt.Errorf("[%s] %v", header, err)
}
