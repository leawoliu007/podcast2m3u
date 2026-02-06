# Podcast2m3u

Podcast2m3u is a robust service to manage and convert Podcast feeds into local M3U playlists, specifically designed for MPD (Music Player Daemon) or other players that support M3U files but not RSS feeds.

## Features

-   **M3U Generation**: Converts podcast RSS feeds into M3U playlists compatible with MPD (`#EXTM3U`, `#EXTINF:-1`).
-   **Configuration Manager**: Manage multiple subscriptions via a YAML configuration file.
-   **Global Output Path**: Define a single directory for all playlists; filenames are automatically generated from subscription names.
-   **Daemon Mode**: Runs in the background with a built-in scheduler (cron) to automatically update playlists.
-   **Database Tracking**: Uses a local SQLite database to track state (prevents unnecessary re-writes/updates).
-   **Cross-Platform**: Compiles to a single binary for Linux, Windows, and macOS (AMD64/ARM64/ARMv7).

## Installation

### From Source
```bash
git clone https://github.com/your-repo/podcast2m3u.git
cd podcast2m3u
go build
```

## Configuration (`config.yaml`)

Create a `config.yaml` file in the same directory as the binary.

### Example Configuration

```yaml
global:
  update_interval: "0 * * * *"   # Default schedule (e.g., every hour)
  database_path: "podcast2m3u.db" # Path to SQLite database
  output_path: "/var/lib/mpd/playlists" # Global directory where ALL M3U files will be saved

subscriptions:
  - name: "VOA Learning English"
    url: "https://learningenglish.voanews.com/podcast/?count=20&zoneId=1579"
    # Resulting file: /var/lib/mpd/playlists/VOA Learning English.m3u
    # Uses global schedule (every hour)
  
  - name: "Tech News"
    url: "https://example.com/rss.xml"
    # Resulting file: /var/lib/mpd/playlists/Tech News.m3u
    cron: "*/30 * * * *" # Override: Run every 30 minutes
```

**Note:** The output filename is generated automatically using the subscription `name` plus the `.m3u` extension, saved inside the global `output_path`.

## Usage

### 1. Manual Run (One-time)
Updates all subscriptions defined in `config.yaml` immediately and then exits.
```bash
./podcast2m3u --config config.yaml
```

### 2. Daemon Mode (Service & Web UI)
Starts the service, runs updates according to the scheduled cron jobs, and **launches the Web Management Interface**.

```bash
./podcast2m3u --config config.yaml --daemon
```

By default, the web interface is available at `http://localhost:8080`.
You can change the port in `config.yaml`:
```yaml
global:
  web_server_port: "9090"
```

### 3. Legacy Mode (Single URL)
For backward compatibility, print a single podcast playlist to stdout:
```bash
./podcast2m3u --podcast https://example.com/feed.xml > output.m3u
```
