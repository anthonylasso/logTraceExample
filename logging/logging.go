package logging

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type TraceKey int

const KEY TraceKey = 1

// can set header explicitly or auto generate on a per-request basis
const TRACE_HEADER string = "X-TraceId"

func contextWithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, KEY, traceId)
}

func TraceIdFromContext(ctx context.Context) (string, bool) {
	traceId, ok := ctx.Value(KEY).(string)
	return traceId, ok
}

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respWrite http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if traceId := request.Header.Get(TRACE_HEADER); traceId != "" {
			ctx = contextWithTraceId(ctx, traceId)
		} else {
			ctx = contextWithTraceId(ctx, uuid.New().String())
		}
		request = request.WithContext(ctx)
		handler.ServeHTTP(respWrite, request)
	})
}

type Logger struct {
	//note: can be any logger we choose
}

func (Logger) Log(ctx context.Context, message string) {
	if traceId, ok := TraceIdFromContext(ctx); ok {
		message = fmt.Sprintf("TraceId: %s - %s", traceId, message)
	}
	fmt.Println(message)
}

func Request(request *http.Request) *http.Request {
	ctx := request.Context()
	if traceId, ok := TraceIdFromContext(ctx); ok {
		request.Header.Add(TRACE_HEADER, traceId)
	}
	return request
}
