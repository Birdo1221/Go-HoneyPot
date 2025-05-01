# Go-HoneyPot: SSH & SMB Honeypots in Go
### Please Consider using fail2ban as additional protection / cover


This code was built to capitalize [AbuseIPDB API](https://www.abuseipdb.com/user/137416).  
> <sup>— This [link] shows reports submitted by me on AbuseIPDB. Feel free to check it out if you're curious or just browsing for fun.</sup>

This repository contains two security honeypot implementations written in Go:
- **SSH Honeypot** (`sshpot.go`) - Detects brute force attempts against SSH servers
- **SMB Honeypot** (`smbpot.go`) - Detects unauthorized SMB access attempts

Both honeypots automatically report malicious IPs to AbuseIPDB and implement temporary bans using iptables.

## Features

### SSH Honeypot
- Listens on multiple common SSH ports (2222, 2200, 22, etc.)
- Simulates OpenSSH server behavior
- Logs all connection attempts
- Reports to AbuseIPDB (category 18 - SSH, 22 - Brute Force)
- 30-minute IP bans using iptables
- Prevents duplicate reports within 15 minutes

### SMB Honeypot
- Listens on SMB port (445)
- Simulates SMB protocol negotiation
- Detailed logging to `smb_attempts.log`
- Reports to AbuseIPDB (categories 14 - SMB, 15 - Brute Force, 18 - SSH)
- 30-minute IP bans using iptables
- Dedicated log file with timestamps

## Installation

### Prerequisites
- Linux system
- Go 1.16+ (to build)
- iptables (for IP banning)
- curl (for AbuseIPDB reporting)
- AbuseIPDB API key (free tier available)

### Build Instructions
```bash
git clone https://github.com/Birdo1221/Go-HoneyPot.git
cd Go-HoneyPot
go build sshpot.go
go build smbpot.go
```

Run with root privileges (required for iptables):

```bash
# Run SSH honeypot
sudo ./sshpot

# Run SMB honeypot
sudo ./smbpot
```

Logging + Warnings
These honeypots are designed to attract malicious traffic.

```bash
SSH Honeypot: Outputs to stdout/stderr
SMB Honeypot: Logs to smb_attempts.log in the same directory

|  Run on a dedicated server or VM
|  Ensure your system is properly secured before deploying
|  Monitor resource usage
|  Dont run on systems with sensitive data
|  Use firewall rules to restrict access if needed
```
