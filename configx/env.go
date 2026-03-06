package configx

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	envProvider "github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

// loadDotenv loads related configuration.
// ignoreErr documents related behavior.
func loadDotenv(files []string, ignoreErr bool) error {
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			if ignoreErr {
				continue
			}
			if os.IsNotExist(err) {
				return fmt.Errorf("configx: dotenv file %q not found: %w", f, err)
			}
			return fmt.Errorf("configx: stat dotenv file %q: %w", f, err)
		}

		if err := godotenv.Load(f); err != nil {
			if ignoreErr {
				continue
			}
			return fmt.Errorf("configx: load dotenv file %q: %w", f, err)
		}

		// Note.
	}
	return nil
}

// loadEnv loads related configuration.
// prefix documents related behavior.
// Note.
// Note.
func loadEnv(k *koanf.Koanf, prefix string) error {
	normalizedPrefix := normalizeEnvPrefix(prefix)

	p := envProvider.Provider(".", envProvider.Opt{
		Prefix: normalizedPrefix,
		TransformFunc: func(k, v string) (string, any) {
			keyWithoutPrefix := strings.TrimPrefix(k, normalizedPrefix)
			keyWithoutPrefix = strings.TrimPrefix(keyWithoutPrefix, "_")

			// Note.
			key := strings.ReplaceAll(strings.ToLower(keyWithoutPrefix), "_", ".")
			return key, v
		},
		EnvironFunc: os.Environ,
	})

	if err := k.Load(p, nil); err != nil {
		return fmt.Errorf("configx: load env prefix %q: %w", normalizedPrefix, err)
	}
	return nil
}

func normalizeEnvPrefix(prefix string) string {
	clean := strings.TrimSpace(prefix)
	if clean == "" {
		return ""
	}
	return strings.TrimSuffix(clean, "_") + "_"
}
