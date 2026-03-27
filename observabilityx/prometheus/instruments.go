package prometheus

import (
	"errors"
	"fmt"

	prom "github.com/prometheus/client_golang/prometheus"
)

func (a *Adapter) counter(metricName string, labels map[string]string) (*counterInstrument, error) {
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
			return nil, fmt.Errorf("register prometheus counter %q: %w", metricName, err)
		}
		existingVec, ok := alreadyRegisteredError.ExistingCollector.(*prom.CounterVec)
		if !ok {
			return nil, fmt.Errorf("prometheus counter %q has unexpected collector type %T", metricName, alreadyRegisteredError.ExistingCollector)
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
			return nil, fmt.Errorf("register prometheus histogram %q: %w", metricName, err)
		}
		existingVec, ok := alreadyRegisteredError.ExistingCollector.(*prom.HistogramVec)
		if !ok {
			return nil, fmt.Errorf("prometheus histogram %q has unexpected collector type %T", metricName, alreadyRegisteredError.ExistingCollector)
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
