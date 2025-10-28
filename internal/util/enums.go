package util

type EventSeverity string

const (
	Info  EventSeverity = "INFO"
	Error EventSeverity = "ERROR"
	Debug EventSeverity = "DEBUG"
	Warn  EventSeverity = "WARNING"
)

type CollectorsTypes string

const (
	CPU    CollectorsTypes = "cpu"
	Memory CollectorsTypes = "memory"
	Net    CollectorsTypes = "net"
)

// just to map the enum
