package main

import (
	"net/http"
	"github.com/uber/jaeger-client-go/examples/ms/config"
	"github.com/uber/jaeger-client-go/examples/ms/delay"
	"github.com/uber/jaeger-client-go/examples/ms/http"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/examples/ms/tracing"
)

func main() {
	NewServer()
}

func NewServer() {
	tracer, closer := tracing.InitJaeger("Frontend")
	defer closer.Close()
	init := false

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("Frontend", opentracing.ChildOf(spanCtx))
		defer span.Finish()

		reply := xhttp.Get(span, "Ms1", config.PortMs1)
		delay.Sleep(config.FrontendDelay, config.FrontendDelayVar, &init)

		w.Write([]byte(reply))
	})

	log.Println("Frontend server up")
	log.Fatal(http.ListenAndServe(config.PortFrontend, nil))
}


