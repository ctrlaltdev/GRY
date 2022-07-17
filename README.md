# Go Redir Yourself

Simple shorturl service

## Install

```sh
curl -fSsL ln.0x.5f.info/getGRY | sh
```

## Run as a service

```
[Unit]
Description=GO REDIR YOURSELF
ConditionPathExists=/usr/local/bin/GRY
After=network.target

[Service]
Type=simple
User=<user>
Group=<user>
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/usr/local/bin
ExecStart=/usr/local/bin/GRY

StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=GRY

[Install]
WantedBy=multi-user.target
```
