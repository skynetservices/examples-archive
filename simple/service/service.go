package main

import (
	"github.com/skynetservices/skynet"
	"github.com/skynetservices/skynet/log"
	"github.com/skynetservices/skynet/service"
	_ "github.com/skynetservices/zkmanager"
	"strings"
)

type TestService struct{}

func (s *TestService) Registered(service *service.Service)   {}
func (s *TestService) Unregistered(service *service.Service) {}
func (s *TestService) Started(service *service.Service)      {}
func (s *TestService) Stopped(service *service.Service) {
}

func NewTestService() *TestService {
	r := &TestService{}
	return r
}

func (s *TestService) Upcase(requestInfo *skynet.RequestInfo, in map[string]interface{}, out map[string]interface{}) (err error) {
	out["data"] = strings.ToUpper(in["data"].(string))
	return
}

func main() {
	testService := NewTestService()

	serviceInfo := skynet.NewServiceInfo("TestService", "1.0.0")

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
