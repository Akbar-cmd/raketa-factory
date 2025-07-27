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

	inventoryV1API "github.com/Akbar-cmd/raketa-factory/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/part"
	inventoryService "github.com/Akbar-cmd/raketa-factory/inventory/internal/service/part"
	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
)

const grpcAddr = 50051

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

	// Регистрируем сервис
	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Рефлексия для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC InventoryService server listening on %d\n", grpcAddr)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
