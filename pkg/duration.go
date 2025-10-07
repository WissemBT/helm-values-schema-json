package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a duration string with support for extended units
// that Go's time.Duration doesn't natively support:
//   - d = 24 hours (day)
//   - w = 7 days (week)
//   - M = 30 days (month)
//   - y = 365 days (year)
//
// Also supports all standard Go time.Duration units: h, m, s, ms, us, ns
//
// Examples:
//   - "24h" = 24 hours
//   - "1d" = 24 hours
//   - "1w" = 168 hours
//   - "1M" = 720 hours
//   - "1y" = 8760 hours
//   - "2d12h30m" = 60.5 hours
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("empty duration string")
	}

	// Try standard time.ParseDuration first
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Pattern to match duration components: number followed by unit
	pattern := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)(y|M|w|d|h|m|s|ms|us|ns)`)
	matches := pattern.FindAllStringSubmatch(s, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %q", s)
	}

	// Verify the entire string was matched (no invalid characters)
	matchedString := ""
	for _, match := range matches {
		matchedString += match[0]
	}
	if matchedString != s {
		return 0, fmt.Errorf("invalid duration format: %q", s)
	}

	var totalDuration time.Duration

	for _, match := range matches {
		valueStr := match[1]
		unit := match[2]

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration value %q: %w", valueStr, err)
		}

		var unitDuration time.Duration
		switch unit {
		case "y":
			unitDuration = time.Duration(value * 365 * 24 * float64(time.Hour))
		case "M":
			unitDuration = time.Duration(value * 30 * 24 * float64(time.Hour))
		case "w":
			unitDuration = time.Duration(value * 7 * 24 * float64(time.Hour))
		case "d":
			unitDuration = time.Duration(value * 24 * float64(time.Hour))
		case "h":
			unitDuration = time.Duration(value * float64(time.Hour))
		case "m":
			unitDuration = time.Duration(value * float64(time.Minute))
		case "s":
			unitDuration = time.Duration(value * float64(time.Second))
		case "ms":
			unitDuration = time.Duration(value * float64(time.Millisecond))
		case "us":
			unitDuration = time.Duration(value * float64(time.Microsecond))
		case "ns":
			unitDuration = time.Duration(value)
		default:
			return 0, fmt.Errorf("unknown duration unit: %q", unit)
		}

		totalDuration += unitDuration
	}

	return totalDuration, nil
}

// MustParseDuration is like ParseDuration but panics on error.
// Useful for static initialization.
func MustParseDuration(s string) time.Duration {
	d, err := ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

// FormatDuration formats a duration using the most appropriate unit(s).
// It will use extended units (y, M, w, d) when appropriate.
func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	var parts []string

	// Handle negative durations
	negative := d < 0
	if negative {
		d = -d
	}

	// Years
	if years := d / (365 * 24 * time.Hour); years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", years))
		d %= 365 * 24 * time.Hour
	}

	// Months (30 days)
	if months := d / (30 * 24 * time.Hour); months > 0 {
		parts = append(parts, fmt.Sprintf("%dM", months))
		d %= 30 * 24 * time.Hour
	}

	// Weeks
	if weeks := d / (7 * 24 * time.Hour); weeks > 0 {
		parts = append(parts, fmt.Sprintf("%dw", weeks))
		d %= 7 * 24 * time.Hour
	}

	// Days
	if days := d / (24 * time.Hour); days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
		d %= 24 * time.Hour
	}

	// Use standard Go formatting for remaining time
	if d > 0 {
		parts = append(parts, d.String())
	}

	result := strings.Join(parts, "")
	if negative {
		result = "-" + result
	}

	return result
}
