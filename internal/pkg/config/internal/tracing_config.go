package internal

import "time"

type TracingConfig struct {
	OTLPHttpEndpoint     string
	OTLPLogHttpEndpoint  string
	MetricReaderInterval time.Duration
}
