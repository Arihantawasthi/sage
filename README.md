# Sage - The Process Manager
> A UNIX-based process manager daemon and CLI tool for managing my own hosted projects.

## ⚙️ Overview

**Sage** is a lightweight daemon and CLI-based process manager built in Go. It allows you to:

- Start, stop and list processes (local services)
- Monitor CPU usage, memory and uptime
- Communicate via a custom binary protocol (SPMP) over a UNIX socket

**NOTE**: Since this project was mainly for learning processes, I've kept the external dependencies to minimum. I'm only using `gopsutils` as an external dependency.

## ✨ Features:
- **Daemonized Process Manager**
- **Custom Protocol (SPMP)**"
- **In-memory Process Store with Monitoring**
- **Real-Time saving of logs from the processes**
- **CLI Tool to manage processes (local services)**

---

## 📡 SPMP — Sage Process Management Protocol

**SPMP** is a custom binary protocol designed specifically for fast and reliable communication between the `sagectl` CLI and the daemon over a UNIX socket. Here's how it works:

### 🔧 Packet Format
Each SPMP packet has the following structure:

| Field          | Size (bytes) | Description                           |
|----------------|--------------|---------------------------------------|
| Magic Bytes    | 2            | Fixed: `"SG"` to validate packet      |
| Version        | 1            | Protocol version (currently `0x1`)    |
| Encoding       | 2            | `TX` (text) or `JS` (JSON)            |
| Type           | 1            | Command type (`START`, `STOP`, etc.)  |
| Payload Size   | 4            | Size of the payload in bytes          |
| Payload        | N            | The actual payload                    |

### 📤 Supported Commands (Types)

- `TypeStart (0x01)` — Start a single service or all
- `TypeStop  (0x02)` — Stop a service
- `TypeList  (0x03)` — Get running services

### 🧬 Encodings

- `TEXT` (TX): Payload is plain service name or keyword (`"all"`)
- `JSON` (JS): Used in responses or structured payloads

---

## 🗃️ Architecture
```
              +------------------+       +-----------------+
              |     sagectl      | <---> |   /tmp/sage.sock|
              +------------------+       +--------+--------+
                                                     |
                                            +--------v--------+
                                            |     SAGE Daemon  |
                                            |------------------|
                                            | Config Loader    |
                                            | SPMP Server      |
                                            | Process Manager  |
                                            +--------+--------+
                                                     |
                               +---------------------+----------------------+
                               |                                            |
                     +---------v----------+                     +-----------v----------+
                     |  gopsutil Monitor  |                     |   stdout/stderr Log   |
                     +--------------------+                     +-----------------------+
```

---

## 🚀 Usage

### Example Config: `config.json`
```json
{
  "serviceMap": {
    "redis": {
      "name": "redis",
      "command": "/usr/bin/redis-server",
      "args": ["--port", "6379"],
      "workingDir": "/usr/local/bin",
      "env": {
        "ENVIRONMENT": "development"
      }
    }
  }
}
```

### Build & Start the Daemon

To use SAGE, you’ll need to build both the **daemon** and the **CLI tool (`sagectl`)**.
### 🛠️ Build the Daemon & CLI

```bash
go build -o sage ./cmd/sage
go build -o sagectl ./cmd/sagectl
```

Place the binaries somewhere accessible in your `$PATH`, or use absolute paths in your systemd unit.

### 🗃️ Optional: Install to `/usr/local/bin`

```bash
sudo cp sage /usr/local/bin/
sudo cp sagectl /usr/local/bin/
```

# 🔧 Setting up systemd service for the daemon
#### Create a unit file at /etc/systemd/system/sage.service
```sudo tee /etc/systemd/system/sage.service > /dev/null <<EOF
[Unit]
Description=SAGE Daemon
After=network.target

[Service]
ExecStart=/usr/local/bin/sage
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target
```

#### Reload systemd to pick up the new unit
```bash
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
```

#### Enable and start the service
```bash
sudo systemctl enable sage
sudo systemctl start sage
```

#### Check status
```bash
systemctl status sage
```

---

### 🚀 Start the Daemon

You can start the daemon using `systemd`.

```bash
systemctl start sage
```

---

### 🧪 Run CLI Commands

```bash
sagectl start redis
sagectl stop redis
sagectl start all
sagectl list
```

---

## 🖥️ Output Sample

```bash
$ sagectl list
+--------+---------+-----------+--------+----------------------+---------+--------+
| SNo.   | PID     | P_NAME    | NAME   | CMD                  | CPU%    | MEM%   |
+--------+---------+-----------+--------+----------------------+---------+--------+
| 1      | 12345   | redis     | redis  | /usr/bin/redis       | 0.2     | 1.1    |
+--------+---------+-----------+--------+----------------------+---------+--------+
```

---

## 💬 Final Thoughts

This project is an exercise in **system programming**, **protocol design**, and **observability**. It's simple, but clean, and gets the job done.
Now that phase 1 is complete, HTTP support and advanced service orchestration will follow in Phase 2.
