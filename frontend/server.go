package frontend

import (
	"net"
	"net/rpc"
	"net/http"
	"ms/config"
	"log"
)

func NewServer() {
	serv := rpc.NewServer()
	front := new(Frontend)
	serv.Register(front)

	// 
	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	//

	serv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	// 
	http.DefaultServeMux = oldMux
	// 

	listener, err := net.Listen("tcp", config.PortFrontend)

	if err != nil {
		log.Panic("Frontend: ", err)
	}

	go http.Serve(listener, mux)
}

