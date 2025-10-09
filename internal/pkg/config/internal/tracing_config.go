package internal

import "time"

type TracingConfig struct {
	OTLPHttpEndpoint     string
	MetricReaderInterval time.Duration
}
