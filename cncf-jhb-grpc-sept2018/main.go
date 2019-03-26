package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"sync"

	"play/cncf/gen"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	//add fluentd logger
	"github.com/fluent/fluent-logger-golang/fluent"

	//add metrics
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"

	//add tracing
	"go.opencensus.io/trace"

	//add metrics exporter
	"net/http"

	"go.opencensus.io/exporter/prometheus"

	//add tracing exporter
	"go.opencensus.io/exporter/jaeger"
)

var fluentdLogger *fluent.Fluent

const fluentdTag = "demo.cncf"

type KV struct {
	sync.Mutex
	store map[string]string

	hits    int64
	nonHits int64
}

func (k *KV) Get(ctx context.Context, in *kv.GetRequest) (*kv.GetResponse, error) {
	_, span := trace.StartSpan(ctx, "(*kv).Get")
	defer span.End()

	k.Lock()
	defer k.Unlock()

	log.Printf("get: %s", in.Key)
	resp := new(kv.GetResponse)
	if val, ok := k.store[in.Key]; ok {
		k.hits++

		span.Annotate([]trace.Attribute{
			trace.Int64Attribute("hits", int64(k.hits)),
		}, "getRequestCount")

		resp.Value = val

		return resp, nil
	}

	k.nonHits++

	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("nonHits", int64(k.nonHits)),
	}, "getRequestCount")

	return nil, status.Errorf(codes.NotFound, "key '%s' not set", in.Key)

}

func (k *KV) Set(ctx context.Context, in *kv.SetRequest) (*kv.SetResponse, error) {
	_, span := trace.StartSpan(ctx, "(*kv).Set")
	defer span.End()

	log.Printf("set: %s = %s", in.Key, in.Value)
	k.Lock()
	defer k.Unlock()

	k.store[in.Key] = in.Value

	return &kv.SetResponse{Ok: true}, nil
}

func NewKVStore() (kv *KV) {
	return &KV{
		store: make(map[string]string),
	}
}

func main() {
	var err error
	port := flag.Int("port", 8080, "grpc port")

	flag.Parse()

	// setup fluentd logger
	fluentdLogger, err = fluent.New(fluent.Config{FluentPort: 24224, FluentHost: "127.0.0.1"})
	if err != nil {
		log.Fatal(err)
	}
	defer fluentdLogger.Close()

	// TLS
	// prepare TLS Config
	tlsCert := "certs/demo.crt"
	tlsKey := "certs/demo.key"
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		log.Fatal(err)
	}

	// create gRPC TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	// TLS END

	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatalf("Failed to register gRPC server views: %v", err)
	}

	if err := createAndRegisterExporters(); err != nil {
		log.Fatalf("Failed to register exporters: %v", err)
	}

	//srv := grpc.NewServer() // without stats/metrics
	srv := grpc.NewServer(
		grpc.StatsHandler(new(ocgrpc.ServerHandler)),
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	kv.RegisterKVServer(srv, NewKVStore())

	log.Printf("starting grpc on :%d\n", *port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv.Serve(lis)
}

func createAndRegisterExporters() error {
	// For demo purposes, set this to always sample.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// 1. Prometheus
	prefix := "kv"
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: prefix,
	})
	if err != nil {
		return fmt.Errorf("Failed to create Prometheus exporter: %v", err)
	}
	view.RegisterExporter(pe)
	// We need to expose the Prometheus collector via an endpoint /metrics
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		log.Fatal(http.ListenAndServe(":9888", mux))
	}()

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:6831",
		Endpoint:      "http://localhost:14268",
		ServiceName:   "cncf-jhb-demo",
	})
	if err == nil {
		// On success, register it as a trace exporter
		trace.RegisterExporter(je)
	}

	return err
}

// general unary interceptor function
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	data := map[string]string{
		"source":   "unaryInterceptor",
		"method":   info.FullMethod,
		"duration": fmt.Sprintf("%v", time.Since(start)),
	}

	if err != nil {
		data["error"] = err.Error()
	}

	log.Printf("request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err) //logging

	if err := fluentdLogger.PostWithTime(fluentdTag, time.Now(), data); err != nil {
		fmt.Println("fluentd logger err:", err)
	}

	return h, err
}

/*
fluentd -c etc\td-agent\td-agent.conf
fluentd code:
	//add fluentd logger
	"github.com/fluent/fluent-logger-golang/fluent"


	logger, err := fluent.New(fluent.Config{FluentPort: 24224, FluentHost: "127.0.0.1"})
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	tag := "cncf-demo"

	logger.PostWithTime(tag, time.Now(), "testing, 456...")

*/
