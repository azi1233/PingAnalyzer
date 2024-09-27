package main

import (
	"bufio"
	"context"

	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/azi1233/PingAnalyzer/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(dstIP string) {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial(":8080", opt)
	if err != nil {
		log.Fatalf("starting grpc-server as client error because :%v\n", err)
	}

	defer cc.Close()

	client := pb.NewPingServiceClient(cc)
	request := &pb.PingRequest{DstIp: dstIP}
	reply, err := client.PingServiceFunc(context.Background(), request)

	if err != nil {
		panic(err)
	}

	fmt.Println(reply)

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var line string

	for {
		scanner.Scan()
		line = scanner.Text()
		line = strings.Fields(line)[1]
		go NewClient(line)

	}

}
