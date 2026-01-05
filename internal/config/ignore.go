package config

import (
	"os"
	"strings"
)

func LoadIgnoreFile() []string {
	data, err := os.ReadFile(".driftignore")
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var ignores []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		ignores = append(ignores, l)
	}
	return ignores
}

func IsIgnored(path string, ignores []string) bool {
	for _, ig := range ignores {
		if ig == path {
			return true
		}
	}
	return false
}
