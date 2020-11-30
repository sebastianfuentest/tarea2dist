package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"papa.com/chat"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
	s := chat.Server{}
	grpcServerChat := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServerChat, &s)

	if err := grpcServerChat.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server over port 9000: %v", err)
	}
	/*
		var conn1 *grpc.ClientConn
		conn1, err1 := grpc.Dial(":9001", grpc.WithInsecure())

		if err1 != nil {
			log.Fatalf("Could not connect: %s", err1)
		}
		defer conn1.Close()

			var conn2 *grpc.ClientConn
			conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())

			if err2 != nil {
				log.Fatalf("Could not connect: %s", err2)
			}
			defer conn2.Close()

			var conn3 *grpc.ClientConn
			conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())

			if err3 != nil {
				log.Fatalf("Could not connect: %s", err3)
			}
			defer conn3.Close()
	*/
}
