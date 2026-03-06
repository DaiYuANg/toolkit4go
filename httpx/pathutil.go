package httpx

import (
	"strings"

	"github.com/samber/lo"
)

func joinRoutePath(basePath, path string) string {
	base := normalizeRoutePrefix(basePath)

	if path == "" {
		if base == "" {
			return "/"
		}
		return base
	}

	cleanPath := path
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	if base == "" {
		return cleanPath
	}

	if cleanPath == "/" {
		return base
	}

	return base + cleanPath
}

func normalizeRoutePrefix(prefix string) string {
	clean := strings.Trim(strings.TrimSpace(prefix), "/")
	return lo.Ternary(clean == "", "", "/"+clean)
}
