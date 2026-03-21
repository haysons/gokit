package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gokithttp "github.com/haysons/gokit/transport/http"
	"github.com/haysons/gokit/middleware"
	example "github.com/haysons/gokit/transport/testdata/helloworld"
)

// GreeterServer 是我们的业务实现
type GreeterServer struct {
	example.UnimplementedGreeterServer
}

// SayHello 实现 gRPC 服务接口
func (s *GreeterServer) SayHello(ctx context.Context, req *example.HelloRequest) (*example.HelloReply, error) {
	return &example.HelloReply{
		Message: "Hello " + req.Name + "!",
	}, nil
}

// 自定义中间件示例
func customMiddleware(next middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		log.Println("[Custom Middleware] Before handler")
		resp, err := next(ctx, req)
		log.Println("[Custom Middleware] After handler")
		return resp, err
	}
}

func main() {
	// 1. 创建 HTTP Server
	srv := gokithttp.NewServer(
		gokithttp.WithAddr(":8080"),
	)

	// 2. 注册中间件
	// 中间件会自动应用到所有 protoc 生成的 gateway 路由
	srv.Use(customMiddleware)        // 自定义中间件

	// 3. 注册 protoc 生成的 gateway 路由
	// 这里使用的是 protoc-gen-grpc-gateway 生成的 RegisterXxxHandlerServer 函数
	greeter := &GreeterServer{}
	// 注意：使用本地生成的代码需要手动导入
	// 在实际项目中，这行代码由 protoc 自动生成
	// RegisterGreeterHandlerServer(ctx, srv.GetMux(), greeter)

	// 为了演示，我们用手动注册的方式
	srv.RegisterHandler("POST", "/helloworld.Greeter/SayHello", greeter.SayHello, &example.HelloRequest{}, &example.HelloReply{})

	// 4. 启动服务器
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := srv.Start(context.Background()); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 5. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 6. 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
