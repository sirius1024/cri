package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
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
	needAlloc := freeMem - freeSwap - 3
	// 占用策略：freeMem-3GB
	fmt.Printf("总内存：%v，总SWAP：%v，可用内存：%v，可用SWAP：%v，预计占用内存（free mem - freeSwap - 3g）：%vGB\n", totalMem, totalSwap, freeMem, freeSwap, needAlloc)
	fmt.Println("开始申请内存...")

	ramMap := make(map[uint64]string, needAlloc)
	for i := uint64(0); i < needAlloc; i++ {
		ramMap[i] = strings.Repeat("a", gb)
	}

	fmt.Println("开始循环交换内存数据...")
	cores := runtime.NumCPU()
	mapMutex := sync.RWMutex{}
	for i := 0; i < cores; i++ {
		go func() {
			mapMutex.Lock()
			ramMap[rand.Uint64()] = strings.Repeat(fmt.Sprint(rand.Intn(9)), gb)
			mapMutex.Unlock()
			time.Sleep(time.Millisecond * 500)
		}()
	}

	fmt.Println("开始占用CPU资源……")
	// 把CPU搞起来
	for i := 0; i < cores; i++ {
		go func() {
			for {
				rand.Intn(math.MaxInt32)
			}
		}()
	}

	fmt.Println("启动成功，日志位于./cri.log")
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
