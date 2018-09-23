package main

import (
	"fmt"
	"github.com/uber/jaeger-client-go/examples/ms/config"
	"github.com/uber/jaeger-client-go/examples/ms/tracing"
	"github.com/uber/jaeger-client-go/examples/ms/http"
	"github.com/opentracing/opentracing-go"
)

func main() {
	fmt.Println(Request())
}

func Request() string {
	/* Start tracer */
	tracer, closer := tracing.InitJaeger("Hello-World")
    defer closer.Close()
    opentracing.SetGlobalTracer(tracer)

    span := tracer.StartSpan("Request")
    span.SetTag("Client", "")

    defer span.Finish()

    //ctx := opentracing.ContextWithSpan(context.Background(), span)

	return xhttp.Get(span, "Frontend", config.PortFrontend)
}
