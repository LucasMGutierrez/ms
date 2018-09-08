package main

import (
	"ms/frontend"
	"ms/ms1"
	"ms/ms2"
	"ms/ms3"
)

func main() {
	// TODO frontend server
	frontend.NewServer()
	ms1.NewServer()
	ms2.NewServer()
	ms3.NewServer()
	//time.Sleep(5 * time.Second)
	//fmt.Println(frontend.Get())
}



