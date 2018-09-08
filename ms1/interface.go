package ms1

import (
	"ms/delay"
	"ms/config"
	"log"
	"net/rpc"
)

type Args int

type Microservice1 int

//* Synchronous call
func (m *Microservice1) Get(args Args, reply *string) error {
	// Can be asynchronous
	log.Print("Getting Ms2 info... ")
	str1 := GetMs2()
	log.Print("Getting Ms3 info... ")
	str2 := GetMs3()
	
	*reply = str1 + " " + str2 + "!"
	
	delay.Sleep(config.Ms1Delay, config.Ms1DelayVar)

	return nil
}
//*/

/* Asynchronous call
func (m *Microservice1) Get(args Args, reply *string) error {
	ch1 := make(chan string)
	ch2 := make(chan string)
	var str1 string
	var str2 string
	log.Print("Getting Ms2 info... ")
	log.Print("Getting Ms3 info... ")
	go func() {
		ch1 <- GetMs2()
	}()

	go func() {
		ch2 <- GetMs3()
	}()

	for i:=0; i < 2; i++ {
		select {
		case str1 = <-ch1:
		case str2 = <-ch2:
		}
	}
	
	*reply = str1 + " " + str2 + "!"
	
	delay.Sleep(config.Ms1Delay, config.Ms1DelayVar)

	return nil
}
*/

func GetMs2() string {
	microservice2, err := rpc.DialHTTP("tcp", "0.0.0.0" + config.PortMs2)
	if err != nil {
		log.Fatal("Ms1: ", err)
	}

	var reply string

	err = microservice2.Call("Microservice2.Get", 0, &reply)
	if err != nil {
		log.Fatal("Ms1: ", err)
	}

	return reply
}

func GetMs3() string {
	microservice3, err := rpc.DialHTTP("tcp", "0.0.0.0" + config.PortMs3)
	if err != nil {
		log.Fatal("Ms1: ", err)
	}

	var reply string

	err = microservice3.Call("Microservice3.Get", 0, &reply)
	if err != nil {
		log.Fatal("Ms1: ", err)
	}

	return reply
}
