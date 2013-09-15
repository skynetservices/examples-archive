package main

import (
	"expvar"
	"flag"
	"fmt"
	"github.com/skynetservices/skynet2/client"
	"github.com/skynetservices/skynet2/config"
	_ "github.com/skynetservices/zkmanager"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var totalRequests = expvar.NewInt("total-requests")
var successfulRequests = expvar.NewInt("successful-requests")

var requests int
var goMaxProcs int
var piDemoClient client.ServiceClientProvider

func main() {
	flagset := flag.NewFlagSet("pidemo", flag.ContinueOnError)
	flagset.IntVar(&requests, "requests", 10, "number of concurrent requests")
	flagset.IntVar(&goMaxProcs, "maxprocs", 1, "GOMAXPROCS")

	runtime.GOMAXPROCS(goMaxProcs)

	c := make(chan os.Signal, 1)
	quitChan := make(chan bool, 1)
	requestChan := make(chan string, requests*3)
	workerQuitChan := make(chan bool, 1)
	workerWaitGroup := new(sync.WaitGroup)

	go watchSignals(c, quitChan)

	pidemoArgs, _ := config.SplitFlagsetFromArgs(flagset, os.Args[1:])
	flagset.Parse(pidemoArgs)

	piDemoClient = client.GetService("PiDemoService", "", "", "")

	startTime := time.Now().UnixNano()
	fmt.Printf("Starting %d Workers\n", requests)
	for i := 0; i < requests; i++ {
		go worker(requestChan, workerWaitGroup, workerQuitChan)
	}

	requestNum := 0

	for {
		select {
		case <-quitChan:
			for i := 0; i < requests; i++ {
				workerQuitChan <- true
			}

			workerWaitGroup.Wait()
			stopTime := time.Now().UnixNano()

			successful, _ := strconv.Atoi(successfulRequests.String())
			total, _ := strconv.Atoi(totalRequests.String())

			failed := total - successful

			percentSuccess := int(float64(successful) / float64(total) * 100)
			percentFailed := int(float64(failed) / float64(total) * 100)

			runtime := (stopTime - startTime) / 1000000
			rqps := float64(total) / (float64(runtime) / 1000)

			fmt.Printf("======================================")
			fmt.Printf("======================================\n")
			fmt.Printf("Completed in %d Milliseconds, %f Requests/s\n",
				runtime, rqps)
			fmt.Printf("\nTotal Requests: %d, Successful: %d (%d%%)",
				total, successful, percentSuccess)
			fmt.Printf(", Failed: %d (%d%%)\n\n", failed, percentFailed)
			return

		default:
			requestChan <- "pidemo"
			requestNum++
		}
	}
}

func worker(requestChan chan string, waitGroup *sync.WaitGroup,
	quitChan chan bool) {

	waitGroup.Add(1)

	for {
		select {
		case <-quitChan:
			waitGroup.Done()
			return

		case service := <-requestChan:
			totalRequests.Add(1)

			switch service {
			case "pidemo":

				randString := strconv.FormatUint(uint64(rand.Uint32()), 35)
				randString = randString + randString + randString

				in := map[string]interface{}{
					"data": randString,
				}

				fmt.Println("Sending TestService request: " + in["data"].(string))

				out := map[string]interface{}{}
				err := piDemoClient.Send(nil, "Upcase", in, &out)

				upper := strings.ToUpper(randString)
				if err == nil && out["data"].(string) == upper {
					successfulRequests.Add(1)
					fmt.Println("TestService returned: " + out["data"].(string))
				}
			}

		}
	}
}

func watchSignals(c chan os.Signal, quitChan chan bool) {
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV,
		syscall.SIGSTOP, syscall.SIGTERM)

	for {
		select {
		case sig := <-c:
			switch sig.(syscall.Signal) {
			// Trap signals for clean shutdown
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT,
				syscall.SIGSEGV, syscall.SIGSTOP, syscall.SIGTERM:

				quitChan <- true
				return
			}
		}
	}
}
