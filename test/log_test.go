package test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/KubeX3/go-log/internal/utils"
	"github.com/KubeX3/go-log/cmd/package"
)

// Helper to capture stdout
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func setupTestConfig(t *testing.T) string {
	// Create a temp directory for logs
	tmpDir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		t.Fatal(err)
	}

	logPath := filepath.Join(tmpDir, "test.log")

	// Mock the DOTENV config
	utils.DOTENV.LogEnabled = true
	utils.DOTENV.LogFilePath = logPath
	utils.DOTENV.Environment = "development"

	return logPath
}

func TestLogger(t *testing.T) {
	logPath := setupTestConfig(t)
	defer os.RemoveAll(filepath.Dir(logPath))

	// Regex for [YYYY-MM-DD - HH:MM:SS]
	dateRegex := `\[\d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2}\]`

	t.Run("should correctly log an ERROR", func(t *testing.T) {
		message := "Database failed"
		output := captureStdout(func() {
			log.LogError(message)
		})

		// Verify Console Output (contains centered tag)
		if !strings.Contains(output, "[  ERROR  ]") {
			t.Errorf("Expected console output to contain [  ERROR  ], got %s", output)
		}

		// Verify File Content
		content, _ := os.ReadFile(logPath)
		matched, _ := regexp.MatchString(dateRegex+` \[  ERROR  \] - `+message, string(content))
		if !matched {
			t.Errorf("File log format mismatch. Got: %s", string(content))
		}
	})

	t.Run("should correctly log a WARNING with location", func(t *testing.T) {
		message := "Disk low"
		location := "Server"
		output := captureStdout(func() {
			log.LogWarning(message, location)
		})

		if !strings.Contains(output, "[ WARNING ]") {
			t.Errorf("Expected console output to contain [ WARNING ]")
		}

		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "- [Server] - Disk low") {
			t.Errorf("File log missing location metadata. Got: %s", string(content))
		}
	})

	t.Run("should correctly log an INFO message", func(t *testing.T) {
		message := "App started"
		captureStdout(func() {
			log.LogInfo(message)
		})

		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  INFO   ] - "+message) {
			t.Errorf("Info log format incorrect")
		}
	})

	t.Run("should correctly log an AUDIT message", func(t *testing.T) {
		message := "User logged in"
		captureStdout(func() {
			log.LogAudit(message, "Auth")
		})

		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  AUDIT  ] - [Auth] - "+message) {
			t.Errorf("Audit log format incorrect")
		}
	})

	t.Run("should correctly log an EVENT message", func(t *testing.T) {
		message := "Backup completed"
		captureStdout(func() {
			log.LogEvent(message)
		})

		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  EVENT  ] - "+message) {
			t.Errorf("Event log format incorrect")
		}
	})

	t.Run("should correctly log a DEBUG message in development", func(t *testing.T) {
		utils.DOTENV.Environment = "development"
		message := "Before DB start"
		output := captureStdout(func() {
			log.LogDebug(message)
		})

		if !strings.Contains(output, "[  DEBUG  ]") {
			t.Errorf("Debug log should be visible in development")
		}
	})

	t.Run("should NOT log DEBUG message in production", func(t *testing.T) {
		utils.DOTENV.Environment = "production"
		output := captureStdout(func() {
			log.LogDebug("Secret stuff")
		})

		if output != "" {
			t.Errorf("Expected no output for debug in production, got %s", output)
		}
	})
}