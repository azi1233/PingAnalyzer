package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	pb "github.com/azi1233/PingAnalyzer/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(dstIP string, count int64, start bool) {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial(":8080", opt)
	if err != nil {
		log.Fatalf("starting grpc-server as client error because :%v\n", err)
	}

	defer cc.Close()

	client := pb.NewPingServiceClient(cc)
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	request := &pb.PingRequestMessage{}
	if start {
		request = &pb.PingRequestMessage{DstIP: dstIP, Start: true, Count: count}
	} else {
		request = &pb.PingRequestMessage{Start: false, Id: count}
	}
	stream, err := client.PingFunc(context.Background(), request)

	if err != nil {
		panic(err)
	}

	for {
		reply, err := stream.Recv()

		if err == io.EOF {
			//	openShellAndWrite(int(reply.Id), "Stream finished")
			//time.Sleep(2 * time.Second)
			//	closeShell(int(reply.Id))
			fmt.Printf("Stream finished")

			break
		}

		if err != nil {
			log.Println("error in reading stream")
		}
		//	outputString := fmt.Sprintf("id is %d, ttl is %d,rtt is%f\n", reply.Id, reply.Ttl, reply.Time)
		//	openShellAndWrite(int(reply.Id), outputString)
		fmt.Printf("id=%d,ttl is %d,rtt is %f\n", reply.Id, reply.Ttl, reply.Time)

	}

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var line string

	for {
		scanner.Scan() //ping IP. STOP ID
		line = scanner.Text()
		if line == "" {
			continue
		}

		line := strings.Fields(line)
		if len(line) < 2 {
			fmt.Println("bad input, more than 1 argument needed")
			continue
		}

		var dstIP string
		var err error
		var count int
		switch line[0] {
		case "ping":
			{
				dstIP = line[1]
				if len(line) > 2 {
					count, err = strconv.Atoi(line[2])
					if err != nil {
						log.Println("bad count input, enter your command again")
						continue
					}

				} else {
					count = 0
				}

				go NewClient(dstIP, int64(count), true)
			}

		case "stop":
			{
				count, err := strconv.Atoi(line[1]) //it is id
				if err != nil {
					log.Println("bad id input, enter your command again")
					continue
				}
				go NewClient(line[0], int64(count), false)
			}
		default:
			{
				fmt.Println("bad input, you can use \"ping\" or \"stop\" in these ways")
				fmt.Println("ping <IP Address>, ping <IP Address> <Count>")
				fmt.Println("stop <ID number>")
			}

		}

	}

}

/*func main() {

	scanner := bufio.NewScanner(os.Stdin)
	var line string
	scanner.Scan()
	line = scanner.Text()
	num, _ := strconv.Atoi(line)

	scanner.Scan()
	line = scanner.Text()
	wait, _ := strconv.Atoi(line)

	for i := 0; i < num; i++ {
		go NewClient("1.1.1.1", 0, true)

	}
	go func() {
		for {
			scanner.Scan()
			line = scanner.Text()
			if line == "panic" {
				for i := 1; i <= num; i++ {
					NewClient("stop", 0, false)
				}
				panic(errors.New("stop"))
			}
		}

	}()
	for i := 0; i < wait; i++ {
		time.Sleep(time.Second)
	}
	for i := 1; i <= num; i++ {
		NewClient("stop", int64(i), false)
	}

}*/
