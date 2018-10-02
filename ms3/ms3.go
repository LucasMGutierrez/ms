package main

import (
	"net/http"
	"github.com/uber/jaeger-client-go/examples/ms/config"
	"github.com/uber/jaeger-client-go/examples/ms/tracing"
	"github.com/uber/jaeger-client-go/examples/ms/delay"
	"github.com/opentracing/opentracing-go"
	"log"
)

func main() {
	NewServer()
}

func NewServer() {
	tracer, closer := tracing.InitJaeger("Ms3")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("Ms3", opentracing.ChildOf(spanCtx))
		defer span.Finish()

		reply := "World!"
		delay.Sleep(config.Ms3Delay, config.Ms3DelayVar)

		w.Write([]byte(reply))
	})

	log.Println("Ms3 server up")
	log.Fatal(http.ListenAndServe(config.PortMs3, nil))
}

