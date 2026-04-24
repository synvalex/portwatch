package config

import (
	"fmt"
	"os"
)

// FindConfigFile searches common locations for a portwatch config file
// and returns the first path that exists. An error is returned if none
// are found and no explicit path was provided.
func FindConfigFile(explicit string) (string, error) {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			return "", fmt.Errorf("config file %q not found", explicit)
		}
		return explicit, nil
	}

	candidates := []string{
		"portwatch.yaml",
		"portwatch.yml",
		"/etc/portwatch/portwatch.yaml",
		"/etc/portwatch/portwatch.yml",
	}

	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			home+"/.config/portwatch/portwatch.yaml",
			home+"/.portwatch.yaml",
		)
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("no config file found; searched: %v", candidates)
}

// MustLoad loads the config from the given path or panics on error.
// Intended for use in tests or simple CLI entry points.
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(fmt.Sprintf("portwatch: config load failed: %v", err))
	}
	return cfg
}
