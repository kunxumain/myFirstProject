package common

import (
	"io"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
// 	cfg := &config.Configuration{
// 		ServiceName: serviceName,
// 		Sampler: &config.SamplerConfig{
// 			Type:  jaeger.SamplerTypeConst,
// 			Param: 1,
// 		},
// 		Reporter: &config.ReporterConfig{
// 			BufferFlushInterval: 1 * time.Second,
// 			LogSpans:            true,
// 			LocalAgentHostPort:  addr,
// 		},
// 	}
// 	return cfg.NewTracer()
// }

const vmwareAddr = "192.168.71.200"

func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: 1 * time.Second,
			LogSpans:            true,
			// 方法1: 使用 HTTP Collector 直接发送
			CollectorEndpoint: "http://" + vmwareAddr + ":14268/api/traces",
			// 注释掉 LocalAgentHostPort，因为我们使用 CollectorEndpoint
			// LocalAgentHostPort:  addr,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Jaeger tracer 初始化成功 - 服务: %s, 端点: %s", serviceName, "http://192.168.71.172:14268/api/traces")
	return tracer, closer, nil
}
