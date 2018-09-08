package ms3

import (
	"net"
	"net/rpc"
	"net/http"
	"ms/config"
	"log"
)

func NewServer() {
	serv := rpc.NewServer()
	microservice := new(Microservice3)
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

	listener, err := net.Listen("tcp", config.PortMs3)

	if err != nil {
		log.Panic("Ms3: ", err)
	}

	http.Serve(listener, mux)
}

