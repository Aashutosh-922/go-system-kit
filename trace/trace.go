package trace

import (
	"context"
)

type key string

const TraceKey key = "trace_id"

func WithTrace(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, TraceKey, id)
}

func GetTrace(ctx context.Context) string {
	v := ctx.Value(TraceKey)

	if v == nil {
		return ""
	}

	return v.(string)
}
