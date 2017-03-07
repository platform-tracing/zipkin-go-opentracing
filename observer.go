package zipkintracer

import (
	opentracing "github.com/opentracing/opentracing-go"
)

// observer is a dispatcher to other observers
type observer struct {
	observers []opentracing.Observer
}

// spanObserver is a dispatcher to other span observers
type spanObserver struct {
	observers []opentracing.SpanObserver
}

// noopSpanObserver is used when there are no observers registered on the
// Tracer or none of them returns span observers
var noopSpanObserver = spanObserver{}

func (o observer) OnStartSpan(sp opentracing.Span, operationName string, options opentracing.StartSpanOptions) opentracing.SpanObserver {
	var spanObservers []opentracing.SpanObserver
	for _, obs := range o.observers {
		spanObs := obs.OnStartSpan(sp, operationName, options)
		if spanObs != nil {
			if spanObservers == nil {
				spanObservers = make([]opentracing.SpanObserver, 0, len(o.observers))
			}
			spanObservers = append(spanObservers, spanObs)
		}
	}
	if len(spanObservers) == 0 {
		return noopSpanObserver
	}

	return spanObserver{observers: spanObservers}
}

func (o spanObserver) OnSetOperationName(operationName string) {
	for _, obs := range o.observers {
		obs.OnSetOperationName(operationName)
	}
}

func (o spanObserver) OnSetTag(key string, value interface{}) {
	for _, obs := range o.observers {
		obs.OnSetTag(key, value)
	}
}

func (o spanObserver) OnFinish(options opentracing.FinishOptions) {
	for _, obs := range o.observers {
		obs.OnFinish(options)
	}
}
