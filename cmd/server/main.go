package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/azi1233/PingAnalyzer/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial(":8080", opt)
	if err != nil {
		log.Fatalf("starign grpc-server as client error because :%v\n", err)
	}

	defer cc.Close()

	client := pb.NewPingServiceClient(cc)
	request := &pb.PingRequest{SrcIp: "1.1.1.1", DstIp: "2.2.2.2", Id: 1}
	reply, err := client.PingServiceFunc(context.Background(), request)

	if err != nil {
		panic(err)
	}

	fmt.Println(reply)

}
