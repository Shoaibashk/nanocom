/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package transfer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands and validates a file path
func ExpandPath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("path is required")
	}
	if strings.HasPrefix(trimmed, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			switch {
			case trimmed == "~":
				trimmed = home
			case len(trimmed) > 1 && (trimmed[1] == '/' || trimmed[1] == '\\'):
				trimmed = filepath.Join(home, trimmed[2:])
			}
		}
	}
	absPath, err := filepath.Abs(trimmed)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
