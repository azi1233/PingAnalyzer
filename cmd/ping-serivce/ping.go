package main

import (
	//"fmt"
	//"github.com/azi1233/PingAnalyzer/api/pb"
	"runtime"

	"github.com/go-ping/ping"
)

const GOOS string = runtime.GOOS

type statStruct struct {
	sent           int
	rec            int
	lossPercentage float64
}

func PingStart(req pingInfo) *chan statStruct {

	pinger, err := ping.NewPinger(req.DstIP)
	if err != nil {
		panic(err)
	}
	if GOOS == "windows" {
		pinger.SetPrivileged(true)

	}

	if GOOS == "linux" {
		pinger.SetPrivileged(true)
		//run your binary as root

	}
	pinger.Count = 4

	c := make(chan statStruct, 10)

	//pinger.OnRecv = func(stats *ping.Packet) {
	//	fmt.Println(stats)
	//	c <- fmt.Sprintln(stats)
	//}

	pinger.OnFinish = func(stats *ping.Statistics) {
		tempStat := statStruct{
			lossPercentage: stats.PacketLoss,
			sent:           stats.PacketsSent,
			rec:            stats.PacketsRecv,
		}
		c <- tempStat
	}

	go func() {
		err := pinger.Run()
		if err != nil {
			panic(err)
		}

	}()

	return &c

}
