package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	// total ram
	const gb = 1024 * 1024 * 1024
	mem := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(mem)
	if err != nil {
		panic(err)
	}
	totalMem := mem.Totalram / gb
	totalSwap := mem.Totalswap / gb
	freeMem := mem.Freeram / gb
	freeSwap := mem.Freeswap / gb
	needAlloc := uint64(0)
	if freeMem >= freeSwap {
		needAlloc = freeMem - freeSwap
	}
	if needAlloc <= 4 {
		needAlloc = 1
	}
	// 占用策略：freeMem-3GB
	fmt.Printf("总内存：%v，总SWAP：%v，可用内存：%v，可用SWAP：%v，预计占用内存freeMem-freeSwap - 4g：%vGB\n", totalMem, totalSwap, freeMem, freeSwap, needAlloc)
	fmt.Println("开始申请内存...")

	ramMap := make(map[uint64]string, needAlloc)
	for i := uint64(0); i < needAlloc; i++ {
		ramMap[i] = strings.Repeat("a", gb)
	}
	cores := runtime.NumCPU()

	fmt.Println("开始占用CPU资源……")
	// 把CPU搞起来
	for i := 0; i < cores-1; i++ {
		go func() {
			for {
				rand.Intn(math.MaxInt32)
			}
		}()
	}

	fmt.Println("启动成功，日志位于./cri.log")
	// 写日志监测
	go func() {
		for {
			time.Sleep(time.Second)
			fileWriting("IO happened.")
		}
	}()

	// 内存刷写
	ch := make(chan string)
	ch2 := make(chan string)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			newData := strings.Repeat(fmt.Sprint(rand.Intn(9)), gb)
			ch <- newData
			runtime.GC()
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 3)
			newData := strings.Repeat(fmt.Sprint(rand.Intn(9)), gb/2)
			ch2 <- newData
			runtime.GC()
		}
	}()

	fmt.Println("开始循环交换内存数据...")
	for {
		select {
		case <-ch:
			fileWriting("New Ram data from goroutine channel 1, size: 1G, refreshed.")
			rdNum := rand.Uint64()
			rdNumReplaced := rand.Uint64()
			ramMap[rdNum] = strings.Replace(ramMap[rdNum], ramMap[rdNum], ramMap[rdNumReplaced], -1)
		case <-ch2:
			fileWriting("new RAM data from goroutine channel 2, size: 512M, refreshed.")
			rdNum := rand.Uint64()
			rdNumReplaced := rand.Uint64()
			ramMap[rdNum] = strings.Replace(ramMap[rdNum], ramMap[rdNum], ramMap[rdNumReplaced], -1)
		}
	}
}

// single core to generate random int.
func fileWriting(s string) {
	f, err := os.OpenFile("./cri.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "[cri]", log.LstdFlags)
	logger.Println(s)
}
