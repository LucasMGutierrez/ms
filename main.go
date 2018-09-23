package main

import (
	"github.com/uber/jaeger-client-go/examples/ms/frontend"
	"github.com/uber/jaeger-client-go/examples/ms/ms1"
	"github.com/uber/jaeger-client-go/examples/ms/ms2"
	"github.com/uber/jaeger-client-go/examples/ms/ms3"
)

func main() {
	// TODO frontend server
	go frontend.NewServer()
	go ms1.NewServer()
	go ms2.NewServer()
	ms3.NewServer()
	//time.Sleep(5 * time.Second)
	//fmt.Println(frontend.Get())
}