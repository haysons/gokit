package metrics

import (
	"context"
	"time"

	"github.com/haysons/gokit/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
)

const (
	metricLabelKind      = "kind"
	metricLabelOperation = "operation"
	metricLabelCode      = "code"
	metricLabelReason    = "reason"
)

const (
	DefaultServerSecondsHistogramName = "server_requests_seconds_bucket"
	DefaultServerRequestsCounterName  = "server_requests_code_total"
	DefaultClientSecondsHistogramName = "client_requests_seconds_bucket"
	DefaultClientRequestsCounterName  = "client_requests_code_total"
)

type Option func(*options)

// WithRequests 统计不同类型和状态的累计请求数
func WithRequests(c metric.Int64Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds 统计请求时间的分布情况
func WithSeconds(histogram metric.Float64Histogram) Option {
	return func(o *options) {
		o.seconds = histogram
	}
}

// DefaultRequestsCounter 请求计数器，构造完成后可通过 WithRequests 统计累计请求数
func DefaultRequestsCounter(meter metric.Meter, name string) (metric.Int64Counter, error) {
	return meter.Int64Counter(name, metric.WithUnit("{call}"))
}

// DefaultSecondsHistogram 请求时间直方图，构造完成后可通过 WithSeconds 统计请求耗时的分布情况
func DefaultSecondsHistogram(meter metric.Meter, name string) (metric.Float64Histogram, error) {
	return meter.Float64Histogram(
		name,
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1),
	)
}

func DefaultSecondsHistogramView(histogramName string) metricsdk.View {
	return func(instrument metricsdk.Instrument) (metricsdk.Stream, bool) {
		if instrument.Name == histogramName {
			return metricsdk.Stream{
				Name:        instrument.Name,
				Description: instrument.Description,
				Unit:        instrument.Unit,
				Aggregation: metricsdk.AggregationExplicitBucketHistogram{
					Boundaries: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
					NoMinMax:   true,
				},
				AttributeFilter: func(attribute.KeyValue) bool {
					return true
				},
			}, true
		}
		return metricsdk.Stream{}, false
	}
}

type options struct {
	// 请求计数器
	requests metric.Int64Counter
	// 请求时间分布直方图
	seconds metric.Float64Histogram
}

// Server is middleware server-side metrics.
func Server(opts ...Option) middleware.Middleware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			// if requests and seconds are nil, return directly
			if op.requests == nil && op.seconds == nil {
				return handler(ctx, req)
			}

			// 自 ctx 中提取请求的元数据
			var (
				code      int
				reason    string
				kind      string
				operation string
			)

			// default code
			startTime := time.Now()
			reply, err := handler(ctx, req)
			if op.requests != nil {
				op.requests.Add(
					ctx, 1,
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
						attribute.Int(metricLabelCode, code),
						attribute.String(metricLabelReason, reason),
					),
				)
			}
			if op.seconds != nil {
				op.seconds.Record(
					ctx, time.Since(startTime).Seconds(),
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
					),
				)
			}
			return reply, err
		}
	}
}

// Client is middleware client-side metrics.
func Client(opts ...Option) middleware.Middleware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			// 自 ctx 中提取请求的元数据
			var (
				code      int
				reason    string
				kind      string
				operation string
			)

			startTime := time.Now()
			reply, err := handler(ctx, req)
			if op.requests != nil {
				op.requests.Add(
					ctx, 1,
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
						attribute.Int(metricLabelCode, code),
						attribute.String(metricLabelReason, reason),
					),
				)
			}
			if op.seconds != nil {
				op.seconds.Record(
					ctx, time.Since(startTime).Seconds(),
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
					),
				)
			}
			return reply, err
		}
	}
}
