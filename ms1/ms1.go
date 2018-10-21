package main

import (
	"net/http"
	"github.com/uber/jaeger-client-go/examples/ms/config"
	"github.com/uber/jaeger-client-go/examples/ms/tracing"
	"github.com/uber/jaeger-client-go/examples/ms/delay"
	"github.com/uber/jaeger-client-go/examples/ms/http"
	"github.com/opentracing/opentracing-go"
	"log"
)

func main() {
	NewServer()
}

func NewServer() {
	tracer, closer := tracing.InitJaeger("Ms1")
	defer closer.Close()
	init := false

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("Ms1", opentracing.ChildOf(spanCtx))
		defer span.Finish()

		str1 := xhttp.Get(span, "Ms2", config.PortMs2)
		str2 := xhttp.Get(span, "Ms3", config.PortMs3)

		delay.Sleep(config.Ms1Delay, config.Ms1DelayVar, &init)
		reply := str1 + str2


		w.Write([]byte(reply))
	})

	log.Println("Ms1 server up")
	log.Fatal(http.ListenAndServe(config.PortMs1, nil))
}

