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

type ApiHealthStatus string

const (
	Healthy       ApiHealthStatus = "healthy"
	Degraded      ApiHealthStatus = "degraded"
	UnhealthyDown ApiHealthStatus = "critical"
	Maintenance   ApiHealthStatus = "maintenance_unavailable"
	ShuttingDown  ApiHealthStatus = "unavailable"
	Initializing  ApiHealthStatus = "initializing"
)

// just to map the enum
