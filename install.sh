#!/bin/bash

# Install dependencies
go mod tidy

if [ ! -d "./build" ]; then
    mkdir "./build"
fi

go build -o ./build/dnsproxy ./cmd/dnsproxy/main.go

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
    config="[dns]
listen = "127.0.0.1:53"

[cache]
cleanup_interval = 2

[upstream]
server = "1.1.1.1:53"

[blocklist]
file = "/etc/dnsproxy/blocked.txt"

[log]
level = 2"

    echo "$config" > "./config/config.toml"
fi

mv ./build/dnsproxy /usr/local/bin/
mkdir -p /etc/dnsproxy/
mv ./config/config.toml /etc/dnsproxy/config.toml
mv ./config/blocklist.txt /etc/dnsproxy/blocked.txt
touch /etc/systemd/system/dnsproxy.service

service="[Unit]
Description=Go DNS Proxy
After=network.target

[Service]
Type=simple

ExecStart=/usr/local/bin/dnsproxy -config=/etc/dnsproxy/config.toml

Restart=always
RestartSec=5

User=root

[Install]
WantedBy=multi-user.target"

echo "$service" > ./etc/systemd/system/dnsproxy.service

systemctl daemo-reload