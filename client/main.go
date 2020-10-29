package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"tracer"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

var (
	err    error
	closer io.Closer
)

func main() {
	tracer.Tracer, closer, err = tracer.NewTracer("servicename", "hongning:6831")
	if err != nil {
		log.Fatal("tracer,NewTracer error(%v)", err)
	}
	defer closer.Close()

	client := &http.Client{Transport: &nethttp.Transport{}}
	span := tracer.Tracer.StartSpan("client")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	req, err := http.NewRequest("GET", "http://hongning:8002/getIP", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(tracer.Tracer, req)
	defer ht.Finish()

	res, err := client.Do(req)
	if err != nil {
		onError(span, err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		onError(span, err)
		return
	}
	log.Printf("Received result: %s\n", body)
}

func onError(span opentracing.Span, err error) {
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Fatalf("client(%v)\n", err)
}
