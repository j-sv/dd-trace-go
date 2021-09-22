// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package httputil // import "gopkg.in/DataDog/dd-trace-go.v1/contrib/internal/httputil"

//go:generate sh -c "go run make_responsewriter.go | gofmt > trace_gen.go"

import (
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	httpinstr "gopkg.in/DataDog/dd-trace-go.v1/internal/appsec/instrumentation/http"
	appsectypes "gopkg.in/DataDog/dd-trace-go.v1/internal/appsec/types"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/dyngo"
)

// TraceConfig defines the configuration for request tracing.
type TraceConfig struct {
	ResponseWriter http.ResponseWriter       // response writer
	Request        *http.Request             // request that is traced
	Service        string                    // service name
	Resource       string                    // resource name
	QueryParams    bool                      // specifies that request query parameters should be appended to http.url tag
	FinishOpts     []ddtrace.FinishOption    // span finish options to be applied
	SpanOpts       []ddtrace.StartSpanOption // additional span options to be applied
}

// TraceAndServe will apply tracing to the given http.Handler using the passed tracer under the given service and resource.
func TraceAndServe(h http.Handler, cfg *TraceConfig) {
	path := cfg.Request.URL.Path
	if cfg.QueryParams {
		path += "?" + cfg.Request.URL.RawQuery
	}
	opts := append([]ddtrace.StartSpanOption{
		tracer.SpanType(ext.SpanTypeWeb),
		tracer.ServiceName(cfg.Service),
		tracer.ResourceName(cfg.Resource),
		tracer.Tag(ext.HTTPMethod, cfg.Request.Method),
		tracer.Tag(ext.HTTPURL, path),
	}, cfg.SpanOpts...)
	if cfg.Request.URL.Host != "" {
		opts = append([]ddtrace.StartSpanOption{
			tracer.Tag("http.host", cfg.Request.URL.Host),
		}, opts...)
	}
	if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(cfg.Request.Header)); err == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}
	span, ctx := tracer.StartSpanFromContext(cfg.Request.Context(), "http.request", opts...)
	defer span.Finish(cfg.FinishOpts...)

	cfg.ResponseWriter = wrapResponseWriter(cfg.ResponseWriter, span)

	op := dyngo.StartOperation(
		httpinstr.HandlerOperationArgs{
			IsTLS:       cfg.Request.TLS != nil,
			Method:      httpinstr.Method(cfg.Request.Method),
			Host:        httpinstr.Host(cfg.Request.Host),
			RequestURI:  httpinstr.RequestURI(cfg.Request.RequestURI),
			RemoteAddr:  httpinstr.RemoteAddr(cfg.Request.RemoteAddr),
			Headers:     httpinstr.Header(cfg.Request.Header),
			UserAgent:   httpinstr.UserAgent(cfg.Request.UserAgent()),
			QueryValues: httpinstr.QueryValues(cfg.Request.URL.Query()),
		},
	)
	op.OnData(func(e *appsectypes.SecurityEvent) {
		// Keep this trace due to the security event
		span.SetTag(ext.SamplingPriority, ext.ManualKeep)
		// Add context to the event
		spanCtx := span.Context()
		// Add the APM context to the event
		e.AddContext(appsectypes.SpanContext{
			TraceID: spanCtx.TraceID(),
			SpanID:  spanCtx.SpanID(),
		})
	})
	defer func() {
		// TODO(julio): get the status code out of the wrapped response writer
		op.Finish(httpinstr.HandlerOperationRes{})
	}()

	h.ServeHTTP(cfg.ResponseWriter, cfg.Request.WithContext(ctx))
}

// responseWriter is a small wrapper around an http response writer that will
// intercept and store the status of a request.
type responseWriter struct {
	http.ResponseWriter
	span   ddtrace.Span
	status int
}

func newResponseWriter(w http.ResponseWriter, span ddtrace.Span) *responseWriter {
	return &responseWriter{w, span, 0}
}

// Write writes the data to the connection as part of an HTTP reply.
// We explicitely call WriteHeader with the 200 status code
// in order to get it reported into the span.
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
// It also sets the status code to the span.
func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.ResponseWriter.WriteHeader(status)
	w.status = status
	w.span.SetTag(ext.HTTPCode, strconv.Itoa(status))
	if status >= 500 && status < 600 {
		w.span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, http.StatusText(status)))
	}
}
