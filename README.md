# 🍯 Go-HoneyPot

> Lightweight SSH & SMB Honeypots written in Go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.16%2B-blue)](https://golang.org/)
[![AbuseIPDB](https://img.shields.io/badge/Protected%20by-AbuseIPDB-green)](https://www.abuseipdb.com/user/137416)

<p align="center">
  <img src="https://raw.githubusercontent.com/Birdo1221/Go-HoneyPot/assets/honeypot-logo.png" alt="Go-HoneyPot Logo" width="200" height="200">
</p>

## 📋 Overview

Go-HoneyPot provides lightweight security monitoring tools that detect and report malicious scanning and unauthorized access attempts. The project consists of two honeypot implementations:

- **SSH Honeypot** - Detects brute force attempts against SSH servers
- **SMB Honeypot** - Identifies unauthorized SMB access attempts

Both honeypots automatically report offending IPs to [AbuseIPDB](https://www.abuseipdb.com/) and temporarily ban them using `iptables`.

> ⚠️ **Security Note**: This tool provides basic detection of abusive scanning and unauthorized login attempts. For enhanced protection, it is strongly recommended to use **Fail2Ban** or other **trusted** security solutions alongside this tool.

## ✨ Features

### SSH Honeypot (`sshpot.go`)

| Feature | Description |
|---------|-------------|
| 🔌 Multi-port Listening | Monitors multiple common SSH ports (2222, 2200, 22, etc.) |
| 🛡️ Server Simulation | Accurately mimics OpenSSH server behavior |
| 📝 Comprehensive Logging | Records all connection attempts with timestamps |
| 🚫 Automated Reporting | Reports to AbuseIPDB (category 18 - SSH, 22 - Brute Force) |
| ⏱️ Temporary Banning | Implements 30-minute IP bans using iptables |
| 🔄 Duplicate Prevention | Prevents duplicate reports within 15 minutes |

### SMB Honeypot (`smbpot.go`) 

| Feature | Description |
|---------|-------------|
| 🔌 Protocol Monitoring | Listens on standard SMB port (445) |
| 🛡️ Protocol Simulation | Simulates SMB protocol negotiation |
| 📝 Detailed Logging | Writes comprehensive logs to `smb_attempts.log` |
| 🚫 Automated Reporting | Reports to AbuseIPDB (categories 14 - SMB, 15 - Brute Force) |
| ⏱️ Temporary Banning | Implements 30-minute IP bans using iptables |
| 📊 Attack Analytics | Tracks attack patterns and frequency |

## 🔧 Installation

### Prerequisites

Before installation, ensure you have the following dependencies installed:

- **Linux system** (Debian/Ubuntu recommended)
- **Go 1.16+** (required to build the project)
- **iptables** (used for IP banning)
- **curl** (for reporting to AbuseIPDB)
- **AbuseIPDB API key** (Free tier available, [sign up here](https://www.abuseipdb.com/))

### Build Instructions

```bash
# Clone the repository
git clone https://github.com/Birdo1221/Go-HoneyPot.git

# Navigate to project directory
cd Go-HoneyPot

# Build both honeypots
go build sshpot.go
go build smbpot.go
```

### Configuration

Before running, you'll need to set your AbuseIPDB API key:

```bash
# Edit the configuration (example - actual implementation may vary)
nano config.json

# Set your API key in the configuration file
{
  "abuseipdb_api_key": "YOUR_API_KEY_HERE"
}
```

## 🚀 Usage

### Running with Root Privileges

Both honeypots require root privileges to manage iptables:

```bash
# Run SSH honeypot
sudo ./sshpot

# Run SMB honeypot
sudo ./smbpot
```

### Running as Background Processes

To run the honeypots in the background:

```bash
# Run SSH Honeypot in the background
nohup sudo ./sshpot &

# Run SMB Honeypot in the background
nohup sudo ./smbpot &
```

> 💡 **Tip**: As an alternative to `nohup`, you can use [`screen`](https://www.geeksforgeeks.org/screen-command-in-linux-with-examples/) or [`tmux`](https://github.com/tmux/tmux/wiki) for better session management.

### Logging

- **SSH Honeypot**: Outputs to stdout/stderr by default
- **SMB Honeypot**: Logs to `smb_attempts.log` in the same directory

## ⚠️ Important Warnings

These honeypots are designed to attract malicious traffic by simulating public services. Please consider the following:

- **Dedicated Environment**: Run on a dedicated server or VM, not your primary system
- **System Security**: Ensure your system is properly secured before deploying
- **Resource Monitoring**: Keep an eye on system resources, especially when under attack
- **Data Protection**: Never run on systems with sensitive data
- **Access Control**: Use firewall rules to restrict access if needed

## 📊 Statistics

The honeypots collect statistics on attack attempts which can be viewed in the log files. For a visual representation of attacks reported through your account, visit your [AbuseIPDB dashboard](https://www.abuseipdb.com/user/137416).

## 🛠️ Advanced Configuration

### Custom Port Configuration

To modify listening ports, edit the source code:

```go
// In sshpot.go
var sshPorts = []int{22, 2222, 2200}  // Modify as needed
```

### Adjusting Ban Duration

To change the default 30-minute ban duration:

```go
// In both sshpot.go and smbpot.go
const banDuration = 30  // Time in minutes, modify as needed
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📬 Contact

Project Link: [https://github.com/Birdo1221/Go-HoneyPot](https://github.com/Birdo1221/Go-HoneyPot)

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/Birdo1221">Birdo1221</a>
</p>
