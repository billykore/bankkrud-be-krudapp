package telemetry

import (
	"context"
	"errors"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// Telemetry wraps the OpenTelemetry.
type Telemetry struct {
	cfg        *config.Configs
	tracer     trace.Tracer
	shutdownFn func(context.Context) error
}

// New creates a new Telemetry instance.
func New(cfg *config.Configs) *Telemetry {
	shutdownFn, err := setupOTelSDK(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	return &Telemetry{
		cfg:        cfg,
		tracer:     otel.Tracer(cfg.App.Name),
		shutdownFn: shutdownFn,
	}
}

// Tracer returns the tracer instance.
func (t *Telemetry) Tracer() trace.Tracer {
	return t.tracer
}

// Shutdown shuts down the tracer.
func (t *Telemetry) Shutdown(ctx context.Context) error {
	return t.shutdownFn(ctx)
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context, cfg *config.Configs) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up exporter.
	exporter, err := newExporter(ctx, cfg.Tracing.OTLPHttpEndpoint)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(ctx, cfg.App.Name, cfg.App.Version, exporter)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(cfg.App.Name)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx, cfg.Tracing.OTLPLogHTTPEndpoint)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func newTracerProvider(ctx context.Context, serviceName string, serviceVersion string, exp sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create a new trace provider with the exporter
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	), nil
}

func newMeterProvider(appName string) (*metric.MeterProvider, error) {
	metricExporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricExporter),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
		)),
	)
	meterProvider.Meter(appName)
	return meterProvider, nil
}

func newLoggerProvider(ctx context.Context, endpoint string) (*log.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpoint(endpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
