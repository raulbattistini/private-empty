package enum

type ApiHealthStatus string

const (
	UnknownOrUnitialized ApiHealthStatus = "unknown"
	Health               ApiHealthStatus = "health"
	Pending              ApiHealthStatus = "op_in_progress"
	Failed               ApiHealthStatus = "failed_ops"
)
