package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"papa.com/chat"
)

//openLis is ponerse a escuchar
func openLis() {

	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("Failed to listen on port 9003: %v", err)
	}
	s := chat.Server{}
	grpcServerChat := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServerChat, &s)

	if err := grpcServerChat.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server over port 9003: %v", err)
	}
}

//func mandarProp(nombreL string, )

func main() {
	//ponerse a escuchar
	openLis()
	//c := chat.NewChatServiceClient(conn)
	//comunicarse con los demas datanodes to do :p

}
