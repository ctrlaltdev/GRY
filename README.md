# Go Redir Yourself

Simple shorturl service

## Run as a container

```sh
# Pull and run the container
docker run -d \
  --name gry \
  -p 3000:3000 \
  -v gry-data:/.GRY \
  ctrlaltdev/gry:latest

# Or using docker compose
version: '3'
services:
  gry:
    image: ctrlaltdev/gry:latest
    container_name: gry
    ports:
      - "3000:3000"
    volumes:
      - gry-data:/.GRY
    restart: unless-stopped

volumes:
  gry-data:
```

## Run as a service

Install the binary
```sh
curl -fSsL ln.0x5f.info/getGRY | sh
```

Create a systemd service file

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

## Configuration

GRY can be configured using environment variables:
The folder used to store the redirects will always be in the user's home directory, and the folder needs be already exist with the correct permissions.
```sh
# Change the port (default: 3000)
export GRY_PORT=8080

# Change the storage location (default: .GRY)
export GRY_FOLDER=storage

# Using with docker run
docker run -d \
  --name gry \
  -p 8080:8080 \
  -v gry-data:/storage \
  -e GRY_PORT=8080 \
  -e GRY_FOLDER=storage \
  ctrlaltdev/gry:latest

# Using with docker compose
version: '3'
services:
  gry:
    image: ctrlaltdev/gry:latest
    container_name: gry
    ports:
      - "8080:8080"
    volumes:
      - gry-data:/storage
    environment:
      - GRY_PORT=8080
      - GRY_FOLDER=storage
    restart: unless-stopped

volumes:
  gry-data:
```

When running as a service, you can add these environment variables to the systemd service file:

```
[Service]
// ... existing service configuration ...
Environment=GRY_PORT=8080
Environment=GRY_FOLDER=storage
// ... rest of service configuration ...
```
