package log_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"arcadium.dev/dmx/assert"
	"arcadium.dev/dmx/log"
)

func TestContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(*slog.Logger) context.Context
		logger *slog.Logger
		verify func(*testing.T, *slog.Logger, *slog.Logger)
	}{
		{
			name: "test nil context",
			setup: func(*slog.Logger) context.Context {
				//lint:ignore SA1012 want to explicitly check nil context
				return log.IntoContext(nil, nil)
			},
			verify: func(t *testing.T, got, _ *slog.Logger) {
				assert.Compare(t, got.Handler(), slog.DiscardHandler)
			},
		},
		{
			name: "test empty context",
			setup: func(_ *slog.Logger) context.Context {
				return context.Background()
			},
			verify: func(t *testing.T, got, want *slog.Logger) {
				assert.Compare(t, got.Handler(), slog.DiscardHandler)
			},
		},
		{
			name: "test logger",
			setup: func(logger *slog.Logger) context.Context {
				return log.IntoContext(context.Background(), logger)
			},
			logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			verify: func(t *testing.T, got, want *slog.Logger) {
				assert.Equal(t, got, want)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := test.setup(test.logger)

			fmt.Printf("l: %v\n", ctx)

			logger := log.FromContext(ctx)
			test.verify(t, logger, test.logger)
		})
	}
}

func TestInfoContext(t *testing.T) {
}
