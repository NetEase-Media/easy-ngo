package xzap

import (
	"go.uber.org/zap"
)

func FieldTraceId(value string) zap.Field {
	return zap.String("trace_id", value)
}

func FieldSpanId(value string) zap.Field {
	return zap.String("span_id", value)
}
