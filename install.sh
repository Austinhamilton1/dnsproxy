#!/bin/bash

# Install dependencies
go mod tidy

if [ ! -d "./build" ]; then
    mkdir "./build"
fi

go build -o ./build/dnsproxy ./cmd/dnsproxy/

if [ ! -d "./config" ]; then 
    mkdir "./config"
fi

if [ ! -f "./config/blocklist.txt" ]; then 
    wget -O "./config/blocklist.txt" "https://raw.githubusercontent.com/cbuijs/oisd/refs/heads/master/small/domains.adblock"
    sed -i 's/\[/# \[/g' "./config/blocklist.txt"
    sed -i 's/!/#/g' "./config/blocklist.txt"
    sed -i 's/||//g' "./config/blocklist.txt"
    sed -i 's/\^//g' "./config/blocklist.txt"
fi

if [ ! -f "./config/config.toml" ]; then
    touch "./config/config.toml"
    cat > ./config/config.toml <<EOF
[dns]
listen = "127.0.0.1:53"

[cache]
enable = true
cleanup_interval = 1

[upstream]
servers = [
    "1.1.1.1:53"
]

[blocklist]
enable = true
file = "/etc/dnsproxy/blocked.txt"

[log]
level = 2
EOF

fi

sudo mv ./build/dnsproxy /usr/local/bin/
sudo mkdir -p /etc/dnsproxy/
sudo mv ./config/config.toml /etc/dnsproxy/config.toml
sudo mv ./config/blocklist.txt /etc/dnsproxy/blocked.txt
sudo touch /etc/systemd/system/dnsproxy.service

sudo tee /etc/systemd/system/dnsproxy.service > /dev/null <<EOF
[Unit]
Description=Go DNS Proxy
After=network.target

[Service]
Type=simple

ExecStart=/usr/local/bin/dnsproxy -config=/etc/dnsproxy/config.toml

Restart=always
RestartSec=5

User=root

[Install]
WantedBy=multi-user.target
EOF

if sudo systemctl is-active --quiet "systemd-resolved"; then
	sudo systemctl stop "systemd-resolved"
fi

sudo systemctl daemon-reload
sudo systemctl enable --now dnsproxy.service
