package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {

	// RAM, 填满内存，每个核心，每0.5s刷写内存数据
	maxLen := 100000000
	ramMap := make(map[int]string, maxLen)
	for i := 0; i < maxLen-1; i++ {
		ramMap[i] = fmt.Sprint(i)
	}

	fmt.Println(len(ramMap))

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	someMapMutex := sync.RWMutex{}
	for i := 0; i < cores; i++ {
		go func() {
			for {
				someMapMutex.Lock()
				ramMap[rand.Intn(maxLen-1)] = fmt.Sprint(rand.Intn(math.MaxInt32))
				someMapMutex.Unlock()
				time.Sleep(time.Millisecond * 500)
			}
		}()
	}

	// 把CPU搞起来
	for i := 0; i < cores; i++ {
		go func() {
			for {
				// time.Sleep(time.Microsecond)
				rand.Intn(math.MaxInt32)
			}
		}()
	}
	// 写日志监测
	for {
		time.Sleep(time.Second)
		fileWriting()
	}
}

// single core to generate random int.
func fileWriting() {
	f, err := os.OpenFile("./cri.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "[cri]", log.LstdFlags)
	logger.Println("IO happend...")
}
