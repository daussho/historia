package trace

import (
	"time"

	"github.com/daussho/historia/domain/common"
)

type Trace struct {
	TraceID      string                        `json:"trace_id" gorm:"column:trace_id;primaryKey"`
	ParentSpanID string                        `json:"parent_span_id" gorm:"column:parent_span_id"`
	SpanID       string                        `json:"span_id" gorm:"column:span_id"`
	SegmentName  string                        `json:"segment_name" gorm:"column:segment_name"`
	Tags         common.SQLMap[string, string] `json:"tags" gorm:"column:tags"`
	StartAt      time.Time                     `json:"start_at" gorm:"column:start_at"`
	EndAt        time.Time                     `json:"end_at" gorm:"column:end_at"`
}
