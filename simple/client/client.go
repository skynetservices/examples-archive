package main

import (
	"fmt"
	"github.com/skynetservices/skynet2/client"
	_ "github.com/skynetservices/zkmanager"
)

func main() {
	// (any version, any region, any host)
	service := client.GetService("TestService", "", "", "")

	// This on the other hand will fail if it can't find a service to
	// connect to
	in := map[string]interface{}{
		"data": "Upcase me!!",
	}
	out := map[string]interface{}{}
	err := service.Send(nil, "Upcase", in, &out)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out["data"].(string))
}
