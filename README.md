# Syslog Sender Application

A cross-platform command-line application written in Go for sending syslog messages to remote syslog servers.

## Features

- **Cross-platform**: Single binary executable for Linux, macOS, and Windows
- **Multiple transports**: Support for UDP and TCP protocols
- **Configurable parameters**: Full control over syslog message components
- **RFC 3164 compliant**: Follows standard syslog message format
- **Easy to use**: Simple command-line interface

## Requirements

### Functional Requirements

1. **Message Transmission**: Send syslog messages to a specified syslog server
2. **Configurable Parameters**:
   - **Address**: Target syslog server hostname or IP address
   - **Port**: Target port number (default: 514 for UDP, 601 for TCP)
   - **Transport**: Protocol selection (UDP or TCP)
   - **Facility**: Syslog facility code (0-23)
   - **Severity**: Syslog severity level (0-7)
   - **Message**: The actual log message content
   - **Hostname**: Custom hostname for the syslog message (default: system hostname)
   - **Program**: Custom program/tag name for the syslog message (default: syslog-sender)

3. **Cross-platform Compatibility**: Single Go binary that runs on multiple operating systems
4. **Error Handling**: Proper error reporting for network and parameter issues

### Non-functional Requirements

- **Performance**: Low latency message sending
- **Reliability**: TCP support for guaranteed delivery
- **Usability**: Intuitive command-line interface
- **Maintainability**: Clean, well-documented code structure

## Syslog Facilities

| Code | Facility |
|------|----------|
| 0    | kernel messages |
| 1    | user-level messages |
| 2    | mail system |
| 3    | system daemons |
| 4    | security/authorization messages |
| 5    | messages generated internally by syslogd |
| 6    | line printer subsystem |
| 7    | network news subsystem |
| 8    | UUCP subsystem |
| 9    | clock daemon |
| 10   | security/authorization messages |
| 11   | FTP daemon |
| 12   | NTP subsystem |
| 13   | log audit |
| 14   | log alert |
| 15   | clock daemon |
| 16-23| local use facilities (local0-local7) |

## Syslog Severity Levels

| Code | Severity | Description |
|------|----------|-------------|
| 0    | Emergency | System is unusable |
| 1    | Alert | Action must be taken immediately |
| 2    | Critical | Critical conditions |
| 3    | Error | Error conditions |
| 4    | Warning | Warning conditions |
| 5    | Notice | Normal but significant condition |
| 6    | Informational | Informational messages |
| 7    | Debug | Debug-level messages |

## Installation

### Prerequisites

- Go 1.19 or later (for building from source)

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd syslog-sender

# Build for current platform
go build -o syslog-sender

# Build for multiple platforms
# Linux
GOOS=linux GOARCH=amd64 go build -o syslog-sender-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o syslog-sender-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o syslog-sender-windows-amd64.exe
```

## Usage

### Command Line Options

```bash
syslog-sender [OPTIONS]

Options:
  -address string
        Syslog server address (default "localhost")
  -port int
        Syslog server port (default 514)
  -transport string
        Transport protocol: udp or tcp (default "udp")
  -facility int
        Syslog facility (0-23) (default 16)
  -severity int
        Syslog severity (0-7) (default 6)
  -message string
        Message to send (required)
  -hostname string
        Custom hostname (default: system hostname)
  -program string
        Custom program/tag name (default: syslog-sender)
  -help
        Show help message
  -version
        Show version information
```

### Examples

#### Basic Usage (UDP)

```bash
# Send a simple informational message
./syslog-sender -message "Application started successfully"

# Send to specific server
./syslog-sender -address "192.168.1.100" -message "Remote log message"
```

#### TCP Transport

```bash
# Send using TCP for reliability
./syslog-sender -transport tcp -port 601 -message "Important system event"
```

#### Custom Facility and Severity

```bash
# Send security alert (facility 4, severity 1)
./syslog-sender -facility 4 -severity 1 -message "Security breach detected"

# Send debug message (facility 16, severity 7)
./syslog-sender -facility 16 -severity 7 -message "Debug: Variable value = 42"
```

#### Custom Hostname

```bash
# Send message with custom hostname
./syslog-sender -hostname "web-server-01" -message "Application deployment completed"

# Send from simulated host
./syslog-sender -hostname "production-db" -facility 3 -severity 6 -message "Database backup completed"

# Hostname with spaces (automatically converted to hyphens)
./syslog-sender -hostname "web server 01" -message "Spaces in hostname handled automatically"
```

#### Custom Program/Tag

```bash
# Send message with custom program name
./syslog-sender -program "nginx" -message "HTTP server started"

