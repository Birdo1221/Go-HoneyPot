# Go-HoneyPot: Lightweight SSH & SMB Honeypots in Go


> <sup> ⚠️ Note: This tool provides basic detection of abusive scanning and unauthorized login attempts. For enhanced protection, it is strongly recommended to use **Fail2Ban** or other **trusted** security solutions alongside this tool. </sup>

> <sup> This code was built to capitalize [AbuseIPDB API](https://www.abuseipdb.com/user/137416). </sup>  
> <sup> — Reports I’ve submitted to AbuseIPDB. Check them out if you're curious or just looking for an example.</sup>

This repository contains two security honeypot implementations written in Go:
- **SSH Honeypot** (**`sshpot.go`**) - Detects brute force attempts against SSH servers
- **SMB Honeypot** ( **`smbpot.go`**) - Detects unauthorized SMB access attempts

> Both honeypots automatically report the IP/(S) to AbuseIPDB and temporary bans using **`iptables`**.

 # Features

> ## SSH Honeypot

```
- Listens on multiple common SSH ports (2222, 2200, 22, etc.)
- Simulates OpenSSH server behavior
- Logs all connection attempts
- Reports to AbuseIPDB (category 18 - SSH, 22 - Brute Force)
- 30-minute IP bans using iptables
- Prevents duplicate reports within 15 minutes

```
>## SMB Honeypot
```
- Listens on SMB port (445)
- Simulates SMB protocol negotiation
- Detailed logging to `smb_attempts.log`
- Reports to AbuseIPDB (categories 14 - SMB, 15 - Brute Force, 18 - SSH)
- 30-minute IP bans using iptables
- Dedicated log file with timestamps
```

> ## Installation
## Prerequisites

> <sup> Before you begin, ensure you have the following dependencies installed: </sup>
```
- **Linux system**
- **Go 1.16+** (required to build the project)
- **iptables** (used for IP banning)
- **curl** (for reporting to AbuseIPDB)
- **AbuseIPDB API key** (Free tier available, sign up on their website)
```

> ## Build Instructions
> ```bash
> git clone https://github.com/Birdo1221/Go-HoneyPot.git
> cd Go-HoneyPot
> go build sshpot.go + go build smbpot.go
> ```

> ## Running with root privileges (required for iptables):

> ```bash
> # Run SSH + SMB honeypot
> sudo ./sshpot
> sudo ./smbpot
> ```

> ## Running with As a background process:

> ```bash
> # Run SSH Honeypot in the background
> nohup sudo ./sshpot &
> 
> # Run SMB Honeypot in the background
> nohup sudo ./smbpot & 
> ```

<sup> An alternative of nohup would be [**`screen`**](https://www.geeksforgeeks.org/screen-command-in-linux-with-examples/).  if you would like to install that and use that instead   Logging and User Warnings:
  These honeypots are designed to attract malicious traffic by actively attempting to capture unauthorized login credentials and related activity.</sup>


> ```bash
> SSH Honeypot: Outputs to stdout/stderr
> SMB Honeypot: Logs to smb_attempts.log in the same directory
> 
> |  Run on a dedicated server or VM
> |  Ensure your system is properly secured before deploying
> |  Monitor resource usage
> |  Dont run on systems with sensitive data
> |  Use firewall rules to restrict access if needed
> ```
