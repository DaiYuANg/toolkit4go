package prometheus

import (
	"errors"
	"fmt"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/samber/oops"
)

func (a *Adapter) counter(metricName string, labels map[string]string) (*counterInstrument, error) {
	if a == nil {
		return nil, oops.In("observabilityx/prometheus").
			With("op", "create_counter", "metric", metricName).
			New("adapter is nil")
	}
	if a.register == nil {
		return nil, oops.In("observabilityx/prometheus").
			With("op", "create_counter", "metric", metricName).
			New("registerer is nil")
	}
	if a.counters == nil {
		return nil, oops.In("observabilityx/prometheus").
			With("op", "create_counter", "metric", metricName).
			New("counter registry is nil")
	}
	if existing, ok := a.counters.Get(metricName); ok {
		return existing, nil
	}

	labelNames := sortedLabelKeys(labels)
	vec := prom.NewCounterVec(prom.CounterOpts{
		Name: metricName,
		Help: "Counter metric for " + metricName,
	}, labelNames)

	if err := a.register.Register(vec); err != nil {
		alreadyRegisteredError, ok := errors.AsType[*prom.AlreadyRegisteredError](err)
		if !ok || alreadyRegisteredError == nil {
			return nil, oops.In("observabilityx/prometheus").
				With("op", "create_counter", "metric", metricName, "label_count", len(labelNames)).
				Wrapf(err, "register prometheus counter")
		}
		existingVec, ok := alreadyRegisteredError.ExistingCollector.(*prom.CounterVec)
		if !ok {
			return nil, oops.In("observabilityx/prometheus").
				With("op", "create_counter", "metric", metricName, "label_count", len(labelNames), "collector_type", fmt.Sprintf("%T", alreadyRegisteredError.ExistingCollector)).
				Errorf("prometheus counter has unexpected collector type")
		}
		vec = existingVec
	}

	inst := &counterInstrument{
		labels: labelNames,
		vec:    vec,
	}

	actual, _ := a.counters.GetOrStore(metricName, inst)
	return actual, nil
}

func (a *Adapter) histogram(metricName string, labels map[string]string) (*histInstrument, error) {
	if a == nil {
		return nil, oops.In("observabilityx/prometheus").
			With("op", "create_histogram", "metric", metricName).
			New("adapter is nil")
	}
	if a.register == nil {
		return nil, oops.In("observabilityx/prometheus").
			With("op", "create_histogram", "metric", metricName).
			New("registerer is nil")
	}
	if a.histograms == nil {
		return nil, oops.In("observabilityx/prometheus").
			With("op", "create_histogram", "metric", metricName).
			New("histogram registry is nil")
	}
	if existing, ok := a.histograms.Get(metricName); ok {
		return existing, nil
	}

	labelNames := sortedLabelKeys(labels)
	vec := prom.NewHistogramVec(prom.HistogramOpts{
		Name:    metricName,
		Help:    "Histogram metric for " + metricName,
		Buckets: a.buckets,
	}, labelNames)

	if err := a.register.Register(vec); err != nil {
		alreadyRegisteredError, ok := errors.AsType[*prom.AlreadyRegisteredError](err)
		if !ok || alreadyRegisteredError == nil {
			return nil, oops.In("observabilityx/prometheus").
				With("op", "create_histogram", "metric", metricName, "label_count", len(labelNames)).
				Wrapf(err, "register prometheus histogram")
		}
		existingVec, ok := alreadyRegisteredError.ExistingCollector.(*prom.HistogramVec)
		if !ok {
			return nil, oops.In("observabilityx/prometheus").
				With("op", "create_histogram", "metric", metricName, "label_count", len(labelNames), "collector_type", fmt.Sprintf("%T", alreadyRegisteredError.ExistingCollector)).
				Errorf("prometheus histogram has unexpected collector type")
		}
		vec = existingVec
	}

	inst := &histInstrument{
		labels: labelNames,
		vec:    vec,
	}

	actual, _ := a.histograms.GetOrStore(metricName, inst)
	return actual, nil
}
