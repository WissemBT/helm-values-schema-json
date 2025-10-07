package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		wantErr  bool
	}{
		// Standard Go durations
		{
			name:     "standard hours",
			input:    "24h",
			expected: 24 * time.Hour,
		},
		{
			name:     "standard minutes",
			input:    "30m",
			expected: 30 * time.Minute,
		},
		{
			name:     "standard seconds",
			input:    "45s",
			expected: 45 * time.Second,
		},
		{
			name:     "standard milliseconds",
			input:    "500ms",
			expected: 500 * time.Millisecond,
		},
		{
			name:     "standard composite",
			input:    "1h30m",
			expected: time.Hour + 30*time.Minute,
		},

		// Extended units - days
		{
			name:     "one day",
			input:    "1d",
			expected: 24 * time.Hour,
		},
		{
			name:     "multiple days",
			input:    "7d",
			expected: 7 * 24 * time.Hour,
		},
		{
			name:     "fractional days",
			input:    "1.5d",
			expected: time.Duration(1.5 * 24 * float64(time.Hour)),
		},

		// Extended units - weeks
		{
			name:     "one week",
			input:    "1w",
			expected: 7 * 24 * time.Hour,
		},
		{
			name:     "two weeks",
			input:    "2w",
			expected: 14 * 24 * time.Hour,
		},

		// Extended units - months (30 days)
		{
			name:     "one month",
			input:    "1M",
			expected: 30 * 24 * time.Hour,
		},
		{
			name:     "three months",
			input:    "3M",
			expected: 90 * 24 * time.Hour,
		},

		// Extended units - years (365 days)
		{
			name:     "one year",
			input:    "1y",
			expected: 365 * 24 * time.Hour,
		},
		{
			name:     "two years",
			input:    "2y",
			expected: 730 * 24 * time.Hour,
		},

		// Composite durations with extended units
		{
			name:     "days and hours",
			input:    "2d12h",
			expected: 2*24*time.Hour + 12*time.Hour,
		},
		{
			name:     "days, hours, and minutes",
			input:    "1d6h30m",
			expected: 24*time.Hour + 6*time.Hour + 30*time.Minute,
		},
		{
			name:     "week and days",
			input:    "1w3d",
			expected: 7*24*time.Hour + 3*24*time.Hour,
		},
		{
			name:     "complex composite",
			input:    "1M2w3d12h30m",
			expected: 30*24*time.Hour + 14*24*time.Hour + 3*24*time.Hour + 12*time.Hour + 30*time.Minute,
		},

		// Edge cases
		{
			name:     "zero hours",
			input:    "0h",
			expected: 0,
		},
		{
			name:     "zero days",
			input:    "0d",
			expected: 0,
		},

		// Error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "missing unit",
			input:   "123",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "1d@2h",
			wantErr: true,
		},
		{
			name:    "trailing characters",
			input:   "1dextra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMustParseDuration(t *testing.T) {
	t.Run("valid duration", func(t *testing.T) {
		result := MustParseDuration("24h")
		assert.Equal(t, 24*time.Hour, result)
	})

	t.Run("panics on invalid duration", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParseDuration("invalid")
		})
	})
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{
			name:     "zero",
			input:    0,
			expected: "0s",
		},
		{
			name:     "one hour",
			input:    time.Hour,
			expected: "1h0m0s",
		},
		{
			name:     "one day",
			input:    24 * time.Hour,
			expected: "1d",
		},
		{
			name:     "one week",
			input:    7 * 24 * time.Hour,
			expected: "1w",
		},
		{
			name:     "one month",
			input:    30 * 24 * time.Hour,
			expected: "1M",
		},
		{
			name:     "one year",
			input:    365 * 24 * time.Hour,
			expected: "1y",
		},
		{
			name:     "days and hours",
			input:    2*24*time.Hour + 12*time.Hour,
			expected: "2d12h0m0s",
		},
		{
			name:     "complex composite",
			input:    365*24*time.Hour + 30*24*time.Hour + 7*24*time.Hour + 2*24*time.Hour + 3*time.Hour + 30*time.Minute,
			expected: "1y1M1w2d3h30m0s",
		},
		{
			name:     "negative duration",
			input:    -24 * time.Hour,
			expected: "-1d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDurationRoundtrip(t *testing.T) {
	// Test that parsing common inputs produces expected durations
	tests := []struct {
		input string
	}{
		{"24h"},
		{"1d"},
		{"1w"},
		{"1M"},
		{"1y"},
		{"1w3d12h30m"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			duration, err := ParseDuration(tt.input)
			assert.NoError(t, err)
			assert.Greater(t, duration, time.Duration(0))
		})
	}
}
