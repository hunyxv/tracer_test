module main

go 1.15

replace tracer => ./tracer

require (
	github.com/opentracing-contrib/go-stdlib v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	tracer v0.0.0-00010101000000-000000000000 // indirect
)
