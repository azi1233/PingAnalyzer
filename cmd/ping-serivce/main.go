package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/azi1233/PingAnalyzer/api/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.PingServiceServer
}

func (s *server) PingServiceFunc(c context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	fmt.Println(req)
	fmt.Println("//////////////")

	return &pb.PingReply{
		Status:  true,
		Time:    req.Id,
		SendTtl: 1101342,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Couldn't generate listen config")
	}

	s := grpc.NewServer()
	pb.RegisterPingServiceServer(s, &server{})
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("couldn't start setver because:%v\n", err)
	}

}
