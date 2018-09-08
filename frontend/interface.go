package frontend

import (
	"ms/delay"
	"ms/config"
	"log"
	"net/rpc"
)

type Args int

type Frontend int

func (f *Frontend) Get(args Args, reply *string) error {
	// Can be asynchronous
	log.Print("Getting Ms1 info... ")
	str := GetString()
	
	*reply = str
	
	delay.Sleep(config.FrontendDelay, config.FrontendDelayVar)

	return nil
}

func GetString() string {
	microservice1, err := rpc.DialHTTP("tcp", "0.0.0.0" + config.PortMs1)
	if err != nil {
		log.Fatal("Frontend: ", err)
	}

	var reply string

	err = microservice1.Call("Microservice1.Get", 0, &reply)
	if err != nil {
		log.Fatal("Frontend: ", err)
	}

	return reply
}