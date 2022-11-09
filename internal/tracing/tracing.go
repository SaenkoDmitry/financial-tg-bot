package tracing

import (
	"flag"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"log"
)

func InitTracing() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	serviceName := flag.String("service", "financial-tg-bot", "financial telegram bot for cost accounting")
	_, err := cfg.InitGlobalTracer(*serviceName)
	if err != nil {
		log.Fatal("Cannot init tracing", zap.Error(err))
	}
}
