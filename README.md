# WaveLogStoat
<img width="807" height="485" alt="image" src="https://github.com/user-attachments/assets/a6f0bc9f-597c-4f0c-8405-f737f5688991" />

# WaveLog Transport CLI

A lightweight CLI application that receives QSO data from any Logbook (which can emit ADIFs via UDP) and forwards it to WaveLog. 
This is a minimal implementation focused on the core QSO transportation functionality, perfect for 32-bit systems and resource-constrained environments.

Perfect for vintage-computing or small-ressources. 

## Features

- **UDP Listener**: Receives Logbook-QSO data on port 2333
- **Dual Format Support**: Handles both XML and ADIF formats from Non-ADIF-Conform loggers like N1MM as well as ADIF-Conform ones
- **Data Normalization**: Automatic power unit conversion and band detection
- **WaveLog Integration**: Direct HTTP API communication with WaveLog
- **Lightweight**: Single binary executable, minimal dependencies
- **Cross-Platform**: Compiles for Windows (32-bit/64-bit), Linux, macOS
- **Configuration**: Simple INI file configuration
- **Testing**: Built-in WaveLog connection test

## Quick Start

### Prerequisites

- Go 1.19 or later (for building)
- WaveLog instance with API access

### Building

```bash
# Clone or copy the source code
cd wavelog-transport-cli

# Download dependencies
go mod tidy

# Build for your current platform
go build -o wavelog-transport

# Build for 32-bit Windows
GOOS=windows GOARCH=386 go build -o wavelog-transport.exe

# Build for 64-bit Linux
GOOS=linux GOARCH=amd64 go build -o wavelog-transport-linux

# Build for other platforms as needed
```

### Configuration

Create a `config.ini` file:

```ini
[wavelog]
url = https://your-wavelog-instance.com
api_key = your-wave-log-api-key
station_profile_id = 1
timeout = 5000

[server]
port = 2333
verbose = true
```

#### Configuration Options

**[wavelog] section:**
- `url`: Your WaveLog instance URL
- `api_key`: WaveLog API key (from WaveLog settings)
- `station_profile_id`: Station profile ID from WaveLog
- `timeout`: HTTP request timeout in milliseconds (default: 5000)

**[server] section:**
- `port`: UDP port to listen on (default: 2333)
- `verbose`: Enable verbose logging (default: false)

### Running

```bash
# Run with default config
./wavelog-transport

# Run with custom config file
./wavelog-transport /path/to/config.ini

# Show help
./wavelog-transport --help

# Test WaveLog connection
./wavelog-transport --test
```

### Logger Setup

In your logger, configure the UDP broadcast settings:

1. Set **UDP Server** to `127.0.0.1:2333`
2. Set **UDP Server format** to **ADIF**

## Usage Examples

```bash
# Basic usage
./wavelog-transport

# Test connection to WaveLog
./wavelog-transport --test

# Use custom config file
./wavelog-transport --config /etc/wavelog.ini

# Verbose mode (set in config.ini)
[server]
verbose = true
```

## Supported QSO Formats

### N1MM proprietary XML Format (beta beta beta)
- Automatic detection and parsing
- Converts USB/LSB to SSB for compatibility

### ADIF Format
- Standard ADIF field parsing
- Supports custom ADIF records

### Data Normalization

- **Power Conversion**: Automatically converts kW/mW to Watts
- **Band Detection**: Calculates band from frequency
- **Mode Compatibility**: Converts USB/LSB to SSB for ADIF compatibility

## Logging

The application creates two log outputs:

1. **Console**: Real-time status messages
2. **File**: `wavelog-transport.log` with detailed logging

Log format: `WL-TRANSPORT: YYYY-MM-DD HH:MM:SS.microseconds message`

## Error Handling

- **Port Conflicts**: Clear error messages if port 2333 is blocked
- **Network Errors**: Automatic retry with timeout handling
- **API Errors**: Detailed WaveLog API error reporting
- **Malformed Data**: Graceful handling of invalid XML/ADIF

## Architecture

```
WSJT-X --UDP(2333)--> UDP Listener --> Format Parser --> Data Normalizer --> HTTP Client --> WaveLog API
```

## System Requirements

- **Memory**: ~2-5MB runtime
- **CPU**: Minimal usage when idle
- **Disk**: ~5-10MB executable size
- **Network**: Internet access to WaveLog instance

## Security

- HTTPS/TLS support for WaveLog communication
- No SSL certificate validation (compatible with self-signed certificates)
- API key stored locally in config file

## Troubleshooting

### Common Issues

1. **"Failed to bind to UDP port 2333"**
   - Another application is using the port
   - Stop the conflicting application or change the port

2. **"WaveLog connection failed"**
   - Check WaveLog URL and API key
   - Verify network connectivity
   - Test with `--test` option

3. **"No QSOs received"**
   - Verify your Mainloggers UDP configuration
   - Check firewall settings
   - Ensure verbose logging is enabled

### Debug Mode

Enable verbose logging in config.ini:

```ini
[server]
verbose = true
```

This will show detailed information about received messages and API calls.

## Development

The project is structured as follows:

```
main.go      - Main application entry point and UDP server
parser.go    - XML/ADIF parsing logic
normalizer.go - Data normalization (power, band)
wavelog.go   - WaveLog API client
go.mod       - Go module definition
README.md    - This file
```

## License

This project is based on the WaveLogGate by DJ7NT, rewritten as a minimal CLI implementation.

## Support

For issues related to:
- **WaveLog API**: Consult your WaveLog documentation
- **Your Logger**: Refer to Loggers manual // printed / faxed / $whatever format
- **This CLI**: Check logs and troubleshooting section
