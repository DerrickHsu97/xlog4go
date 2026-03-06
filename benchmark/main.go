package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	logger "git.xiaojukeji.com/shield-arch/dlog4go"
)

var (
	wg sync.WaitGroup
)

func zerologBench() {
	start := time.Now()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				logger.Info("info=%v", "infotest")
			}
			wg.Done()

		}()
	}
	wg.Wait()
	logger.Close()
	cost := time.Since(start)
	fmt.Printf("zerolog time cost %v \n", cost)
}

func main() {
	if err := logger.SetupLogWithConf("./log.json"); err != nil {
		panic(err)
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	zerologBench()
}
