# Lightarr
Lightarr is a home automation tool that connects your Plex server to your Wiz light bulbs. Add this 
lightweight and simple to use app to your *arr collection and enjoy!

# Motivation
One of the most annoying parts of a movie night is deciding who has to get up and turn all 
the lights off. But now, whether it's with your family, significant other, friends or even a solo 
movie night, you can leave it to Lightarr! Simply install it on your home server or PC and configure 
it to your liking and never worry about getting up right before you start the movie.

## Quick Start
First get your Plex token: [Find X-Plex-Token](https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/)

### Docker (Recommended)

```bash
# Pull image
docker pull aradd7/lightarr:latest

# Create docker-compose.yml and .env, then run
docker compose up -d
```

**docker-compose.yml:**
```yaml
services:
  lightarr:
    image: aradd7/lightarr:latest
    container_name: lightarr
    environment:
      - PORT=${PORT:-10100}
      - X_PLEX_TOKEN=${X_PLEX_TOKEN}
      - WIZ_SUBNET=${WIZ_SUBNET:-192.168.1.0/24}
      - X_PLEX_CLIENT_IDENTIFIER=lightarr-app
      - DB_PATH=/data/lightarr.db
    volumes:
      - ./data:/data
    ports:
      - "${PORT:-10100}:${PORT:-10100}"
    restart: unless-stopped
```

**.env:**
```env
# Your Plex authentication token (required)
X_PLEX_TOKEN=your_token_here

# Optional: customize port (defaults to 10100, change if needed)
PORT=10100

# Optional: network subnet (defaults to 192.168.1.0/24)
WIZ_SUBNET=192.168.1.0/24
```

### Plex Setup
1. Go to Plex Settings → Webhooks
2. Add webhook: `http://<lightarr-ip>:10100/plexhook`

## Using the Web UI

Access the interface at `http://<lightarr-ip>:10100` (replace with your server's IP and port)

### Initial Setup

1. **Head to Plex tab**
    - Click "Plex Accounts" to fetch and save users from your Plex server
    - Click "Plex Devices" to fetch and save playback devices (TVs, phones, etc.)
2. **Head to Bulbs tab**
    - See discovered Wiz bulbs and give them friendly names

### Creating Rules

Rules automate your lights based on who's watching, what device they're using, and what's happening.

**Step 1: Head to Rules tab and click Add Rule**
- **Account(s)**: Which Plex user(s) trigger this rule (e.g., your account, family members)
- **Device(s)**: Which playback device(s) trigger this rule (e.g., Living Room TV, Bedroom Chromecast)
- **Event(s)**: Which playback events trigger this rule (pause, stop, resume, play)

**Step 2: Add Actions (Command/Bulb Pairs)**

Each rule can control multiple bulbs with different commands. Click "Add Command/Bulb Pair" to create actions:

- **Commands**:
  - **Turn Off**: Power off selected bulbs
  - **Turn On**: Power on selected bulbs
  - **Dim**: Set brightness (1-100%)
  - **Change Color**: Set RGB color (0-255 per channel)
  - **Change Temperature**: Set warmth (2200-6500K)

- **Select Bulbs**: Choose which bulbs this command applies to
  - Each bulb can only appear in one pair per rule
  - No overlap between pairs

**Example Rule:**
```
Account: Your Account
Device: Living Room TV
Event: media.play

Actions:
  Pair 1: Turn Off → Kitchen Light, Hallway Light
  Pair 2: Dim 15% → Couch Light, Floor Lamp
  Pair 3: Change Temperature 2700K → Ambient Light
```

When you press play on the Living Room TV:
- Kitchen and hallway lights turn off
- Couch light and floor lamp dim to 15%
- Ambient light sets to warm 2700K

## Features
- **Plex Integration**: Responds to media playback events via webhooks
- **Rule-Based Automation**: Configure which users, devices, and events trigger which lights
- **Auto-Discovery**: Scans network subnet to find Wiz bulbs
- **Web UI**: React frontend for managing rules
- **Multi-Platform**: Docker images for amd64/arm64

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `10100` | HTTP server port |
| `X_PLEX_TOKEN` | *required* | Plex authentication token |
| `WIZ_SUBNET` | `192.168.1.0/24` | Network subnet to scan for bulbs |
| `DB_PATH` | `/data/lightarr.db` | SQLite database path |

### Development

```bash
# Clone repo
git clone https://github.com/yourusername/lightarr.git
cd lightarr

# Set up environment
cp .env.example .env
nano .env  # Fill in X_PLEX_TOKEN

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Build frontend
cd frontend && npm run build && cd ..

# Run locally
go run .
```

## Tech Stack

- **Backend**: Go 1.21+
- **Frontend**: React + Vite + TanStack Query
- **Database**: SQLite + Goose migrations
- **Deployment**: Docker multi-platform (amd64/arm64)
- **CI/CD**: GitHub Actions → Docker Hub

## Network Requirements

- Lightarr uses standard Docker bridge networking
- Port 10100/TCP for webhooks and web UI
- UDP communication to Wiz bulbs (port 38899)
- Subnet scanning discovers bulbs without broadcast
