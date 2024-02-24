package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "mygrpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

// Hello 関数は、リクエストからnameフィールドを取り出して、Hello、[名前]というレスポンスを返す
func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello %s!", req.GetName()),
	}, nil
}

// NewMyServer はサービス構造体のコンストラクタを定義
func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバー作成
	s := grpc.NewServer()

	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	// サーバーリフレクションの設定
	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Panicln("stopping gRPC server...")
	s.GracefulStop()
}