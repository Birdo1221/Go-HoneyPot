package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	AbuseIPDBAPIKey = "Temp_Key" // Replace with your actual AbuseIPDB API key
	BanDuration     = 30 * time.Minute
)

var (
	ports          = []int{2222, 2200, 22222, 50000, 3389, 1337, 10001, 222, 2022, 2181, 23, 2000, 830, 2002, 5353, 8081, 6000, 5900}
	reportedIPs    = make(map[string]time.Time)
	reportedIPsMux = &sync.Mutex{}
)

func main() {
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			startSSHHoneypot(p)
		}(port)
	}

	wg.Wait()
}

func startSSHHoneypot(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			log.Printf("Port %d is already in use. Skipping...\n", port)
			return
		}
		log.Printf("Failed to start server on port %d: %v\n", port, err)
		return
	}
	defer listener.Close()

	log.Printf("Starting SSH honeypot on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection on port %d: %v\n", port, err)
			continue
		}

		go handleConnection(conn, port)
	}
}

func handleConnection(conn net.Conn, port int) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	log.Printf("Connection attempt from %s on port %d\n", remoteAddr, port)

	// Simulate SSH banner
	conn.Write([]byte("SSH-2.0-OpenSSH_7.9p1 Debian-10+deb10u2\r\n"))

	// Read client SSH banner
	reader := bufio.NewReader(conn)
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading from %s: %v\n", remoteAddr, err)
	}

	// Simulate authentication request
	conn.Write([]byte("SSH-2.0-OpenSSH_7.9p1 Debian-10+deb10u2\r\n"))

	// Record the attempt
	recordAttempt(remoteAddr, port)

	// Report to AbuseIPDB
	reportToAbuseIPDB(remoteAddr)

	// Ban the IP
	banIP(remoteAddr)

	// Send a fake message and close
	conn.Write([]byte("Permission denied (publickey,password).\r\n"))
}

func recordAttempt(ip string, port int) {
	log.Printf("SSH login attempt detected from %s on port %d\n", ip, port)
}

func reportToAbuseIPDB(ip string) {
	reportedIPsMux.Lock()
	defer reportedIPsMux.Unlock()

	// Check if we've reported this IP recently (within 15 minutes)
	if lastReported, exists := reportedIPs[ip]; exists && time.Since(lastReported) < 15*time.Minute {
		log.Printf("Skipping report for IP %s as it was reported recently\n", ip)
		return
	}

	// Prepare the curl command to report to AbuseIPDB
	cmd := exec.Command("curl", "https://api.abuseipdb.com/api/v2/report",
		"--data-urlencode", fmt.Sprintf("ip=%s", ip),
		"-d", "categories=18,22",
		"--data-urlencode", "comment=[Birdo SSH Honeypot] SSH login attempt",
		"-H", fmt.Sprintf("Key: %s", AbuseIPDBAPIKey),
		"-H", "Accept: application/json")

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to report IP %s to AbuseIPDB: %v\n", ip, err)
		return
	}

	reportedIPs[ip] = time.Now()
	log.Printf("Reported IP %s to AbuseIPDB successfully\n", ip)
}

func banIP(ip string) {
	// Ban the IP using iptables
	banCmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	err := banCmd.Run()
	if err != nil {
		log.Printf("Failed to ban IP %s: %v\n", ip, err)
		return
	}

	log.Printf("Banned IP %s successfully\n", ip)

	// Schedule unban after BanDuration
	time.AfterFunc(BanDuration, func() {
		unbanCmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
		err := unbanCmd.Run()
		if err != nil {
			log.Printf("Failed to unban IP %s: %v\n", ip, err)
			return
		}
		log.Printf("Unbanned IP %s successfully\n", ip)
	})
}
