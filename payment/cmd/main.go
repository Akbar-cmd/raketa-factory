package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1API "github.com/Akbar-cmd/raketa-factory/payment/internal/api/payment/v1"
	paymentService "github.com/Akbar-cmd/raketa-factory/payment/internal/service/payment"
	paymentV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/payment/v1"
)

const grpcAddr = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcAddr))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcAddr)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
