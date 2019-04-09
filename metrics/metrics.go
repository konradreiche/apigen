package metrics

import "time"

// Counter is the most basic and default metric type. It is treated as a count
// of a type of event per second and are typically averaged over on minute.
// That is, when looking at a graph, you are usually seeing the average number
// of events per second during a one-minute period.
type Counter interface {
	With(Tag) Counter
	Add(delta int64)
}

// Gauge is a constant metric type. Gauges are not subject to averaging and do
// not change unless they are changed. That is, once a gauge value is set, it
// will be a flat line on the graph until it is changed again.
//
// Examples: system load, active WebSocket connections
type Gauge interface {
	With(Tag) Gauge
	Set(value float64)
}

// TimeHistogram calculate the statistical distribution of any kind of value
// over time. An example of a histogram is the request latency.
type TimeHistogram interface {
	With(Tag) TimeHistogram
	Observe(d time.Duration)
}

// Tag is a key-value pair which can be associated with a metric in order to
// add additional information.
type Tag struct {
	Key   string
	Value string
}
