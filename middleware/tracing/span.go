package tracing

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/haysons/gokit/transport"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

// setClientSpan 向 client span 中添加属性信息
func setClientSpan(ctx context.Context, span trace.Span, m any) {
	var (
		attrs     []attribute.KeyValue
		remote    string
		operation string
		rpcKind   string
	)
	tr, ok := transport.FromClientContext(ctx)
	if ok {
		operation = tr.Operation()
		rpcKind = tr.Kind().String()
		switch tr.Kind() {
		case transport.KindHTTP:
			if req := tr.Request(); req != nil {
				if httpReq, ok := req.(*http.Request); ok {
					attrs = append(attrs, semconv.HTTPMethodKey.String(httpReq.Method))
					attrs = append(attrs, semconv.HTTPTargetKey.String(httpReq.URL.Path))
					attrs = append(attrs, semconv.HTTPRouteKey.String(tr.PathTemplate()))
					remote = httpReq.Host
				}
			}
		case transport.KindGRPC:
			remote, _ = parseTarget(tr.Endpoint())
		}
	}
	attrs = append(attrs, semconv.RPCSystemKey.String(rpcKind))
	_, mAttrs := parseFullMethod(operation)
	attrs = append(attrs, mAttrs...)
	if remote != "" {
		attrs = append(attrs, peerAttr(remote)...)
	}
	if p, ok := m.(proto.Message); ok {
		attrs = append(attrs, attribute.Key("send_msg.size").Int(proto.Size(p)))
	}

	span.SetAttributes(attrs...)
}

// setServerSpan 向 server span 添加属性信息
func setServerSpan(ctx context.Context, span trace.Span, m any) {
	var (
		attrs     []attribute.KeyValue
		remote    string
		operation string
		rpcKind   string
	)
	tr, ok := transport.FromServerContext(ctx)
	if ok {
		operation = tr.Operation()
		rpcKind = tr.Kind().String()
		switch tr.Kind() {
		case transport.KindHTTP:
			if req := tr.Request(); req != nil {
				if httpReq, ok := req.(*http.Request); ok {
					attrs = append(attrs, semconv.HTTPMethodKey.String(httpReq.Method))
					attrs = append(attrs, semconv.HTTPTargetKey.String(httpReq.URL.Path))
					attrs = append(attrs, semconv.HTTPRouteKey.String(tr.PathTemplate()))
					remote = httpReq.RemoteAddr
				}
			}
		case transport.KindGRPC:
			if p, ok := peer.FromContext(ctx); ok {
				remote = p.Addr.String()
			}
		}
	}
	attrs = append(attrs, semconv.RPCSystemKey.String(rpcKind))
	_, mAttrs := parseFullMethod(operation)
	attrs = append(attrs, mAttrs...)
	attrs = append(attrs, peerAttr(remote)...)
	if p, ok := m.(proto.Message); ok {
		attrs = append(attrs, attribute.Key("recv_msg.size").Int(proto.Size(p)))
	}

	span.SetAttributes(attrs...)
}

// parseFullMethod 获取 rpc 的 service 以及 method
func parseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 { //nolint:mnd
		// Invalid format, does not follow `/package.service/method`.
		return name, []attribute.KeyValue{attribute.Key("rpc.operation").String(fullMethod)}
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCServiceKey.String(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethodKey.String(method))
	}
	return name, attrs
}

// peerAttr 获取访问的 ip 及端口
func peerAttr(addr string) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return []attribute.KeyValue(nil)
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return []attribute.KeyValue{
		semconv.NetPeerIPKey.String(host),
		semconv.NetPeerPortKey.String(port),
	}
}

func parseTarget(endpoint string) (address string, err error) {
	var u *url.URL
	u, err = url.Parse(endpoint)
	if err != nil {
		if u, err = url.Parse("http://" + endpoint); err != nil {
			return "", err
		}
		return u.Host, nil
	}
	if len(u.Path) > 1 {
		return u.Path[1:], nil
	}
	return endpoint, nil
}
