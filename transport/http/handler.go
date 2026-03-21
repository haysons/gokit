package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/haysons/gokit/middleware"
	"github.com/haysons/gokit/transport"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// HandlerFunc represents a gRPC-style handler function
// It has the same signature as gRPC unary handlers
type HandlerFunc func(ctx context.Context, req any) (any, error)

// RegisterHandler registers a gRPC-style handler to the HTTP server
//
// Example:
//
//	type GreeterServer struct {}
//	func (s *GreeterServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
//	    return &HelloReply{Message: "Hello " + req.Name}, nil
//	}
//
//	srv := http.NewServer()
//	srv.Use(middleware.Logging)
//	srv.RegisterHandler("POST", "/api/hello", GreeterServer.SayHello, &HelloRequest{}, &HelloReply{})
func (s *Server) RegisterHandler(method, pattern string, handlerFunc interface{}, reqProto, respProto proto.Message) {
	s.registerHandlerWithMiddleware(method, pattern, handlerFunc, reqProto, respProto, s.middlewares)
}

// RegisterHandlerWithMux registers a handler with custom mux options
func (s *Server) RegisterHandlerWithMux(method, pattern string, handlerFunc interface{}, reqProto, respProto proto.Message, muxOpts ...gruntime.ServeMuxOption) {
	s.registerHandlerWithMiddleware(method, pattern, handlerFunc, reqProto, respProto, s.middlewares)
}

func (s *Server) registerHandlerWithMiddleware(method, pattern string, handlerFunc interface{}, reqProto, respProto proto.Message, middlewares []middleware.Middleware) {
	// Convert handlerFunc to HandlerFunc
	handler := s.wrapHandler(handlerFunc, reqProto, respProto, middlewares)

	// Get HTTP method
	httpMethod := method
	if httpMethod == "" {
		httpMethod = "POST"
	}

	// Register to grpc-gateway mux using HandlePath (takes string pattern)
	_ = s.mux.HandlePath(httpMethod, pattern, handler)
}

func (s *Server) wrapHandler(handlerFunc interface{}, reqProto, respProto proto.Message, middlewares []middleware.Middleware) func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		ctx := r.Context()

		// Extract metadata from HTTP headers
		md := metadata.MD{}
		for k, v := range r.Header {
			md[k] = v
		}
		ctx = metadata.NewIncomingContext(ctx, md)

		// Create transport layer info and inject into context
		headerCarrier := make(http.Header)
		replyHeader := make(http.Header)
		tr := &Transport{
			reqHeader:   headerCarrier,
			replyHeader: replyHeader,
		}
		ctx = transport.InjectServerContext(ctx, tr)

		// Get request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.writeError(w, r, err)
			return
		}

		// Decode request (JSON -> Protobuf)
		req := reflect.New(reflect.TypeOf(reqProto).Elem()).Interface().(proto.Message)
		if len(body) > 0 {
			if err := protojson.Unmarshal(body, req); err != nil {
				s.writeError(w, r, err)
				return
			}
		}

		// Build handler chain
		h := func(ctx context.Context, req any) (any, error) {
			// Get the actual handler function using reflection
			handlerValue := reflect.ValueOf(handlerFunc)
			handlerType := handlerValue.Type()

			// Call handler with correct argument types
			if handlerType.Kind() == reflect.Func {
				results := handlerValue.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					reflect.ValueOf(req),
				})
				if len(results) == 2 {
					if !results[1].IsNil() {
						return nil, results[1].Interface().(error)
					}
					return results[0].Interface(), nil
				}
			}
			return nil, nil
		}

		// Apply middlewares (from outermost to innermost)
		for i := len(middlewares) - 1; i >= 0; i-- {
			nextHandler := h
			h = middlewares[i](nextHandler)
		}

		// Execute handler chain
		resp, err := h(ctx, req)
		if err != nil {
			s.writeError(w, r, err)
			return
		}

		// Encode response (Protobuf -> JSON)
		respBytes, err := protojson.Marshal(resp.(proto.Message))
		if err != nil {
			s.writeError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respBytes)
	}
}

func (s *Server) writeError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	errorResp := map[string]string{"error": err.Error()}
	jsonBytes, _ := json.Marshal(errorResp)
	w.Write(jsonBytes)
}
