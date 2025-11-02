package metrics

type MetricSinkComp struct {
	stores []Storer
}

func NewMetricSink(stores ...Storer) *MetricSinkComp {
	return &MetricSinkComp{stores: stores}
}

func (ms *MetricSinkComp) Receive(m Metric) {
	for _, store := range ms.stores {
		store.Receive(m)
	}
}
