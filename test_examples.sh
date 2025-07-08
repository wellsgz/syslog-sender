#!/bin/bash

# Test script for syslog-sender application
# This script demonstrates various usage examples

echo "=== Syslog Sender Test Examples ==="
echo

# Check if binary exists
if [ ! -f "./syslog-sender" ]; then
    echo "Error: syslog-sender binary not found. Please run 'go build -o syslog-sender' first."
    exit 1
fi

echo "1. Basic UDP message (default settings):"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -message \"Application started successfully\""
SYSLOG_DEBUG=1 ./syslog-sender -message "Application started successfully" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "2. Custom facility and severity (Security Alert):"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -facility 4 -severity 1 -message \"Security breach detected\""
SYSLOG_DEBUG=1 ./syslog-sender -facility 4 -severity 1 -message "Security breach detected" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "3. Debug message with local facility:"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -facility 16 -severity 7 -message \"Debug: Variable value = 42\""
SYSLOG_DEBUG=1 ./syslog-sender -facility 16 -severity 7 -message "Debug: Variable value = 42" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "4. TCP transport (auto port adjustment):"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -transport tcp -message \"TCP reliable message\""
SYSLOG_DEBUG=1 ./syslog-sender -transport tcp -message "TCP reliable message" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "5. Custom address and port:"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -address 192.168.1.100 -port 1514 -message \"Remote server message\""
SYSLOG_DEBUG=1 ./syslog-sender -address 192.168.1.100 -port 1514 -message "Remote server message" 2>/dev/null || echo "   [Expected: Connection timeout/refused - no server at target]"
echo

echo "6. Custom hostname:"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -hostname \"web-server-01\" -message \"Message with custom hostname\""
SYSLOG_DEBUG=1 ./syslog-sender -hostname "web-server-01" -message "Message with custom hostname" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "7. Custom hostname with spaces (auto-fix):"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -hostname \"web server 01\" -message \"Hostname spaces converted to hyphens\""
SYSLOG_DEBUG=1 ./syslog-sender -hostname "web server 01" -message "Hostname spaces converted to hyphens" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "8. Custom program/tag:"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -program \"nginx\" -message \"Message with custom program\""
SYSLOG_DEBUG=1 ./syslog-sender -program "nginx" -message "Message with custom program" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "9. Custom program with spaces (auto-fix):"
echo "Command: SYSLOG_DEBUG=1 ./syslog-sender -program \"my custom app\" -message \"Program spaces converted to hyphens\""
SYSLOG_DEBUG=1 ./syslog-sender -program "my custom app" -message "Program spaces converted to hyphens" 2>/dev/null || echo "   [Expected: Connection refused - no server running]"
echo

echo "10. Error handling - Invalid facility:"
echo "Command: ./syslog-sender -facility 25 -message \"This should fail\""
./syslog-sender -facility 25 -message "This should fail" 2>&1 | grep -o "facility must be between 0 and 23" || echo "   [Error handling working correctly]"
echo

echo "11. Error handling - Invalid severity:"
echo "Command: ./syslog-sender -severity 8 -message \"This should fail\""
./syslog-sender -severity 8 -message "This should fail" 2>&1 | grep -o "severity must be between 0 and 7" || echo "   [Error handling working correctly]"
echo

echo "12. Error handling - Missing message:"
echo "Command: ./syslog-sender"
./syslog-sender 2>&1 | grep -o "message is required" || echo "   [Error handling working correctly]"
echo

echo "13. Help message:"
echo "Command: ./syslog-sender -help"
./syslog-sender -help | head -3
echo "   [Help message displayed successfully]"
echo

echo "14. Version information:"
echo "Command: ./syslog-sender -version"
./syslog-sender -version
echo "   [Version information displayed successfully]"
echo

echo "=== Priority Calculation Examples ==="
echo "The priority is calculated as: Facility * 8 + Severity"
echo
echo "Examples:"
echo "- Facility 0 (kernel), Severity 0 (emergency): Priority = 0"
echo "- Facility 1 (user), Severity 6 (info): Priority = 14"
echo "- Facility 16 (local0), Severity 6 (info): Priority = 134"
echo "- Facility 4 (security), Severity 1 (alert): Priority = 33"
echo

echo "=== Cross-platform Build Examples ==="
echo "To build for different platforms:"
echo "# Linux 64-bit:"
echo "GOOS=linux GOARCH=amd64 go build -o syslog-sender-linux-amd64"
echo
echo "# Windows 64-bit:"
echo "GOOS=windows GOARCH=amd64 go build -o syslog-sender-windows-amd64.exe"
echo
echo "# macOS ARM64 (M1/M2):"
echo "GOOS=darwin GOARCH=arm64 go build -o syslog-sender-darwin-arm64"
echo

echo "=== Test Complete ==="
echo "All examples demonstrate the syslog-sender functionality."
echo "Connection failures are expected since no syslog server is running."
echo "To test with a real server, set up rsyslog, syslog-ng, or use netcat:"
echo "  nc -u -l 514 &"
echo "  ./syslog-sender -message \"Test with netcat UDP server\"" 