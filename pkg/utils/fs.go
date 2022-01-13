package utils

import (
	"os"
)

func CopyFile(src string, dst string) error {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, bytesRead, 0755)
	if err != nil {
		return err
	}
	return nil
}
