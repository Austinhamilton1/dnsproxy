# DNSProxy

A lightweight, high-performance DNS proxy written in Go.

DNSProxy forwards DNS requests from clients to one or more upstream DNS servers while providing a simple configuration system and native Linux service support. It is designed to be easy to deploy on servers, embedded devices, and home networks.

---

## Features

- Lightweight and fast
- Written entirely in Go
- Minimal dependencies
- Configurable upstream DNS servers
- Linux systemd service installation
- Simple configuration file
- Easy deployment

---

## Installation

### Clone the repository

```bash
git clone https://github.com/<username>/dnsproxy.git
cd dnsproxy
```

### Build

```bash
go build -o dnsproxy ./cmd/dnsproxy
```

or use the provided installation script:

```bash
chmod +x install.sh
./install.sh
```

The installer will:

- Build the project
- Install the executable
- Create the configuration directory (if necessary)
- Install a systemd service
- Enable the service

---

# Configuration

The proxy uses a configuration file located at:

```
config/config.toml
```

Example configuration:

```toml
[dns]
listen = "0.0.0.0:53"

[upstream]
servers = [
    "1.1.1.1:53",
    "192.168.0.5:53"
]

[cache]
enable = false

[blocklist]
enable = false

[log]
level = 0
```

---

# Running Manually

To start the proxy directly:

```bash
./dnsproxy -config {config file}
```

or

```bash
./build/dnsproxy -config {config file}
```

The proxy will begin listening for DNS requests using the configured settings.

---

# Running as a Service

If installed using the installation script:

Start the service

```bash
sudo systemctl start dnsproxy
```

Stop the service

```bash
sudo systemctl stop dnsproxy
```

Restart the service

```bash
sudo systemctl restart dnsproxy
```

Check service status

```bash
sudo systemctl status dnsproxy
```

Enable automatic startup

```bash
sudo systemctl enable dnsproxy
```

Disable automatic startup

```bash
sudo systemctl disable dnsproxy
```

View logs

```bash
journalctl -u dnsproxy -f
```

---

# User Guide

## Basic Usage

1. Install DNSProxy.
2. Configure your upstream DNS servers.
3. Start the service.
4. Configure your device or router to use the DNSProxy server as its DNS resolver.

Once configured, all DNS queries sent to the proxy will automatically be forwarded to the configured upstream server(s).

---

## Changing Configuration

1. Edit the configuration file.

```bash
sudo nano config/config.toml
```

2. Save the file.

3. Restart the service.

```bash
sudo systemctl restart dnsproxy
```

The new configuration will now be active.

---

## Updating

Pull the latest source:

```bash
git pull
```

Rebuild:

```bash
go build -o build/dnsproxy ./cmd/dnsproxy
```

Restart the service:

```bash
sudo systemctl restart dnsproxy
```

---

# Troubleshooting

### The service will not start

Check the service status:

```bash
sudo systemctl status dnsproxy
```

### View runtime logs

```bash
journalctl -u dnsproxy -f
```

### Verify port 53 is available

```bash
sudo ss -tulpn | grep :53
```

or

```bash
sudo netstat -tulpn | grep :53
```

If another DNS server (such as `systemd-resolved`, `dnsmasq`, or `named`) is already using port 53, it must be stopped or reconfigured.

---

# Building from Source

Requirements:

- Go 1.22 or newer
- Git

Build:

```bash
go mod tidy
go build ./cmd/dnsproxy
```

---

# Project Structure

```
cmd/
    dnsproxy/       Application entry point

internal/
    ...             Internal packages

config/
    config.toml     Configuration
    blocked.txt     List of blocked hosts 

build/
    dnsproxy        Compiled executable
```

---

# Contributing

Contributions are welcome!

Feel free to submit issues, feature requests, or pull requests.

---

# License

This project is licensed under the MIT License.
