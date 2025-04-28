package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	AbuseIPDBAPIKey = "Temp_Key" // Replace with your actual AbuseIPDB API key
	SMBPort         = 445
	BanDuration     = 30 * time.Minute
	ReportInterval  = 15 * time.Minute
)

var (
	reportedIPs    = make(map[string]time.Time)
	reportedIPsMux = &sync.Mutex{}
	logFile        *os.File
)

func main() {
	// Set up logging
	err := setupLogging()
	if err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}
	defer logFile.Close()

	log.Println("Starting SMB honeypot...")
	startSMBHoneypot(SMBPort)
}

func setupLogging() error {
	var err error
	logFile, err = os.OpenFile("smb_attempts.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)
	return nil
}

func startSMBHoneypot(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			log.Printf("Port %d is already in use. Skipping...\n", port)
			return
		}
		log.Fatalf("Failed to start server on port %d: %v\n", port, err)
	}
	defer listener.Close()

	log.Printf("SMB honeypot listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	log.Printf("Connection attempt from %s\n", remoteAddr)

	// Simulate SMB protocol negotiation
	_, err := conn.Write([]byte("\x00\x00\x00\x00\xff\x53\x4d\x42\x72\x00\x00\x00\x00"))
	if err != nil {
		log.Printf("Error writing SMB response to %s: %v\n", remoteAddr, err)
		return
	}

	// Try to read data (simulate credentials being sent)
	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading from %s: %v\n", remoteAddr, err)
	} else {
		log.Printf("Received data from %s: %q\n", remoteAddr, strings.TrimSpace(data))
	}

	// Report and ban the IP
	reportToAbuseIPDB(remoteAddr)
	go banIP(remoteAddr)
}

func reportToAbuseIPDB(ip string) {
	reportedIPsMux.Lock()
	defer reportedIPsMux.Unlock()

	// Check if we've reported this IP recently
	if lastReported, exists := reportedIPs[ip]; exists && time.Since(lastReported) < ReportInterval {
		log.Printf("Skipping report for IP %s (reported recently)\n", ip)
		return
	}

	log.Printf("Reporting IP %s to AbuseIPDB\n", ip)
	cmd := exec.Command("curl", "https://api.abuseipdb.com/api/v2/report",
		"--data-urlencode", fmt.Sprintf("ip=%s", ip),
		"-d", "categories=18,14,15",
		"--data-urlencode", "comment=[Birdo SMB Honeypot] SMB unauthorized attempt",
		"-H", fmt.Sprintf("Key: %s", AbuseIPDBAPIKey),
		"-H", "Accept: application/json")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to report IP %s: %v\nOutput: %s\n", ip, err, string(output))
		return
	}

	reportedIPs[ip] = time.Now()
	log.Printf("Successfully reported IP %s to AbuseIPDB\n", ip)
}

func banIP(ip string) {
	log.Printf("Banning IP %s\n", ip)

	// Ban the IP using iptables
	banCmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	if output, err := banCmd.CombinedOutput(); err != nil {
		log.Printf("Failed to ban IP %s: %v\nOutput: %s\n", ip, err, string(output))
		return
	}

	log.Printf("Successfully banned IP %s\n", ip)

	// Schedule unban after BanDuration
	time.AfterFunc(BanDuration, func() {
		log.Printf("Unbanning IP %s\n", ip)
		unbanCmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
		if output, err := unbanCmd.CombinedOutput(); err != nil {
			log.Printf("Failed to unban IP %s: %v\nOutput: %s\n", ip, err, string(output))
			return
		}
		log.Printf("Successfully unbanned IP %s\n", ip)
	})
}
