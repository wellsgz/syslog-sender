package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSyslogConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  SyslogConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  16,
				Severity:  6,
				Message:   "test message",
				Hostname:  "test-host",
				Program:   "test-program",
			},
			wantErr: false,
		},
		{
			name: "empty message",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  16,
				Severity:  6,
				Message:   "",
			},
			wantErr: true,
		},
		{
			name: "invalid facility",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  25,
				Severity:  6,
				Message:   "test message",
			},
			wantErr: true,
		},
		{
			name: "invalid severity",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  16,
				Severity:  8,
				Message:   "test message",
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      0,
				Transport: "udp",
				Facility:  16,
				Severity:  6,
				Message:   "test message",
			},
			wantErr: true,
		},
		{
			name: "invalid transport",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "invalid",
				Facility:  16,
				Severity:  6,
				Message:   "test message",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewSyslogClient(tt.config)
			err := client.validateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name     string
		config   SyslogConfig
		contains []string
	}{
		{
			name: "basic message formatting",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  16,
				Severity:  6,
				Message:   "test message",
				Hostname:  "test-host",
				Program:   "test-program",
			},
			contains: []string{"<134>", "test-host", "test-program:", "test message"},
		},
		{
			name: "priority calculation",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  4,
				Severity:  1,
				Message:   "security alert",
				Hostname:  "security-host",
				Program:   "security-app",
			},
			contains: []string{"<33>", "security-host", "security-app:", "security alert"},
		},
		{
			name: "hostname with spaces",
			config: SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  16,
				Severity:  6,
				Message:   "test message",
				Hostname:  "test host name",
				Program:   "test program",
			},
			contains: []string{"<134>", "test-host-name", "test-program:", "test message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewSyslogClient(tt.config)
			message, err := client.formatMessage()
			if err != nil {
				t.Errorf("formatMessage() error = %v", err)
				return
			}

			for _, expected := range tt.contains {
				if !strings.Contains(message, expected) {
					t.Errorf("formatMessage() = %v, expected to contain %v", message, expected)
				}
			}
		})
	}
}

func TestPriorityCalculation(t *testing.T) {
	tests := []struct {
		facility int
		severity int
		expected int
	}{
		{0, 0, 0},    // kernel emergency
		{1, 6, 14},   // user info
		{16, 6, 134}, // local0 info
		{4, 1, 33},   // security alert
		{23, 7, 191}, // local7 debug
	}

	for _, tt := range tests {
		t.Run("priority calculation", func(t *testing.T) {
			config := SyslogConfig{
				Address:   "localhost",
				Port:      514,
				Transport: "udp",
				Facility:  tt.facility,
				Severity:  tt.severity,
				Message:   "test message",
			}

			client := NewSyslogClient(config)
			message, err := client.formatMessage()
			if err != nil {
				t.Errorf("formatMessage() error = %v", err)
				return
			}

			expectedPriority := tt.facility*8 + tt.severity
			expectedPrefix := fmt.Sprintf("<%d>", expectedPriority)
			if !strings.HasPrefix(message, expectedPrefix) {
				t.Errorf("Expected priority %d, got message: %s", expectedPriority, message)
			}
		})
	}
}

func TestNewSyslogClient(t *testing.T) {
	config := SyslogConfig{
		Address:   "localhost",
		Port:      514,
		Transport: "udp",
		Facility:  16,
		Severity:  6,
		Message:   "test message",
	}

	client := NewSyslogClient(config)
	if client == nil {
		t.Error("NewSyslogClient() returned nil")
	}

	if client.config.Address != config.Address {
		t.Errorf("Expected address %s, got %s", config.Address, client.config.Address)
	}
}

func TestMessageFormat(t *testing.T) {
	config := SyslogConfig{
		Address:   "localhost",
		Port:      514,
		Transport: "udp",
		Facility:  16,
		Severity:  6,
		Message:   "test message",
		Hostname:  "test-host",
		Program:   "test-program",
	}

	client := NewSyslogClient(config)
	message, err := client.formatMessage()
	if err != nil {
		t.Errorf("formatMessage() error = %v", err)
		return
	}

	// Check RFC 3164 format: <PRI>TIMESTAMP HOSTNAME TAG: MESSAGE
	// Priority should be at the beginning
	if !strings.HasPrefix(message, "<134>") {
		t.Errorf("Expected priority <134>, got: %s", message)
	}

	// Check that hostname is in the message
	if !strings.Contains(message, "test-host") {
		t.Errorf("Expected hostname 'test-host' in message: %s", message)
	}

	// Check that program is in the message
	if !strings.Contains(message, "test-program:") {
		t.Errorf("Expected program 'test-program:' in message: %s", message)
	}

	// Check that the actual message content is present
	if !strings.Contains(message, "test message") {
		t.Errorf("Expected message content 'test message' in: %s", message)
	}
}

func TestSpaceHandling(t *testing.T) {
	config := SyslogConfig{
		Address:   "localhost",
		Port:      514,
		Transport: "udp",
		Facility:  16,
		Severity:  6,
		Message:   "test message with spaces",
		Hostname:  "host with spaces",
		Program:   "program with spaces",
	}

	client := NewSyslogClient(config)
	message, err := client.formatMessage()
	if err != nil {
		t.Errorf("formatMessage() error = %v", err)
		return
	}

	// Hostname spaces should be replaced with hyphens
	if !strings.Contains(message, "host-with-spaces") {
		t.Errorf("Expected hostname with hyphens, got: %s", message)
	}

	// Program spaces should be replaced with hyphens
	if !strings.Contains(message, "program-with-spaces:") {
		t.Errorf("Expected program with hyphens, got: %s", message)
	}

	// Message spaces should be preserved
	if !strings.Contains(message, "test message with spaces") {
		t.Errorf("Expected message spaces to be preserved, got: %s", message)
	}
}

func TestVersionConstants(t *testing.T) {
	if AppName == "" {
		t.Error("AppName constant is empty")
	}

	if AppVersion == "" {
		t.Error("AppVersion constant is empty")
	}

	if AppAuthor == "" {
		t.Error("AppAuthor constant is empty")
	}

	// Check version format (basic semantic versioning)
	versionParts := strings.Split(AppVersion, ".")
	if len(versionParts) != 3 {
		t.Errorf("Version should follow semantic versioning (x.y.z), got: %s", AppVersion)
	}
}

func BenchmarkFormatMessage(b *testing.B) {
	config := SyslogConfig{
		Address:   "localhost",
		Port:      514,
		Transport: "udp",
		Facility:  16,
		Severity:  6,
		Message:   "benchmark test message",
		Hostname:  "benchmark-host",
		Program:   "benchmark-program",
	}

	client := NewSyslogClient(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.formatMessage()
		if err != nil {
			b.Errorf("formatMessage() error = %v", err)
		}
	}
}

func BenchmarkValidateConfig(b *testing.B) {
	config := SyslogConfig{
		Address:   "localhost",
		Port:      514,
		Transport: "udp",
		Facility:  16,
		Severity:  6,
		Message:   "benchmark test message",
	}

	client := NewSyslogClient(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.validateConfig()
		if err != nil {
			b.Errorf("validateConfig() error = %v", err)
		}
	}
}
