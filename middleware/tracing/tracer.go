package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
	opt    *options
}

// NewTracer 创建一个 tracer
func NewTracer(kind trace.SpanKind, opts ...Option) *Tracer {
	op := options{
		propagator: propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}),
		tracerName: "gokit-otel-tracer",
	}
	for _, o := range opts {
		o(&op)
	}
	if op.tracerProvider == nil {
		op.tracerProvider = otel.GetTracerProvider()
	}

	switch kind {
	case trace.SpanKindClient:
		return &Tracer{tracer: op.tracerProvider.Tracer(op.tracerName), kind: kind, opt: &op}
	case trace.SpanKindServer:
		return &Tracer{tracer: op.tracerProvider.Tracer(op.tracerName), kind: kind, opt: &op}
	default:
		panic(fmt.Sprintf("unsupported span kind: %v", kind))
	}
}

// Start 开启一个 span
func (t *Tracer) Start(ctx context.Context, spanName string, carrier propagation.TextMapCarrier) (context.Context, trace.Span) {
	if t.kind == trace.SpanKindServer {
		// 自 server header 中提取 trace 元数据并写入 ctx
		ctx = t.opt.propagator.Extract(ctx, carrier)
	}
	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(t.kind))
	if t.kind == trace.SpanKindClient {
		// 将 ctx 中的 trace 元数据写入 client header
		t.opt.propagator.Inject(ctx, carrier)
	}
	return ctx, span
}

// End 结束一个 span
func (t *Tracer) End(_ context.Context, span trace.Span, reply any, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "OK")
	}

	if p, ok := reply.(proto.Message); ok {
		if t.kind == trace.SpanKindServer {
			span.SetAttributes(attribute.Key("send_msg.size").Int(proto.Size(p)))
		} else {
			span.SetAttributes(attribute.Key("recv_msg.size").Int(proto.Size(p)))
		}
	}
	span.End()
}
