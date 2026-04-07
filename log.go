package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KubeX3/go-log/internal/utils"
	"github.com/KubeX3/go-log/internal/utils/enums"
)

// ANSI color codes
const (
	ColorReset   = "\x1b[0m"
	ColorRed     = "\x1b[31m"
	ColorGreen   = "\x1b[32m"
	ColorYellow  = "\x1b[33m"
	ColorMagenta = "\x1b[35m"
	ColorCyan    = "\x1b[36m"
	ColorGray    = "\x1b[90m"
)

/**
 * Appends the log message to file synchronously (simulated via standard IO)
 */
func logFile(logWithColor string, logWithoutColor string) {
	logEnabled := utils.DOTENV.LogEnabled

	if logEnabled {
		filePath := utils.DOTENV.LogFilePath
		dir := filepath.Dir(filePath)

		// Synchronous directory creation (recursive)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
		}

		// Synchronous file append
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error handling file operation: %v\n", err)
		} else {
			defer f.Close()
			if _, err := f.WriteString(logWithoutColor + "\n"); err != nil {
				fmt.Printf("Error writing to log file: %v\n", err)
			}
		}
	}
	fmt.Println(logWithColor)
}

// Format date & time
func getFormattedDateTime() string {
	now := time.Now()
	return fmt.Sprintf("[%d-%02d-%02d - %02d:%02d:%02d]",
		now.Year(), int(now.Month()), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

/**
 * Centers text inside brackets to a fixed width of 9 characters.
 */
func getPaddedType(logType string) string {
	width := 9
	length := len(logType)

	if length >= width {
		return "[" + logType + "]"
	}

	leftPadding := (width - length) / 2
	rightPadding := width - length - leftPadding

	return fmt.Sprintf("[%s%s%s]",
		strings.Repeat(" ", leftPadding),
		logType,
		strings.Repeat(" ", rightPadding))
}

// Internal helper to construct the log strings
func writeLog(lType enums.LogType, color string, message string, location []string) {
	typeStr := getPaddedType(string(lType))
	timeStr := getFormattedDateTime()
	
	loc := ""
	if len(location) > 0 && location[0] != "" {
		loc = location[0]
	}

	locPart := ""
	locPartPlain := ""
	if loc != "" {
		locPart = fmt.Sprintf(" - %s[%s]%s", ColorGray, loc, ColorReset)
		locPartPlain = fmt.Sprintf(" - [%s]", loc)
	}

	logWithColor := fmt.Sprintf("%s%s%s %s%s%s%s - %s",
		ColorGray, timeStr, ColorReset,
		color, typeStr, ColorReset,
		locPart, message)

	logWithoutColor := fmt.Sprintf("%s %s%s - %s",
		timeStr, typeStr, locPartPlain, message)

	logFile(logWithColor, logWithoutColor)
}

// --- Exported Functions ---

// Error
func LogError(message string, location ...string) {
	writeLog(enums.ERROR, ColorRed, message, location)
}

func LogErrorF(format string, args ...interface{}) {
	writeLog(enums.ERROR, ColorRed, fmt.Sprintf(format, args...), nil)
}

// Warning
func LogWarning(message string, location ...string) {
	writeLog(enums.WARNING, ColorYellow, message, location)
}

func LogWarningF(format string, args ...interface{}) {
	writeLog(enums.WARNING, ColorYellow, fmt.Sprintf(format, args...), nil)
}

// Info
func LogInfo(message string, location ...string) {
	writeLog(enums.INFO, ColorGreen, message, location)
}

func LogInfoF(format string, args ...interface{}) {
	writeLog(enums.INFO, ColorGreen, fmt.Sprintf(format, args...), nil)
}

// Audit
func LogAudit(message string, location ...string) {
	writeLog(enums.AUDIT, ColorCyan, message, location)
}

func LogAuditF(format string, args ...interface{}) {
	writeLog(enums.AUDIT, ColorCyan, fmt.Sprintf(format, args...), nil)
}

// Event
func LogEvent(message string, location ...string) {
	writeLog(enums.EVENT, ColorMagenta, message, location)
}

func LogEventF(format string, args ...interface{}) {
	writeLog(enums.EVENT, ColorMagenta, fmt.Sprintf(format, args...), nil)
}

// Debug
func LogDebug(message string, location ...string) {
	if utils.DOTENV.Environment == "production" {
		return
	}
	writeLog(enums.DEBUG, ColorReset, message, location)
}

func LogDebugF(format string, args ...interface{}) {
	if utils.DOTENV.Environment == "production" {
		return
	}
	writeLog(enums.DEBUG, ColorReset, fmt.Sprintf(format, args...), nil)
}

// LogFatal 
func LogFatal(message string, location ...string) {
    LogError(message, location...);
    os.Exit(1); // Kill the app
}

func LogFatalF(format string, args ...interface{}) {
    LogErrorF(format, args...);
    os.Exit(1);
}