package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	pb "imageelev/api/v1"
	srv "imageelev/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const (
	portgRPC uint32 = 3144
	portHTTP uint32 = 8081
)

var (
	isProfile = flag.Bool("debug", false, "Enable pprof server")
)

func main() {
	flag.Parse()
	log.Println("Starting ImagesByElevationServer ...")
	servEndPoint := fmt.Sprintf("localhost:%d", portgRPC)
	lis, err := net.Listen("tcp", servEndPoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	/// Create HTTP REST API
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcServ := grpc.NewServer()
	pb.RegisterImageByElevationServer(grpcServ, srv.NewService())
	go func() {
		log.Println("Staring gRPC server on " + servEndPoint)
		log.Fatalln(grpcServ.Serve(lis))
	}()

	/// Create debug pprof server
	if *isProfile {
		log.Println("Starting pprof server on localhost:8088")
		go func() {
			log.Fatalln(http.ListenAndServe(":8088", nil))
		}()
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterImageByElevationHandlerFromEndpoint(ctx, mux, servEndPoint, opts)

	if err != nil {
		log.Fatalf("Can`t RegisterImageByElevationHandlerFromEndpoint err:%v\n", err)
	}
	servHTTP := fmt.Sprintf("localhost:%v", portHTTP)
	log.Println("Starting HTTP server on " + servHTTP)
	log.Fatalln(http.ListenAndServe(servHTTP, mux))
}
