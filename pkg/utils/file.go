package utils

import (
	"path/filepath"
	"strings"
)

func Join(elem ...string) string {
	return strings.ReplaceAll(filepath.Join(elem...), "\\", "/")
}
