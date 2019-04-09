package dogstatsd

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/konradreiche/apigen/metrics"
	"github.com/sirupsen/logrus"
)

type counter struct {
	name       string
	sampleRate float64
	tags       []metrics.Tag
	client     *statsd.Client
	log        *logrus.Logger
}

// NewCounter returns a new instance of counter.
func NewCounter(name string, sampleRate float64, c *statsd.Client, l *logrus.Logger) metrics.Counter {
	return &counter{name: name, sampleRate: sampleRate, client: c, log: l}
}

// With adds a tag to the observation.
func (c *counter) With(tag metrics.Tag) metrics.Counter {
	return &counter{
		name:       c.name,
		sampleRate: c.sampleRate,
		tags:       append(c.tags, tag),
		client:     c.client,
		log:        c.log,
	}
}

// Add increments a value by the given delta and writes it to the StatsD client.
func (c *counter) Add(delta int64) {
	var tags []string
	for _, tag := range c.tags {
		tags = append(tags, fmt.Sprintf("%s:%s", tag.Key, tag.Value))
	}
	err := c.client.Count(c.name, delta, tags, c.sampleRate)
	if err != nil {
		c.log.Error(err)
	}
}

type gauge struct {
	name       string
	sampleRate float64
	tags       []metrics.Tag
	client     *statsd.Client
	log        *logrus.Logger
}

// NewGauge returns a new instance of gauge.
func NewGauge(name string, sampleRate float64, c *statsd.Client, l *logrus.Logger) metrics.Gauge {
	return &gauge{name: name, sampleRate: sampleRate, client: c, log: l}
}

// With adds a tag to the observation.
func (g *gauge) With(tag metrics.Tag) metrics.Gauge {
	return &gauge{
		name:       g.name,
		sampleRate: g.sampleRate,
		tags:       append(g.tags, tag),
		client:     g.client,
		log:        g.log,
	}
}

// Set updates the gauge to the given value.
func (g *gauge) Set(value float64) {
	var tags []string
	for _, tag := range g.tags {
		tags = append(tags, fmt.Sprintf("%s:%s", tag.Key, tag.Value))
	}
	err := g.client.Gauge(g.name, value, tags, g.sampleRate)
	if err != nil {
		g.log.Error(err)
	}
}

type timeHistogram struct {
	name       string
	sampleRate float64
	tags       []metrics.Tag
	client     *statsd.Client
	log        *logrus.Logger
}

func NewTimeHistogram(name string, sampleRate float64, c *statsd.Client, l *logrus.Logger) metrics.TimeHistogram {
	return &timeHistogram{name: name, sampleRate: sampleRate, client: c, log: l}
}

// With adds a tag to the observation.
func (th *timeHistogram) With(tag metrics.Tag) metrics.TimeHistogram {
	return &timeHistogram{
		name:       th.name,
		sampleRate: th.sampleRate,
		tags:       append(th.tags, tag),
		client:     th.client,
		log:        th.log,
	}
}

func (th *timeHistogram) Observe(value time.Duration) {
	var tags []string
	for _, tag := range th.tags {
		tags = append(tags, fmt.Sprintf("%s:%s", tag.Key, tag.Value))
	}
	err := th.client.Timing(th.name, value, tags, th.sampleRate)
	if err != nil {
		th.log.Error(err)
	}
}
