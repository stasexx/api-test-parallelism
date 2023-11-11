package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	counter   int
	mutex     sync.Mutex
	waitGroup sync.WaitGroup
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/counter", counterHandler)

	waitGroup.Add(1)
	go startServer()

	sendSequentialRequests()

	sendParallelRequests()

	waitGroup.Wait()
}

const (
	numRequests = 100
	baseURL     = "http://localhost:8080/counter"
)

func sendSequentialRequests() {
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		response, err := http.Get(baseURL)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		_, _ = ioutil.ReadAll(response.Body)
		_ = response.Body.Close()
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Sequential: Elapsed time: %s\n", elapsedTime)
}

func sendParallelRequests() {
	startTime := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var counter int

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			response, err := http.Get(baseURL)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			_, _ = ioutil.ReadAll(response.Body)
			_ = response.Body.Close()

			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Parallel: Elapsed time: %s, Counter: %d\n", elapsedTime, counter)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is a simple parallel web server!")
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	counter++
	fmt.Fprintf(w, "Counter: %d", counter)
}

func startServer() {
	defer waitGroup.Done()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
