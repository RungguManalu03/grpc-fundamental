package main

import (
	"context"
	"fmt"
	"golang-grpc/golang-grpc/proto"
	"golang-grpc/repository"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	dbUrl := "postgres://postgres:03Manalu01@35.223.1.56:5432/db_golang_grpc"
	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userServer := &repository.UserRepository{DB: dbPool}
	proto.RegisterUserServiceServer(grpcServer, userServer)

	fmt.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
