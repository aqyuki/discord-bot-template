package logging

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want slog.Level
	}{
		{
			name: "level debug",
			arg:  "debug",
			want: slog.LevelDebug,
		},
		{
			name: "level info",
			arg:  "info",
			want: slog.LevelInfo,
		},
		{
			name: "level warn",
			arg:  "warn",
			want: slog.LevelWarn,
		},
		{
			name: "level error",
			arg:  "error",
			want: slog.LevelError,
		},
		{
			name: "level unknown",
			arg:  "unknown",
			want: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := convertLevel(tt.arg)
			assert.Equal(t, tt.want, actual, "convertLevel(%s) = %v, want %v", tt.arg, actual, tt.want)
		})
	}
}

func TestFromEnv(t *testing.T) {
	t.Run("develop mode", func(t *testing.T) {
		t.Setenv(EnvLogMode, "develop")
		actual := FromEnv()
		assert.NotNil(t, actual)
	})

	t.Run("production mode", func(t *testing.T) {
		t.Setenv(EnvLogMode, "production")
		actual := FromEnv()
		assert.NotNil(t, actual)
	})

	t.Run("empty mode", func(t *testing.T) {
		t.Setenv(EnvLogMode, "")
		actual := FromEnv()
		assert.NotNil(t, actual)
	})
}

func TestNew(t *testing.T) {
	t.Run("develop mode", func(t *testing.T) {
		actual := New(true, "debug")
		assert.NotNil(t, actual)
	})

	t.Run("production mode", func(t *testing.T) {
		actual := New(false, "info")
		assert.NotNil(t, actual)
	})
}

func TestDefaultLogger(t *testing.T) {
	t.Run("default logger", func(t *testing.T) {
		actual := DefaultLogger()
		assert.NotNil(t, actual)
	})
}

func TestContext(t *testing.T) {
	t.Run("context with logger", func(t *testing.T) {
		ctx := context.Background()
		logger := DefaultLogger()

		actual := WithContext(ctx, logger)
		assert.NotNil(t, actual)

		actualLogger := FromContext(actual)
		assert.NotNil(t, actualLogger)
		assert.Equal(t, logger, actualLogger)
	})

	t.Run("context without logger", func(t *testing.T) {
		ctx := context.Background()
		actual := FromContext(ctx)
		assert.NotNil(t, actual)
	})
}
