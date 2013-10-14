package main

import (
	"github.com/skynetservices/skynet"
	"github.com/skynetservices/skynet/log"
	"github.com/skynetservices/skynet/service"
	"github.com/skynetservices/skynet/stats"
	_ "github.com/skynetservices/zkmanager"
	"strings"
	"time"
)

var led *LED = NewLED()
var registered bool

const (
	OFF = iota
	RED
	GREEN
	BLUE
)

type LedReporter struct {
	blinkChan chan int
}

func (r *LedReporter) UpdateHostStats(host string, stats stats.Host) {}
func (r *LedReporter) MethodCalled(method string)                    {}
func (r *LedReporter) MethodCompleted(method string, duration int64, err error) {
	if err != nil {
		r.blinkChan <- RED
	} else {
		r.blinkChan <- GREEN
	}
}

func (r *LedReporter) watch() {
	var t *time.Ticker = time.NewTicker(1 * time.Minute)

	for {
		select {
		case color := <-r.blinkChan:
			if t != nil {
				t.Stop()
			}

			led.Off()

			switch color {
			case RED:
				led.Red(true)
			case BLUE:
				led.Blue(true)
			case GREEN:
				led.Green(true)
			}

			t = time.NewTicker(100 * time.Millisecond)
		case <-t.C:
			led.Off()

			if registered {
				led.Blue(true)
			} else {
				led.Off()
			}
		}
	}
}

func NewLedReporter() (r *LedReporter) {
	r = &LedReporter{
		blinkChan: make(chan int, 10000),
	}
	go r.watch()

	return
}

type PiDemoService struct {
}

func (s *PiDemoService) Registered(service *service.Service) {
	registered = true
	led.Blue(true)
}
func (s *PiDemoService) Unregistered(service *service.Service) {
	registered = false
	led.Blue(false)
}

func (s *PiDemoService) Started(service *service.Service) {}
func (s *PiDemoService) Stopped(service *service.Service) {
}

func NewPiDemoService() *PiDemoService {
	r := &PiDemoService{}
	return r
}

func (s *PiDemoService) Upcase(requestInfo *skynet.RequestInfo, in map[string]interface{}, out map[string]interface{}) (err error) {
	out["data"] = strings.ToUpper(in["data"].(string))
	return
}

func main() {
	stats.AddReporter(NewLedReporter())

	testService := NewPiDemoService()

	serviceInfo := skynet.NewServiceInfo("PiDemoService", "1.0.0")

	service := service.CreateService(testService, serviceInfo)

	// handle panic so that we remove ourselves from the pool in case
	// of catastrophic failure
	defer func() {
		service.Shutdown()
		if err := recover(); err != nil {
			log.Panic("Unrecovered error occured: ", err)
		}
	}()

	waiter := service.Start()

	// waiting on the sync.WaitGroup returned by service.Start() will
	// wait for the service to finish running.
	waiter.Wait()
}
