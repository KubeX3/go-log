# 📝 Go-Log

A high-performance, zero-dependency logging utility for Go applications. It automatically manages directory creation and generates structured `.log` files with both ANSI-colored terminal output and clean, plain-text file persistence.

## ✨ Features

* **Dual Output:** Logs beautifully colored messages to the terminal while saving clean, uncolored text to your log files.
* **Auto-Directory Creation:** Automatically creates the log folder structure if it doesn't exist.
* **Five Log Levels:** Built-in support for `INFO`, `WARNING`, `ERROR`, `AUDIT`, `EVENT`, and `DEBUG`.
* **Contextual Tagging:** Attach a `location` or package name to your logs for lightning-fast debugging.
* **Zero-Config Ready:** Works out of the box, or can be customized via Environment Variables.
* **Thread Safe:** Designed to handle concurrent logging in Go routines safely.

---

## 📦 Installation

Install the package using `go get`:

```bash
go get github.com/KubeX3/go-log
```

---

## ⚙️ Configuration

By default, file logging is enabled and logs are saved to `./logs/system.log`.

Configuration Options

| Variable | Description | Default |
| :--- | :--- | :--- |
| `LOG_ENABLED` | Set to `true` to write logs to a file, or `false` to only show them in the console. | `true` |
| `LOG_FILE_PATH` | The relative or absolute path where the log file should be created. | `./logs/system.log` |
| `Environment` | Set to `production` to optimize logging and suppress verbose debug output. | `development` |

### Example `.env` Setup
```env
# Enable or disable file logging (true/false)
LOG_ENABLED="true"

# Custom location for your log files
LOG_FILE_PATH="./src/storage/logs/application.log"

# Environment setting (affects log verbosity)
Environment="production"
```

---

## 🚀 Usage

Import the package into your Go project:

```go
package main

import (
	"github.com/KubeX3/go-log"
)

func main() {
	// 1. Basic Logging
	log.LogInfo("Server successfully started on port 8080", "")
	log.LogEvent("Daily database backup triggered", "")

	// 2. Logging with Context/Location (Highly Recommended)
	// The second argument tags the log with a specific module or function name.
	log.LogWarning("High memory usage detected", "SystemMonitor")
	log.LogError("Failed to connect to Redis cluster", "CacheService")
	log.LogAudit("User password updated successfully", "AuthModule")
}
```

### 💻 Terminal Output (With ANSI Colors)
```txt
[2026-04-02 - 17:05:15] [  INFO   ] - Server successfully started on port 8080
[2026-04-02 - 17:05:15] [  EVENT  ] - Daily database backup triggered
[2026-04-02 - 17:05:16] [ WARNING ] - [SystemMonitor] - High memory usage detected
[2026-04-02 - 17:05:16] [  ERROR  ] - [CacheService] - Failed to connect to Redis cluster
[2026-04-02 - 17:05:17] [  AUDIT  ] - [AuthModule] - User password updated successfully
```

## 🛠️ API Reference

All functions share the same signature: `functionName(message: string, location?: string): void`

| Function | Color | Best Used For |
| :--- | :--- | :--- |
| `logInfo()` | 🟢 Green | Standard system operations and startup messages. |
| `logWarning()` | 🟡 Yellow | Non-critical issues or approaching limits. |
| `logError()` | 🔴 Red | Fatal exceptions and system failures. |
| `logAudit()` | 🔵 Cyan | Security events, logins, and authorization. |
| `logEvent()` | 🟣 Magenta | Business logic milestones and cron jobs. |
| `logDebug()` | ⚪ White | Verbose data for deep troubleshooting. |


---

## 👨‍💻 Development

Want to contribute to the project?

1. Clone the repository:

    ```sh
    git clone https://github.com/KubeX3/go-log.git
    ```

2. Run the test suite:

    ```sh
    go test ./...
    ```

### 📜 License

Designed and developed by <b>KubeX3</b>.

Licensed under the <b>MIT License</b>.