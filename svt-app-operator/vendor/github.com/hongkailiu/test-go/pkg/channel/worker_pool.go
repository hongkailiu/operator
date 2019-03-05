// ref. https://stackoverflow.com/questions/18267460/how-to-use-a-goroutine-pool
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hongkailiu/test-go/pkg/lib/logger"
	"github.com/hongkailiu/test-go/pkg/lib/util"
)

var (
	log = logger.Logger
)

func worker(index int, linkChan chan job, resultChan chan result, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	log.WithField("index", index).Info("work started")
	defer wg.Done()
	log.Info(fmt.Sprintf("Getting job ..."))
	for j := range linkChan {
		n := util.GetRandomInt(10)
		if n > 3 {
			doSilly()
		}
		log.Info(fmt.Sprintf("doing #%s, #%s, %d", j.param1, j.param2, n))
		time.Sleep(time.Duration(n) * time.Second)
		resultChan <- result{j.param1, j.param2}
		log.Info(fmt.Sprintf("finished #%s, #%s, %d", j.param1, j.param2, n))
	}
	log.WithField("index", index).Info("work finished")
}
func doSilly() {
	//panic("panic panic panic")
	// or
	//n = 0
	//log.Info(3/n)
}

type job struct {
	param1 string
	param2 string
}

type result struct {
	param1 string
	param2 string
}

func main() {
	jobSlice := make([]job, 50)
	for i := 0; i < 50; i++ {
		jobSlice[i] = job{fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", i+2)}
	}

	poolSize := 10
	//buffered channel could make sending/receiving smoother
	//but it also works without buffering
	//jobCh := make(chan job)
	jobCh := make(chan job, poolSize)
	rCh := make(chan result)
	wg := new(sync.WaitGroup)

	go func() {
		for r := range rCh {
			fmt.Printf("result is %s\n", r.param1)
		}
	}()

	// Adding routines to workgroup and running then
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go worker(i, jobCh, rCh, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, j := range jobSlice {
		log.WithField("aaa", "bbb").Info("send job")
		jobCh <- j
	}

	// Closing channel (waiting in goroutines won't continue any more)
	close(jobCh)

	// Waiting for all goroutines to finish (otherwise they die as main routine dies)
	wg.Wait()

	close(rCh)

}
