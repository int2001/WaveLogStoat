package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func sendToWaveLog(adifString string, qso QSO) error {
	// Prepare payload
	payload := WaveLogPayload{
		Key:             config.WaveLog.APIKey,
		StationProfileID: config.WaveLog.StationProfileID,
		Type:            "adif",
		String:          adifString,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Prepare request URL
	apiURL := strings.TrimSuffix(config.WaveLog.URL, "/") + "/api/qso"

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WL-Transport-v1.0")

	// Create HTTP client with timeout
	timeout := time.Duration(config.WaveLog.Timeout) * time.Millisecond
	client := &http.Client{
		Timeout: timeout,
	}

	if verbose {
		logger.Printf("Sending QSO to WaveLog: %s on %s", qso.CALL, qso.FREQ)
		logger.Printf("API URL: %s", apiURL)
		logger.Printf("Payload: %s", string(jsonData))
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var waveLogResponse WaveLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&waveLogResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Check response status
	if waveLogResponse.Status == "created" {
		logger.Printf("✓ QSO successfully added: %s on %s MHz", qso.CALL, qso.FREQ)
	} else {
		var errorMsg string
		if len(waveLogResponse.Messages) > 0 {
			errorMsg = strings.Join(waveLogResponse.Messages, ", ")
		}
		return fmt.Errorf("QSO not added (status: %s): %s", waveLogResponse.Status, errorMsg)
	}

	return nil
}

// Test function to verify WaveLog connectivity
func testWaveLogConnection() error {
	// Create a test ADIF record
	testADIF := `<ADIF_VER:5>5.0<EOH>
<TEST_CALL:6>K0TEST<QSO_DATE:8>20240101<TIME_ON:6>120000<MODE:4>FT8<FREQ:6>14.074<BAND:3>20M<EOR>`

	// Prepare payload
	payload := WaveLogPayload{
		Key:             config.WaveLog.APIKey,
		StationProfileID: config.WaveLog.StationProfileID,
		Type:            "adif",
		String:          testADIF,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Prepare request URL (use dry run endpoint if available)
	apiURL := strings.TrimSuffix(config.WaveLog.URL, "/") + "/api/qso"

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WL-Transport-v1.0-Test")

	// Create HTTP client with timeout
	timeout := time.Duration(config.WaveLog.Timeout) * time.Millisecond
	client := &http.Client{
		Timeout: timeout,
	}

	logger.Printf("Testing WaveLog connection to: %s", apiURL)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var waveLogResponse WaveLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&waveLogResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	logger.Printf("WaveLog connection test - Status: %d, Response: %s", resp.StatusCode, waveLogResponse.Status)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		logger.Printf("✓ WaveLog connection successful")
		return nil
	}

	return fmt.Errorf("WaveLog connection failed: HTTP %d - %s", resp.StatusCode, waveLogResponse.Status)
}