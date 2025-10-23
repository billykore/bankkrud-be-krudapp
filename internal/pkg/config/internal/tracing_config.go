package internal

import "time"

type TracingConfig struct {
	OTLPHttpEndpoint     string
	OTLPLogHTTPEndpoint  string
	MetricReaderInterval time.Duration
}
