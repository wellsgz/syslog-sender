package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// Version information
const (
	AppName    = "syslog-sender"
	AppVersion = "0.2.0"
	AppAuthor  = "https://github.com/wellsgz/syslog-sender"
)

// SyslogConfig holds the configuration for the syslog message
type SyslogConfig struct {
	Address   string
	Port      int
	Transport string
	Facility  int
	Severity  int
	Message   string
	Hostname  string
	Program   string
}

// SyslogClient handles sending syslog messages
type SyslogClient struct {
	config SyslogConfig
}

// NewSyslogClient creates a new syslog client with the given configuration
func NewSyslogClient(config SyslogConfig) *SyslogClient {
	return &SyslogClient{config: config}
}

// validateConfig validates the syslog configuration
func (s *SyslogClient) validateConfig() error {
	if s.config.Message == "" {
		return fmt.Errorf("message is required")
	}

	if s.config.Facility < 0 || s.config.Facility > 23 {
		return fmt.Errorf("facility must be between 0 and 23")
	}

	if s.config.Severity < 0 || s.config.Severity > 7 {
		return fmt.Errorf("severity must be between 0 and 7")
	}

	if s.config.Port < 1 || s.config.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	transport := strings.ToLower(s.config.Transport)
	if transport != "udp" && transport != "tcp" {
		return fmt.Errorf("transport must be 'udp' or 'tcp'")
	}
	s.config.Transport = transport

	return nil
}

// formatMessage creates a RFC 3164 compliant syslog message
func (s *SyslogClient) formatMessage() (string, error) {
	// Calculate priority: Facility * 8 + Severity
	priority := s.config.Facility*8 + s.config.Severity

	// Get current timestamp in RFC 3164 format
	timestamp := time.Now().Format("Jan  2 15:04:05")

	// Get hostname (use custom hostname if provided, otherwise system hostname)
	var hostname string
	if s.config.Hostname != "" {
		hostname = s.config.Hostname
	} else {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
	}

	// Replace spaces with hyphens in hostname to prevent syslog parsing issues
	hostname = strings.ReplaceAll(hostname, " ", "-")

	// Get program/tag (use custom program if provided, otherwise default)
	var program string
	if s.config.Program != "" {
		program = s.config.Program
	} else {
		program = "syslog-sender"
	}

	// Replace spaces with hyphens in program to prevent syslog parsing issues
	program = strings.ReplaceAll(program, " ", "-")

	// Format: <PRI>TIMESTAMP HOSTNAME TAG: MESSAGE
	message := fmt.Sprintf("<%d>%s %s %s: %s",
		priority, timestamp, hostname, program, s.config.Message)

	return message, nil
}

// SendUDP sends the syslog message using UDP
func (s *SyslogClient) SendUDP(message string) error {
	// Resolve UDP address
	serverAddr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address %s: %v", serverAddr, err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to UDP server: %v", err)
	}
	defer conn.Close()

	// Set write timeout
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// Send message
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send UDP message: %v", err)
	}

	return nil
}

// SendTCP sends the syslog message using TCP
func (s *SyslogClient) SendTCP(message string) error {
	// Create TCP connection
	serverAddr := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	conn, err := net.DialTimeout("tcp", serverAddr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to TCP server: %v", err)
	}
	defer conn.Close()

	// Set write timeout
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// Send message (TCP syslog messages often end with newline)
	messageWithNewline := message + "\n"
	_, err = conn.Write([]byte(messageWithNewline))
	if err != nil {
		return fmt.Errorf("failed to send TCP message: %v", err)
	}

	return nil
}

// Send sends the syslog message using the configured transport
func (s *SyslogClient) Send() error {
	// Validate configuration
	if err := s.validateConfig(); err != nil {
		return err
	}

	// Format message
	message, err := s.formatMessage()
	if err != nil {
		return fmt.Errorf("failed to format message: %v", err)
	}

	// Debug output if enabled
	if os.Getenv("SYSLOG_DEBUG") == "1" {
		fmt.Printf("Debug: Sending message: %s\n", message)
		fmt.Printf("Debug: Target: %s:%d (%s)\n", s.config.Address, s.config.Port, s.config.Transport)
	}

	// Send message based on transport
	switch s.config.Transport {
	case "udp":
		return s.SendUDP(message)
	case "tcp":
		return s.SendTCP(message)
	default:
		return fmt.Errorf("unsupported transport: %s", s.config.Transport)
	}
}

// printUsage prints the usage information
func printUsage() {
	fmt.Printf("Syslog Sender v%s - Send syslog messages to remote servers\n\n", AppVersion)
	fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  %s -message \"Application started\"\n", os.Args[0])
	fmt.Printf("  %s -address 192.168.1.100 -transport tcp -message \"TCP message\"\n", os.Args[0])
	fmt.Printf("  %s -facility 4 -severity 1 -message \"Security alert\"\n", os.Args[0])
	fmt.Printf("  %s -hostname \"custom-host\" -message \"Message with custom hostname\"\n", os.Args[0])
	fmt.Printf("  %s -program \"my-app\" -message \"Message with custom program\"\n", os.Args[0])
	fmt.Printf("\nFacilities (0-23): 0=kernel, 1=user, 2=mail, 3=daemon, 4=security, 16-23=local0-7\n")
	fmt.Printf("Severities (0-7): 0=emergency, 1=alert, 2=critical, 3=error, 4=warning, 5=notice, 6=info, 7=debug\n")
}

func main() {
	// Define command line flags
	var config SyslogConfig
	var showHelp bool
	var showVersion bool

	flag.StringVar(&config.Address, "address", "localhost", "Syslog server address")
	flag.IntVar(&config.Port, "port", 514, "Syslog server port")
	flag.StringVar(&config.Transport, "transport", "udp", "Transport protocol (udp or tcp)")
	flag.IntVar(&config.Facility, "facility", 16, "Syslog facility (0-23)")
	flag.IntVar(&config.Severity, "severity", 6, "Syslog severity (0-7)")
	flag.StringVar(&config.Message, "message", "", "Message to send (required)")
	flag.StringVar(&config.Hostname, "hostname", "", "Custom hostname (default: system hostname)")
	flag.StringVar(&config.Program, "program", "", "Custom program/tag name (default: syslog-sender)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showVersion, "version", false, "Show version information")

	// Custom usage function
	flag.Usage = printUsage

	// Parse command line arguments
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("%s version %s\n", AppName, AppVersion)
		fmt.Printf("Project: %s\n", AppAuthor)
		os.Exit(0)
	}

	// Show help if requested
	if showHelp {
		printUsage()
		os.Exit(0)
	}

	// Check if message is provided
	if config.Message == "" {
		fmt.Fprintf(os.Stderr, "Error: message is required\n\n")
		printUsage()
		os.Exit(1)
	}

	// Adjust default port for TCP if not explicitly set
	if flag.Lookup("port").Value.String() == "514" && strings.ToLower(config.Transport) == "tcp" {
		config.Port = 601
	}

	// Create syslog client
	client := NewSyslogClient(config)

	// Send message
	if err := client.Send(); err != nil {
		log.Fatalf("Failed to send syslog message: %v", err)
	}

	// Success message
	if os.Getenv("SYSLOG_DEBUG") == "1" {
		fmt.Printf("Message sent successfully to %s:%d via %s\n",
			config.Address, config.Port, strings.ToUpper(config.Transport))
	}
}
