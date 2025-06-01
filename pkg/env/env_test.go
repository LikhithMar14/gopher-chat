package env

import (
	"os"
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		want     string
	}{
		{
			name:     "existing env var",
			key:      "TEST_STRING",
			value:    "test-value",
			fallback: "fallback",
			want:     "test-value",
		},
		{
			name:     "non-existing env var",
			key:      "NON_EXISTING",
			value:    "",
			fallback: "fallback",
			want:     "fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := GetString(tt.key, tt.fallback); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback int
		want     int
	}{
		{
			name:     "valid integer",
			key:      "TEST_INT",
			value:    "42",
			fallback: 0,
			want:     42,
		},
		{
			name:     "invalid integer",
			key:      "TEST_INVALID_INT",
			value:    "not-a-number",
			fallback: 0,
			want:     0,
		},
		{
			name:     "non-existing env var",
			key:      "NON_EXISTING",
			value:    "",
			fallback: 0,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := GetInt(tt.key, tt.fallback); got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback time.Duration
		want     time.Duration
	}{
		{
			name:     "valid duration",
			key:      "TEST_DURATION",
			value:    "5s",
			fallback: time.Second,
			want:     5 * time.Second,
		},
		{
			name:     "invalid duration",
			key:      "TEST_INVALID_DURATION",
			value:    "not-a-duration",
			fallback: time.Second,
			want:     time.Second,
		},
		{
			name:     "non-existing env var",
			key:      "NON_EXISTING",
			value:    "",
			fallback: time.Second,
			want:     time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := GetDuration(tt.key, tt.fallback); got != tt.want {
				t.Errorf("GetDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
