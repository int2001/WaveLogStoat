package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func normalizeQSO(qso QSO) QSO {
	// Normalize power
	qso.POWER = normalizePower(qso.POWER)

	// Calculate band from frequency
	if qso.FREQ != "" {
		qso.BAND = calculateBand(qso.FREQ)
	}

	return qso
}

func normalizePower(powerStr string) string {
	if powerStr == "" {
		return powerStr
	}

	// Remove whitespace and convert to lowercase
	powerStr = strings.TrimSpace(strings.ToLower(powerStr))

	// Extract numeric value
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)`)
	match := re.FindStringSubmatch(powerStr)
	if len(match) < 2 {
		return powerStr // Return original if no valid number found
	}

	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return powerStr // Return original if parsing fails
	}

	// Convert based on unit
	if strings.Contains(powerStr, "kw") {
		value *= 1000 // kW to W
	} else if strings.Contains(powerStr, "mw") {
		value *= 0.001 // mW to W
	}
	// If it's just 'w' or no unit, assume it's already in watts

	// Return as string without decimal if it's a whole number
	if value == math.Floor(value) {
		return fmt.Sprintf("%.0f", value)
	}
	return fmt.Sprintf("%.3f", value)
}

func calculateBand(freqStr string) string {
	freq, err := strconv.ParseFloat(freqStr, 64)
	if err != nil {
		return ""
	}

	// Band definitions (frequencies in MHz)
	// These are standard amateur radio bands
	bandMap := []struct {
		name  string
		lower float64
		upper float64
	}{
		{"160M", 1.800, 2.000},
		{"80M", 3.500, 4.000},
		{"60M", 5.330, 5.400},
		{"40M", 7.000, 7.300},
		{"30M", 10.100, 10.150},
		{"20M", 14.000, 14.350},
		{"17M", 18.068, 18.168},
		{"15M", 21.000, 21.450},
		{"12M", 24.890, 24.990},
		{"10M", 28.000, 29.700},
		{"6M", 50.000, 54.000},
		{"2M", 144.000, 148.000},
		{"1.25M", 222.000, 225.000},
		{"70CM", 420.000, 450.000},
		{"33CM", 902.000, 928.000},
		{"23CM", 1240.000, 1300.000},
	}

	for _, band := range bandMap {
		if freq >= band.lower && freq <= band.upper {
			return band.name
		}
	}

	return ""
}