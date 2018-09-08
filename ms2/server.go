package ms2

import (
	"net"
	"net/rpc"
	"net/http"
	"ms/config"
	"log"
)

func NewServer() {
	serv := rpc.NewServer()
	microservice := new(Microservice2)
	serv.Register(microservice)

	// 
	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	//

	serv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	// 
	http.DefaultServeMux = oldMux
	// 

	listener, err := net.Listen("tcp", config.PortMs2)

	if err != nil {
		log.Panic("Ms2: ", err)
	}

	go http.Serve(listener, mux)
}

