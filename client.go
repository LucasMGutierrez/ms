package main

import (
	"ms/config"
	"fmt"
	"net/rpc"
	"log"
)

func main() {
	fmt.Println(Request())
}

func Request() string {
	front, err := rpc.DialHTTP("tcp", "0.0.0.0" + config.PortFrontend)
	if err != nil {
		log.Fatal("Client: ", err)
	}

	var reply string

	err = front.Call("Frontend.Get", 0, &reply)
	if err != nil {
		log.Fatal("Client: ", err)
	}

	return reply
}