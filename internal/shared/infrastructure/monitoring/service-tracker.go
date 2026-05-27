package monitoring

import "time"

type ServiceTracker struct {
	model     string
	operation string
	attrs     map[string]any
	startedAt time.Time
}

func StartService(model, operation string, attrs map[string]any) *ServiceTracker {
	tracker := &ServiceTracker{
		model:     model,
		operation: operation,
		attrs:     attrs,
		startedAt: time.Now(),
	}
	Service(model, operation, "started", attrs, nil, 0)
	return tracker
}

func (t *ServiceTracker) Success(extra map[string]any) {
	attrs := mergeAttrs(t.attrs, extra)
	Service(t.model, t.operation, "success", attrs, nil, time.Since(t.startedAt))
}

func (t *ServiceTracker) Fail(err error, extra map[string]any) {
	attrs := mergeAttrs(t.attrs, extra)
	Service(t.model, t.operation, "failed", attrs, err, time.Since(t.startedAt))
}

func mergeAttrs(base, extra map[string]any) map[string]any {
	if len(base) == 0 && len(extra) == 0 {
		return nil
	}
	merged := make(map[string]any, len(base)+len(extra))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range extra {
		merged[key] = value
	}
	return merged
}
