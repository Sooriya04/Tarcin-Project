package main

import (
	"backend/internal/config"
	"backend/internal/database"
	grpcservice "backend/internal/grpc"
	"backend/internal/grpc/pb"
	"backend/internal/repository"
	"context"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	log.Printf("[gRPC] --> %s", info.FullMethod)

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		log.Printf("[gRPC] <-- %s | ERROR: %v | Duration: %s", info.FullMethod, err, duration)
	} else {
		log.Printf("[gRPC] <-- %s | SUCCESS | Duration: %s", info.FullMethod, duration)
	}

	return resp, err
}

func main() {
	// Setup MultiWriter to write logs to stdout and server.log
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open server.log: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// 1. Load configuration (Environment variables)
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection
	db := database.InitDB(cfg.DBUrl)
	defer db.Close()

	// 3. Setup Repository
	internRepo := repository.NewInternRepository(db)

	// 4. Setup and Start gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on :50051: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	perfService := grpcservice.NewPerformanceService(db, internRepo)
	pb.RegisterPerformanceServiceServer(grpcServer, perfService)

	// Enable reflection so grpcurl works without needing proto files
	reflection.Register(grpcServer)

	log.Println("gRPC Server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
