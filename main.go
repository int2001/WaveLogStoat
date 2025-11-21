package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

// Configuration structure
type Config struct {
	WaveLog struct {
		URL              string `ini:"url"`
		APIKey           string `ini:"api_key"`
		StationProfileID string `ini:"station_profile_id"`
		Timeout          int    `ini:"timeout"`
	} `ini:"wavelog"`
	Server struct {
		Port    int  `ini:"port"`
		Verbose bool `ini:"verbose"`
	} `ini:"server"`
}

// WaveLog API payload structure
type WaveLogPayload struct {
	Key             string `json:"key"`
	StationProfileID string `json:"station_profile_id"`
	Type            string `json:"type"`
	String          string `json:"string"`
}

// WaveLog API response structure
type WaveLogResponse struct {
	Status   string   `json:"status"`
	Messages []string `json:"messages,omitempty"`
}

// QSO structure for internal processing
type QSO struct {
	CALL             string
	MODE             string
	QSO_DATE_OFF     string
	QSO_DATE         string
	TIME_OFF         string
	TIME_ON          string
	RST_RCVD         string
	RST_SENT         string
	FREQ             string
	FREQ_RX          string
	OPERATOR         string
	COMMENT          string
	POWER            string
	STX              string
	SRX              string
	STX_STRING       string
	SRX_STRING       string
	RTX              string
	MYCALL           string
	GRIDSQUARE       string
	MY_GRIDSQUARE    string
	STATION_CALLSIGN string
	BAND             string
	NAME             string
	QTH              string
	STATE            string
	COUNTRY          string
	CQZ              string
	ITUZ             string
	CONT             string
	IOTA             string
	DXCC             string
	PROP_MODE        string
	SAT_NAME         string
	SAT_MODE         string
	// Contest-specific fields (ADIF compliant only)
	CONTEST_ID       string
	PREFIX           string
	// Additional WaveLog-supported fields
	SUBMODE          string
	QSLMSG           string
	NOTES            string
	EMAIL            string
	DARC_DOK         string
	SOTA_REF         string
	WWFF_REF         string
	POTA_REF         string
	CNTY             string
	REGION           string
	LAT              string
	LON              string
	ANT_AZ           string
	ANT_EL           string
	ANT_PATH         string
	A_INDEX          string
	K_INDEX          string
	SFI              string
	RX_PWR           string
	Created          bool
	Fail             interface{}
}

const (
	AppName    = "WavelogStoat"
	AppVersion = "0.0.2"
)

var (
	config   Config
	verbose  bool
	logFile  *os.File
	logger   *log.Logger
)

func init() {
	// Initialize logging
	logFile, err := os.OpenFile("wavelog-transport.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(io.MultiWriter(os.Stdout, logFile), "WL-TRANSPORT: ", log.LstdFlags|log.Lmicroseconds)
}

func main() {
	// Parse command line arguments
	configFile := "config.ini"
	testMode := false

	for i, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			printUsage()
			return
		} else if arg == "--test" || arg == "-t" {
			testMode = true
		} else if arg == "--config" || arg == "-c" {
			if i+1 < len(os.Args) {
				configFile = os.Args[i+1]
			}
		} else if i == 1 && !strings.HasPrefix(arg, "-") {
			configFile = arg
		}
	}

	// Load configuration
	if err := loadConfig(configFile); err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	verbose = config.Server.Verbose

	if testMode {
		logger.Printf("Running in test mode")
		if err := testWaveLogConnection(); err != nil {
			logger.Fatalf("WaveLog connection test failed: %v", err)
		}
		logger.Printf("WaveLog connection test passed")
		return
	}

	logger.Printf("Starting WaveLog Transport CLI on port %d", config.Server.Port)

	// Start UDP server
	if err := startUDPServer(); err != nil {
		logger.Fatalf("Failed to start UDP server: %v", err)
	}
}

func printUsage() {
	fmt.Println("WaveLog Transport CLI - Lightweight QSO transport from WSJT-X to WaveLog")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  wavelog-transport [options] [config.ini]")
	fmt.Println("  wavelog-transport --help")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("  -t, --test           Test WaveLog connection")
	fmt.Println("  -c, --config FILE    Use specified config file")
	fmt.Println("")
	fmt.Println("Default config file: config.ini")
	fmt.Println("")
	fmt.Println("Example config.ini:")
	fmt.Println("[wavelog]")
	fmt.Println("url = https://wavelog.example.com")
	fmt.Println("api_key = your-api-key")
	fmt.Println("station_profile_id = 1")
	fmt.Println("timeout = 5000")
	fmt.Println("")
	fmt.Println("[server]")
	fmt.Println("port = 2333")
	fmt.Println("verbose = true")
}

