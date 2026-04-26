package observability

import (
	"log/slog"
	"os"
)

// NewLogger creates a JSON structured logger.
// Attaches service and project to every log line.
// Call once in main() and set as default:
//
//	slog.SetDefault(observability.NewLogger("api"))
func NewLogger(serviceName string) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel(),
		}),
	).With(
		"service", serviceName,
		"project", "city-explorer",
	)
}

// logLevel reads LOG_LEVEL from environment.
// Defaults to Info.
// Set LOG_LEVEL=debug in .env for verbose output.
func logLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
