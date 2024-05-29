package main

import (
	"context"
	"erp-logger-service/logger"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const grpcPort = "50001"

type LogServer struct {
	logger.UnimplementedLoggerServiceServer

	Models *Models
}

func (server *LogServer) WriteLog(ctx context.Context, req *logger.LogRequest) (*logger.LogResponse, error) {
	err := server.Models.LogEntry.Insert(&LogEntry{
		Event:         req.Event,
		CallerService: req.CallerService,
		Timestamp:     req.Timestamp,
		Details:       req.Details,
	})

	return &logger.LogResponse{
		WriteSuccess: err != nil,
	}, err
}

func (app *Config) startGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))

	if err != nil {
		log.Fatalf("Failed to listen to grpc: %v", err)
	}

	s := grpc.NewServer()

	logger.RegisterLoggerServiceServer(s, &LogServer{
		Models: &Models{},
	})

	log.Printf("grpc server started on port %s", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen to grpc port %v", grpcPort)
	}
}
