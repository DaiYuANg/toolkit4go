package authx

import (
	"log/slog"

	"github.com/DaiYuANg/arcgo/observabilityx"
)

func normalizeLogger(logger *slog.Logger) *slog.Logger {
	return observabilityx.NormalizeLogger(logger)
}
