package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Metrics struct {
	histogramVec *prometheus.HistogramVec
}

type MetricsTimer struct {
	histogramVec *prometheus.HistogramVec
	start        time.Time
	last         time.Time
}

func NewMetrics(namespace string, name string, labelName string, help string) Metrics {
	histogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: name,
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{labelName})
	prometheus.Register(histogramVec)
	return Metrics{
		histogramVec,
	}
}

func (metrics Metrics) NewTimer() MetricsTimer {
	now := time.Now()
	return MetricsTimer{
		histogramVec: metrics.histogramVec,
		start:        now,
		last:         now,
	}
}

func (t *MetricsTimer) ObserveTotal() {
	(*t.histogramVec).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
