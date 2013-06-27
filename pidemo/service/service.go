package main

import (
	"github.com/skynetservices/skynet2"
	"github.com/skynetservices/skynet2/log"
	"github.com/skynetservices/skynet2/service"
	"github.com/skynetservices/zkmanager"
	"os"
	"strings"
	"time"
)

type PiDemoService struct {
	Led *LED
}

func (s *PiDemoService) Registered(service *service.Service) {
	s.Led.Blue(true)
}
func (s *PiDemoService) Unregistered(service *service.Service) {
	s.Led.Blue(false)
}

func (s *PiDemoService) Started(service *service.Service) {}
func (s *PiDemoService) Stopped(service *service.Service) {
}

func NewPiDemoService() *PiDemoService {
	r := &PiDemoService{
		Led: NewLED(),
	}
	return r
}

func (s *PiDemoService) Upcase(requestInfo *skynet.RequestInfo, in map[string]interface{}, out map[string]interface{}) (err error) {
	out["data"] = strings.ToUpper(in["data"].(string))
	return
}

func main() {
	log.SetLogLevel(log.DEBUG)
	skynet.SetServiceManager(zkmanager.NewZookeeperServiceManager(os.Getenv("SKYNET_ZOOKEEPER"), 1*time.Second))

	testService := NewPiDemoService()

	config, _ := skynet.GetServiceConfig()

	if config.Name == "" {
		config.Name = "PiDemoService"
	}

	if config.Version == "unknown" {
		config.Version = "1"
	}

	if config.Region == "unknown" {
		config.Region = "Clearwater"
	}

	service := service.CreateService(testService, config)

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
