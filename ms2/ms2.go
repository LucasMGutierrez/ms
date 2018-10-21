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
	tracer, closer := tracing.InitJaeger("Ms2")
	defer closer.Close()
	init := false

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("Ms1", opentracing.ChildOf(spanCtx))
		defer span.Finish()

		reply := "Hello "
		delay.Sleep(config.Ms2Delay, config.Ms2DelayVar, &init)

		w.Write([]byte(reply))
	})

	log.Println("Ms2 server up")
	log.Fatal(http.ListenAndServe(config.PortMs2, nil))
}

