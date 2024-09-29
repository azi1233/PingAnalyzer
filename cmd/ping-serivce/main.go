package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/azi1233/PingAnalyzer/api/pb"
	"google.golang.org/grpc"
)

type PingInfo struct {
	DstIP     string
	Count     int64
	Interval  int64
	StartTime time.Time
	EndTime   time.Time
	ch        *chan streamResponse
	chStop    *chan bool //stopping the pinger
	stop      bool       //stopping the for loop
}

type Server struct {
	pb.PingServiceServer
	DataSotre map[uint64]PingInfo
	LastId    uint64
	mu        *sync.Mutex
}

var wg *sync.WaitGroup

func (s *Server) PingFunc(msg *pb.PingRequestMessage, stream pb.PingService_PingFuncServer) error {
	if msg.Start {

		ch, chStop, err := MyPing(msg.DstIP, msg.Count)

		if err != nil {
			panic(fmt.Errorf("ping chanle dosen't return: %v\n", err))
		}
		temp := PingInfo{
			DstIP:     msg.DstIP,
			Count:     msg.Count,
			Interval:  msg.Interval,
			StartTime: time.Now(),
			ch:        ch,
			chStop:    chStop,
		}
		s.mu.Lock()
		s.LastId++
		id := s.LastId
		s.DataSotre[s.LastId] = temp
		s.mu.Unlock()
		//go func(id uint64, info PingInfo) {
		for {
			msg := <-*s.DataSotre[id].ch
			pong := &pb.PongReplyStream{
				Result: false,
				Time:   float32(msg.rtt.Milliseconds()),
				Ttl:    int32(msg.ttl),
				Status: msg.status,
				Id:     int64(id),
				DstIP:  s.DataSotre[id].DstIP,
			}
			if err := stream.Send(pong); err != nil {
				// You need to handle this error later
			}

			if s.DataSotre[id].stop {
				break
			}
		}
		//	}(s.LastId-1,temp)
	} else {
		temp, ok := s.DataSotre[uint64(msg.Id)]
		if !ok {
			fmt.Printf("%d id not available\n", msg.Id)

		} else {
			temp.stop = true
			s.DataSotre[uint64(msg.Id)] = temp
			*s.DataSotre[uint64(msg.Id)].chStop <- true
			//wg.Add(1)
			//	wg.Wait()//further analysis needed
			s.mu.Lock()
			delete(s.DataSotre, uint64(msg.Id))
			s.mu.Unlock()
		}

	}
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(fmt.Errorf("listener error: %v\n", err))
	}

	s := grpc.NewServer()
	pb.RegisterPingServiceServer(s, &Server{DataSotre: make(map[uint64]PingInfo), LastId: 0, mu: &sync.Mutex{}})

	if err := s.Serve(listener); err != nil {
		panic(fmt.Errorf("server start error: %v\n", err))
	}
}
