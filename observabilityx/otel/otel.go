package otel

import (
	"context"
	"log/slog"
	"strings"

	collectionmapping "github.com/DaiYuANg/arcgo/collectionx/mapping"
	"github.com/DaiYuANg/arcgo/observabilityx"
	"github.com/DaiYuANg/arcgo/pkg/option"
	"github.com/samber/oops"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultTracerName = "github.com/DaiYuANg/arcgo"
	defaultMeterName  = "github.com/DaiYuANg/arcgo"
)

// Option configures OTel observability integration.
type Option func(*config)

type config struct {
	logger *slog.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

// WithLogger sets logger used by this adapter.
func WithLogger(logger *slog.Logger) Option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}

// WithTracer sets tracer used by this adapter.
func WithTracer(tracer trace.Tracer) Option {
	return func(cfg *config) {
		cfg.tracer = tracer
	}
}

// WithMeter sets meter used by this adapter.
func WithMeter(meter metric.Meter) Option {
	return func(cfg *config) {
		cfg.meter = meter
	}
}

// New creates an OTel-backed observability adapter.
func New(opts ...Option) observabilityx.Observability {
	cfg := config{
		logger: slog.Default(),
		tracer: otel.Tracer(defaultTracerName),
		meter:  otel.Meter(defaultMeterName),
	}
	option.Apply(&cfg, opts...)

	return &adapter{
		logger:     observabilityx.NormalizeLogger(cfg.logger),
		tracer:     cfg.tracer,
		meter:      cfg.meter,
		counters:   collectionmapping.NewConcurrentMap[string, metric.Int64Counter](),
		histograms: collectionmapping.NewConcurrentMap[string, metric.Float64Histogram](),
	}
}

type adapter struct {
	logger *slog.Logger
	tracer trace.Tracer
	meter  metric.Meter

	counters   *collectionmapping.ConcurrentMap[string, metric.Int64Counter]
	histograms *collectionmapping.ConcurrentMap[string, metric.Float64Histogram]
}

func (a *adapter) Logger() *slog.Logger {
	return observabilityx.NormalizeLogger(a.logger)
}

func (a *adapter) StartSpan(
	ctx context.Context,
	name string,
	attrs ...observabilityx.Attribute,
) (context.Context, observabilityx.Span) {
	return startTraceSpan(normalizeContext(ctx), a.tracer, normalizeSpanName(name), attrs)
}

func (a *adapter) AddCounter(
	ctx context.Context,
	name string,
	value int64,
	attrs ...observabilityx.Attribute,
) {
	if value == 0 {
		return
	}

	counter, err := a.counter(name)
	if err != nil {
		a.Logger().Warn("create metric counter failed", "name", name, "error", err.Error())
		return
	}
	counter.Add(normalizeContext(ctx), value, metric.WithAttributes(toOTelAttributes(attrs)...))
}

func (a *adapter) RecordHistogram(
	ctx context.Context,
	name string,
	value float64,
	attrs ...observabilityx.Attribute,
) {
	histogram, err := a.histogram(name)
	if err != nil {
		a.Logger().Warn("create metric histogram failed", "name", name, "error", err.Error())
		return
	}
	histogram.Record(normalizeContext(ctx), value, metric.WithAttributes(toOTelAttributes(attrs)...))
}

func (a *adapter) counter(name string) (metric.Int64Counter, error) {
	if a == nil {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_counter").
			New("adapter is nil")
	}
	clean := strings.TrimSpace(name)
	if clean == "" {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_counter").
			New("metric counter name is empty")
	}
	if a.meter == nil {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_counter", "metric", clean).
			New("meter is nil")
	}

	if existing, ok := a.counters.Get(clean); ok {
		return existing, nil
	}

	created, err := a.meter.Int64Counter(clean)
	if err != nil {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_counter", "metric", clean).
			Wrapf(err, "create OTel counter")
	}

	actual, _ := a.counters.GetOrStore(clean, created)
	return actual, nil
}

func (a *adapter) histogram(name string) (metric.Float64Histogram, error) {
	if a == nil {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_histogram").
			New("adapter is nil")
	}
	clean := strings.TrimSpace(name)
	if clean == "" {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_histogram").
			New("metric histogram name is empty")
	}
	if a.meter == nil {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_histogram", "metric", clean).
			New("meter is nil")
	}

	if existing, ok := a.histograms.Get(clean); ok {
		return existing, nil
	}

	created, err := a.meter.Float64Histogram(clean)
	if err != nil {
		return nil, oops.In("observabilityx/otel").
			With("op", "create_histogram", "metric", clean).
			Wrapf(err, "create OTel histogram")
	}

	actual, _ := a.histograms.GetOrStore(clean, created)
	return actual, nil
}

type otelSpan struct {
	span trace.Span
}

func (s otelSpan) End() {
	if s.span != nil {
		s.span.End()
	}
}

func (s otelSpan) RecordError(err error) {
	if s.span != nil && err != nil {
		s.span.RecordError(err)
	}
}

func (s otelSpan) SetAttributes(attrs ...observabilityx.Attribute) {
	if s.span == nil || len(attrs) == 0 {
		return
	}
	s.span.SetAttributes(toOTelAttributes(attrs)...)
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}

	return context.Background()
}

func normalizeSpanName(name string) string {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return "operation"
	}

	return cleanName
}

//nolint:spancheck // span ownership is transferred to the returned observabilityx.Span wrapper.
func startTraceSpan(ctx context.Context, tracer trace.Tracer, name string, attrs []observabilityx.Attribute) (context.Context, observabilityx.Span) {
	nextCtx, span := tracer.Start(ctx, name, trace.WithAttributes(toOTelAttributes(attrs)...))
	return nextCtx, otelSpan{span: span}
}
