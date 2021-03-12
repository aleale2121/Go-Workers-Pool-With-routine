package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	id       int
	randomNo int
}
type Factorial struct {
	job    Job
	result uint64
}

var jobs = make(chan Job, 10)
var factorials = make(chan Factorial, 10)

func calculateFactorial(number int) uint64 {
	if number >= 1 {
		return uint64(number) * calculateFactorial(number-1)
	} else {
		return 1
	}
}

func work(wg *sync.WaitGroup) {
	for job := range jobs {
		factorial := Factorial{job, calculateFactorial(job.randomNo)}
		time.Sleep(2 * time.Second)
		factorials <- factorial
	}
	wg.Done()
}
func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go work(&wg)
	}
	wg.Wait()
	close(factorials)
}
func allocate(noOfJobs int) {
	for i := 0; i < noOfJobs; i++ {
		randomNo := rand.Intn(20)
		job := Job{i, randomNo}
		jobs <- job
	}
	close(jobs)
}
func displayFactorials(done chan bool) {
	for factorial := range factorials {
		fmt.Printf("Job id %d, input random no %d , factorial %v\n", factorial.job.id, factorial.job.randomNo, factorial.result)
	}
	done <- true
}
func main() {
	startTime := time.Now()
	noOfJobs := 20
	go allocate(noOfJobs)
	done := make(chan bool)

	go displayFactorials(done)

	noOfWorkers := 10
	createWorkerPool(noOfWorkers)

	<-done
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}
