package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sample-grpc-server/service"
	"syscall"

	"sample-grpc-server/config"
	"sample-grpc-server/database"
	"sample-grpc-server/interceptor"
	"sample-grpc-server/pb"
	"sample-grpc-server/server"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	config.LoadConfig()
}

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(interceptor.RecoveryFunc),
	}

	qer := database.NewQuery(db)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			interceptor.AuthInterceptor(qer),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	)

	pb.RegisterBackendServiceServer(s, server.NewServer(qer, service.NewHash(), service.NewAuth()))
	reflection.Register(s)

	go func() {
		log.Println("listening server with port:8080")
		if err := s.Serve(ln); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	<-quit

	log.Println("stopping gRPC server...")
	s.GracefulStop()
	log.Println("grpc server shutdown completed")
}
