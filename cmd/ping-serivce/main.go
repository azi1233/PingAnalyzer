package main

import (
	"context"
	"fmt"

	//"os/exec"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/azi1233/PingAnalyzer/api/pb"
)

type pingInfo struct {
	SrcIP       string
	DstIP       string
	Interval    int
	RequestTime time.Time
	ch          *chan statStruct
}

type server struct {
	pb.PingServiceServer
	store         (map[uint64]pingInfo)
	NextID        uint64 //always increasing
	lastID        uint64
	UnallocatedID []uint64
}

func (s *server) PingServiceFunc(c context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	fmt.Println("//////////////")

	//if start request
	tempPing := pingInfo{
		SrcIP:       req.SrcIp,
		DstIP:       req.DstIp,
		Interval:    int(req.Interval),
		RequestTime: time.Now(),
	}
	if len(s.UnallocatedID) != 0 {

		if s.NextID > s.UnallocatedID[0] {
			//s.store[s.UnallocatedID[0]] = tempPing
			s.lastID = s.UnallocatedID[0]
			s.UnallocatedID = s.UnallocatedID[1:]

		} else {
			//	s.store[s.NextID] = tempPing
			s.lastID = s.NextID
			s.NextID++
		}

	} else {
		s.lastID = s.NextID
		s.NextID++
	}

	tempCh := PingStart(tempPing)
	tempPing.ch = tempCh
	s.store[s.lastID] = tempPing

	msg := <-*s.store[s.lastID].ch

	return &pb.PingReply{Status: true, Time: int32(msg.sent), SendTtl: int32(msg.rec), ReceivedTtl: int32(msg.lossPercentage)}, nil

}

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Couldn't generate listen config")
	}

	s := grpc.NewServer()
	pb.RegisterPingServiceServer(s, &server{NextID: 1, UnallocatedID: []uint64{}, store: make(map[uint64]pingInfo)})
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("couldn't start setver because:%v\n", err)
	}

}
