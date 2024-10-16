package trace

import (
	"context"
	"time"

	"github.com/daussho/historia/utils/clock"
	"github.com/google/uuid"
)

type Span struct {
	TraceID      string            `json:"trace_id"`
	ParentSpanID string            `json:"parent_span_id"`
	SpanID       string            `json:"span_id"`
	SegmentName  string            `json:"segment_name"`
	Tags         map[string]string `json:"tags"`
	StartAt      time.Time         `json:"start_at"`
	EndAt        time.Time         `json:"end_at"`
}

type contextKey string

const (
	TraceIDKey      contextKey = "internal-trace-id"
	ParentSpanIDKey contextKey = "internal-parent-span-id"
	tracesKey       contextKey = "internal-traces"
)

func StartSpanWithCtx(ctx context.Context, segmentName string, tags map[string]string) (*Span, context.Context) {
	var spans []*Span
	if tempSpans := ctx.Value(tracesKey); tempSpans != nil {
		spans = tempSpans.([]*Span)
	}

	var traceID string
	if tempTraceID := ctx.Value(TraceIDKey); tempTraceID != nil {
		traceID = tempTraceID.(string)
	} else {
		traceID = uuid.NewString()
	}

	var parentSpanID string
	if tempParentSpanID := ctx.Value(ParentSpanIDKey); tempParentSpanID != nil {
		parentSpanID = tempParentSpanID.(string)
	}

	spanID := uuid.NewString()
	span := &Span{
		TraceID:      traceID,
		ParentSpanID: parentSpanID,
		SpanID:       spanID,
		SegmentName:  segmentName,
		Tags:         tags,
		StartAt:      clock.Now(),
		EndAt:        clock.Now(),
	}

	spans = append(spans, span)

	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	ctx = context.WithValue(ctx, ParentSpanIDKey, spanID)
	ctx = context.WithValue(ctx, tracesKey, spans)

	return span, ctx
}

func (s *Span) Finish() {
	s.EndAt = clock.Now()
}

func (s *Span) FinishAndSubmit(ctx context.Context) {
	s.Finish()

	traces := ctx.Value(tracesKey)
	if traces == nil {
		return
	}

	// log.Println(utils.JsonStringify(traces))
}
