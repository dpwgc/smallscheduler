package common

import (
	"os"
)

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0766)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
