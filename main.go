package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hunyxv/tracer_test/tracer"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

func main() {
	var (
		err    error
		closer io.Closer
	)

	tracer.Tracer, closer, err = tracer.NewTracer("servicename", "hongning:6831")
	if err != nil {
		log.Fatal("tracer,NewTracer error(%v)", err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer.Tracer)

	http.HandleFunc("/getIP", getIP)
	log.Printf("Starting server on port %d", 8002)
	err = http.ListenAndServe(
		fmt.Sprintf(":%d", 8002),
		nethttp.Middleware(tracer.Tracer, http.DefaultServeMux),
	)
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}

func getIP(w http.ResponseWriter, r *http.Request) {
	log.Print("Received getIP request")

	client := &http.Client{Transport: &nethttp.Transport{}}

	var span opentracing.Span
	if parent := opentracing.SpanFromContext(r.Context()); parent != nil {
		pctx := parent.Context()
		fmt.Println(pctx)
		fmt.Println(r.Context())
		if tracer := tracer.Tracer; tracer != nil {
			span = tracer.StartSpan("getIP", opentracing.ChildOf(pctx))
		}
	} else {
		span = tracer.Tracer.StartSpan("getIP")
	}
	defer span.Finish()

	span.SetTag(string(ext.Component), "getIP")

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	req, err := http.NewRequest("GET", "http://icanhazip.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)
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
	io.WriteString(w, fmt.Sprintf("ip %s", body))
}

func onError(span opentracing.Span, err error) {
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Fatalf("client(%v)\n", err)
}
