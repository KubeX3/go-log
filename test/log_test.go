package test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/KubeX3/go-log"
	"github.com/KubeX3/go-log/internal/utils"
)

// Helper to capture stdout for console verification
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

// Sets up a temporary environment for the logger
func setupTestConfig(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		t.Fatal(err)
	}

	logPath := filepath.Join(tmpDir, "test.log")

	// Reset Mock DOTENV config
	utils.DOTENV.LogEnabled = true
	utils.DOTENV.LogFilePath = logPath
	utils.DOTENV.Environment = "development"

	return logPath
}

func TestLogger(t *testing.T) {
	logPath := setupTestConfig(t)
	defer os.RemoveAll(filepath.Dir(logPath))

	// Regex for timestamp: [YYYY-MM-DD - HH:MM:SS]
	dateRegex := `\[\d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2}\]`

	// --- Standard Logging Tests ---

	t.Run("LogError should log message and file content", func(t *testing.T) {
		msg := "Critical error"
		output := captureStdout(func() { log.LogError(msg) })

		if !strings.Contains(output, "[  ERROR  ]") {
			t.Errorf("Console missing ERROR tag")
		}
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  ERROR  ] - "+msg) {
			t.Errorf("File missing error message")
		}
	})

	t.Run("LogWarning should include location", func(t *testing.T) {
		msg := "Low memory"
		loc := "Worker-1"
		captureStdout(func() { log.LogWarning(msg, loc) })

		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "- ["+loc+"] - "+msg) {
			t.Errorf("File missing location metadata")
		}
	})

	t.Run("LogInfo should log correctly", func(t *testing.T) {
		msg := "System online"
		captureStdout(func() { log.LogInfo(msg) })
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  INFO   ] - "+msg) {
			t.Error("Info log format incorrect")
		}
	})

	t.Run("LogAudit should log correctly", func(t *testing.T) {
		msg := "User updated profile"
		captureStdout(func() { log.LogAudit(msg, "AdminPanel") })
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  AUDIT  ] - [AdminPanel] - "+msg) {
			t.Error("Audit log format incorrect")
		}
	})

	t.Run("LogEvent should log correctly", func(t *testing.T) {
		msg := "Job finished"
		captureStdout(func() { log.LogEvent(msg) })
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "[  EVENT  ] - "+msg) {
			t.Error("Event log format incorrect")
		}
	})

	// --- Formatted (...F) Logging Tests ---

	t.Run("LogInfoF should interpolate values", func(t *testing.T) {
		host := "localhost"
		port := 443
		expected := "Starting HTTPS Server on localhost:443"
		
		output := captureStdout(func() {
			log.LogInfoF("Starting HTTPS Server on %s:%d", host, port)
		})

		if !strings.Contains(output, expected) {
			t.Errorf("Formatted console output mismatch. Got: %s", output)
		}
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), expected) {
			t.Error("Formatted file content mismatch")
		}
	})

	t.Run("LogErrorF should interpolate correctly", func(t *testing.T) {
		code := 502
		expected := "Bad Gateway: 502"
		captureStdout(func() { log.LogErrorF("Bad Gateway: %d", code) })
		
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), expected) {
			t.Error("LogErrorF failed to write formatted string to file")
		}
	})

	t.Run("LogWarningF should interpolate correctly", func(t *testing.T) {
		expected := "Disk at 95%"
		captureStdout(func() { log.LogWarningF("Disk at %d%%", 95) })
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), expected) {
			t.Error("LogWarningF format mismatch")
		}
	})

	t.Run("LogAuditF should interpolate correctly", func(t *testing.T) {
		user := "JaneDoe"
		captureStdout(func() { log.LogAuditF("User %s deleted record", user) })
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "User JaneDoe deleted record") {
			t.Error("LogAuditF format mismatch")
		}
	})

	t.Run("LogEventF should interpolate correctly", func(t *testing.T) {
		captureStdout(func() { log.LogEventF("Batch %s processed", "A1") })
		content, _ := os.ReadFile(logPath)
		if !strings.Contains(string(content), "Batch A1 processed") {
			t.Error("LogEventF format mismatch")
		}
	})

	// --- Debug Logic Tests ---

	t.Run("LogDebug and LogDebugF should show in development", func(t *testing.T) {
		utils.DOTENV.Environment = "development"
		out1 := captureStdout(func() { log.LogDebug("Dev mode active") })
		out2 := captureStdout(func() { log.LogDebugF("Pointer: %p", &t) })

		if out1 == "" || out2 == "" {
			t.Error("Debug logs should not be empty in development")
		}
	})

	t.Run("LogDebug and LogDebugF should NOT show in production", func(t *testing.T) {
		utils.DOTENV.Environment = "production"
		out1 := captureStdout(func() { log.LogDebug("Secret") })
		out2 := captureStdout(func() { log.LogDebugF("Secret %d", 123) })

		if out1 != "" || out2 != "" {
			t.Errorf("Debug logs should be suppressed in production")
		}
	})

	t.Run("File format should match date regex", func(t *testing.T) {
		content, _ := os.ReadFile(logPath)
		// Check the last line for correct timestamp format
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		lastLine := lines[len(lines)-1]
		
		matched, _ := regexp.MatchString("^"+dateRegex, lastLine)
		if !matched {
			t.Errorf("Log timestamp format mismatch. Got: %s", lastLine)
		}
	})

	// --- Fatal Logging Tests ---

    t.Run("LogFatal should exit with code 1 and log message", func(t *testing.T) {
        // We run the test in a subprocess because LogFatal calls os.Exit(1)
        if os.Getenv("BE_CRASHER") == "1" {
            log.LogFatal("Critical failure")
            return
        }

        // Spawn a subprocess of the current test
        cmd := exec.Command(os.Args[0], "-test.run=TestLogger/LogFatal_should_exit_with_code_1_and_log_message")
        cmd.Env = append(os.Environ(), "BE_CRASHER=1")
        err := cmd.Run()

        // Check if the exit code was 1
        if e, ok := err.(*exec.ExitError); ok && !e.Success() {
            if e.ExitCode() != 1 {
                t.Errorf("Expected exit status 1, got %d", e.ExitCode())
            }
        } else {
            t.Fatalf("Process ran successfully, but it should have exited with status 1")
        }
    })

    t.Run("LogFatalF should exit with code 1 and log formatted message", func(t *testing.T) {
        if os.Getenv("BE_CRASHER_F") == "1" {
            log.LogFatalF("Fatal error code: %d", 500)
            return
        }

        cmd := exec.Command(os.Args[0], "-test.run=TestLogger/LogFatalF_should_exit_with_code_1_and_log_formatted_message")
        cmd.Env = append(os.Environ(), "BE_CRASHER_F=1")
        err := cmd.Run()

        if e, ok := err.(*exec.ExitError); ok && !e.Success() {
            if e.ExitCode() != 1 {
                t.Errorf("Expected exit status 1, got %d", e.ExitCode())
            }
        } else {
            t.Fatalf("Process ran successfully, but it should have exited with status 1")
        }
    })
}