# Send message from custom application
./syslog-sender -program "my-custom-app" -facility 16 -severity 6 -message "Application event logged"

# Program with spaces (automatically converted to hyphens)
./syslog-sender -program "my application" -message "Spaces in program name handled automatically"
```

#### Complete Configuration

```bash
./syslog-sender \
  -address "syslog.example.com" \
  -port 514 \
  -transport udp \
  -facility 16 \
  -severity 4 \
  -hostname "custom-host" \
  -program "my-app" \
  -message "Custom configuration test message"
```

## Implementation Details

### Architecture

The application consists of several key components:

1. **Command Line Parser**: Handles argument parsing and validation
2. **Syslog Client**: Core functionality for message formatting and transmission
3. **Transport Layer**: Abstracts UDP and TCP connections
4. **Message Formatter**: Creates RFC 3164 compliant syslog messages

### Message Format

The application generates syslog messages following RFC 3164 format:

```
<PRI>TIMESTAMP HOSTNAME TAG: MESSAGE
```

Where:
- **PRI**: Priority value calculated as (Facility × 8 + Severity)
- **TIMESTAMP**: Current timestamp in RFC 3164 format
- **HOSTNAME**: Custom or system hostname (spaces automatically converted to hyphens)
- **TAG**: Custom or default program name (spaces automatically converted to hyphens)
- **MESSAGE**: User-provided message content

### Space Handling

To maintain RFC 3164 compliance and prevent parsing issues:
- **Hostname**: Any spaces in custom hostnames are automatically replaced with hyphens
- **Program**: Any spaces in custom program names are automatically replaced with hyphens
- **Message**: Spaces in messages are preserved as-is

### Error Handling

The application provides comprehensive error handling for:

- Invalid command line arguments
- Network connection failures
- DNS resolution errors
- Transport-specific errors
- Message formatting issues

### Dependencies

- Standard Go library only (no external dependencies)
- Uses `net` package for network operations
- Uses `flag` package for command line parsing
- Uses `log/syslog` package for message formatting

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Integration Testing

```bash
# Test with local syslog server (rsyslog/syslog-ng)
./syslog-sender -message "Test message"

# Test with netcat as UDP server
nc -u -l 514 &
./syslog-sender -message "UDP test message"

# Test with netcat as TCP server
nc -l 601 &
./syslog-sender -transport tcp -port 601 -message "TCP test message"
```

## Troubleshooting

### Common Issues

1. **Permission Denied (Port < 1024)**
   - Use `sudo` for privileged ports on Unix systems
   - Or use non-privileged ports (> 1024)

2. **Connection Refused**
   - Verify target server is running and listening
   - Check firewall settings
   - Confirm port number and transport protocol

3. **DNS Resolution Errors**
   - Use IP address instead of hostname
   - Check DNS configuration

### Debug Mode

Enable verbose logging by setting environment variable:

```bash
export SYSLOG_DEBUG=1
./syslog-sender -message "Debug test"
```

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Changelog

### Version 0.1.0 (2025-07-08)

#### Features
- **Cross-platform syslog client**: Send syslog messages to remote servers via UDP or TCP
- **RFC 3164 compliance**: Properly formatted syslog messages with priority calculation
- **Configurable parameters**: All syslog components are user-configurable
  - Server address and port
  - Transport protocol (UDP/TCP with auto port adjustment)
  - Facility codes (0-23)
  - Severity levels (0-7)
  - Custom hostname (with automatic space handling)
  - Custom program/tag name (with automatic space handling)
- **Automatic space handling**: Spaces in hostname and program fields automatically converted to hyphens for RFC compliance
- **Debug mode**: `SYSLOG_DEBUG=1` environment variable for message inspection
- **Comprehensive validation**: Input parameter validation with helpful error messages
- **Network timeouts**: Connection and write timeouts to prevent hanging
- **Version information**: `-version` flag to display version details

#### Implementation Details
- **Language**: Go (Golang) for cross-platform compatibility
- **Dependencies**: Standard library only, no external dependencies
- **Message format**: `<PRI>TIMESTAMP HOSTNAME TAG: MESSAGE`
- **Priority calculation**: `Facility × 8 + Severity`
- **Default ports**: UDP 514, TCP 601
- **Error handling**: Comprehensive error reporting for all failure scenarios

#### Cross-platform Support
- Linux (x86_64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (x86_64)
- Single binary deployment for each platform

## Security Considerations

- Messages are sent in plain text
- No authentication mechanism
- Consider using TLS for sensitive environments
- Validate input parameters to prevent injection attacks 