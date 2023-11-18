package jaeger

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

type Config struct {
	ServiceName string `yaml:"service_name" env:"JAEGER_SERVICE_NAME"`
	Address     string `yaml:"address" env:"JAEGER_ADDRESS"`
	LogSpans    bool   `yaml:"log_spans" env:"JAEGER_LOG_SPANS"`
}

func InitJaeger(cfg Config) (opentracing.Tracer, io.Closer, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1.0,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.LogSpans,
			LocalAgentHostPort: cfg.Address,
		},
	}

	return jaegerCfgInstance.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
}
