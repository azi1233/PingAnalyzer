package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-ping/ping"
)

const GOOS string = runtime.GOOS

type streamResponse struct {
	rtt    time.Duration
	ttl    int
	status bool
}

func MyPing(ip string, count int64) (*chan streamResponse, *chan bool, error) {
	pinger, err := ping.NewPinger(ip)
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

	if count == 0 {
		pinger.Count = -1
	} else {
		pinger.Count = int(count)

	}
	c := make(chan streamResponse, 10)
	cStop := make(chan bool, 1)

	pinger.OnRecv = func(p *ping.Packet) {
		temp := streamResponse{
			ttl:    p.Ttl,
			rtt:    p.Rtt,
			status: true,
		}
		c <- temp
	}

	go func() {
		msg := <-cStop
		if msg {
			pinger.Stop()
			//wg.Done()//should add another chanle to do that

		}
	}()

	go func() {
		err := pinger.Run()
		if err != nil {
			panic(fmt.Errorf("the pinger.Run error:%v/n", err))
		}

	}()

	return &c, &cStop, nil
}
