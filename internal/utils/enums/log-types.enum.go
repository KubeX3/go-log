package enums

type LogType string

const (
	AUDIT   LogType = "AUDIT"
	DEBUG   LogType = "DEBUG"
	ERROR   LogType = "ERROR"
	EVENT   LogType = "EVENT"
	INFO    LogType = "INFO"
	WARNING LogType = "WARNING"
)
