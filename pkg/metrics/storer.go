package metrics

type Storer interface {
	Receive(m Metric)
	GetLatest() []Metric
	GetHistory(source string, limit int) []Metric
}