func loadConfig(filename string) error {
	// Set default values
	config.WaveLog.Timeout = 5000
	config.Server.Port = 2333
	config.Server.Verbose = false

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create default config file
		logger.Printf("Creating default config file: %s", filename)
		if err := createDefaultConfig(filename); err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
		logger.Printf("Please edit %s with your WaveLog settings and restart", filename)
		return fmt.Errorf("default config created - please configure and restart")
	}

	cfg, err := ini.Load(filename)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	if err := cfg.MapTo(&config); err != nil {
		return fmt.Errorf("failed to map config: %v", err)
	}

	// Validate required settings
	if config.WaveLog.URL == "" || config.WaveLog.APIKey == "" || config.WaveLog.StationProfileID == "" {
		return fmt.Errorf("missing required WaveLog configuration (url, api_key, station_profile_id)")
	}

	return nil
}

func createDefaultConfig(filename string) error {
	cfg := ini.Empty()

	wavelogSec := cfg.Section("wavelog")
	wavelogSec.Key("url").SetValue("https://your-wavelog-url.com")
	wavelogSec.Key("api_key").SetValue("your-api-key-here")
	wavelogSec.Key("station_profile_id").SetValue("1")
	wavelogSec.Key("timeout").SetValue("5000")

	serverSec := cfg.Section("server")
	serverSec.Key("port").SetValue("2333")
	serverSec.Key("verbose").SetValue("true")

	return cfg.SaveTo(filename)
}

func startUDPServer() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind to UDP port %d: %v", config.Server.Port, err)
	}
	defer conn.Close()

	logger.Printf("UDP server listening on port %d", config.Server.Port)

	buffer := make([]byte, 4096)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := string(buffer[:n])
		logger.Printf("Received %d bytes from %s", n, clientAddr.String())

		if verbose {
			logger.Printf("Message content: %s", message)
		}

		// Process the message asynchronously
		go processMessage(message)
	}
}

func processMessage(message string) {
	// Detect format and parse
	if strings.Contains(message, "xml") {
		// XML format typically contains single QSO
		processSingleQSO(message, true)
	} else {
		// ADIF format - check for multiple QSOs separated by <EOR>
		if strings.Contains(message, "<EOR>") {
			processMultipleQSOs(message)
		} else {
			// Single QSO without explicit <EOR> tag
			processSingleQSO(message, false)
		}
	}
}

func processMultipleQSOs(adifPayload string) {
	// Split by <EOR> and process each QSO
	// Note: Keep the <EOR> tag for proper ADIF parsing
	qsoRecords := strings.Split(adifPayload, "<EOR>")

	processedCount := 0
	for i, qsoRecord := range qsoRecords {
		// Skip empty records (last element might be empty after split)
		qsoRecord = strings.TrimSpace(qsoRecord)
		if qsoRecord == "" {
			continue
		}

		// Add back the <EOR> tag for proper parsing (except for last record)
		if i < len(qsoRecords)-1 {
			qsoRecord += "<EOR>"
		}

		if verbose {
			logger.Printf("Processing QSO %d of %d", processedCount+1, len(qsoRecords)-1)
		}

		success := processSingleQSO(qsoRecord, false)
		if success {
			processedCount++
		}
	}

	if processedCount > 1 {
		logger.Printf("Successfully processed %d QSOs from batch payload", processedCount)
	}
}

func processSingleQSO(message string, isXML bool) bool {
	var qso QSO
	var err error

	// Parse the QSO
	if isXML {
		qso, err = parseXMLMessage(message)
	} else {
		qso, err = parseADIFMessage(message)
	}

	if err != nil {
		logger.Printf("Failed to parse message: %v", err)
		return false
	}

	// Normalize data
	qso = normalizeQSO(qso)

	// Generate ADIF string
	adifString := generateADIF(qso)

	// Send to WaveLog
	if err := sendToWaveLog(adifString, qso); err != nil {
		logger.Printf("Failed to send QSO to WaveLog: %v", err)
		return false
	}

	return true
}
