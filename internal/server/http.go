package server

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	hUtils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"go.opentelemetry.io/otel/attribute"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/registry/consul"
	httpHandler "github.com/yosuarichel/billing-engine/handler/http"
	"github.com/yosuarichel/billing-engine/pkg/config"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

func StartHTTP(ctx context.Context, cfg *config.AppConfig, handler *httpHandler.HttpHandler) {
	klog.Infof("Starting HTTP Server on :%d ...", cfg.HTTPPort)
	appName := utils.GetAppName()

	// Init OpenTelemetry provider
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(appName),
		provider.WithExportEndpoint("otel-collector:4317"),
		provider.WithInsecure(),
		provider.WithResourceAttribute(attribute.String("env", cfg.Env)),
	)
	defer func() {
		if err := p.Shutdown(ctx); err != nil {
			klog.CtxErrorf(ctx, "Failed to shutdown OTel provider: %+v", err)
		}
	}()

	// Init HTTP tracer
	tracer, traceCfg := tracing.NewServerTracer()

	// Init Consul register
	consulClient, err := consulapi.NewClient(&consulapi.Config{
		Address: fmt.Sprintf("%s:%d", cfg.Consul.Host, cfg.Consul.Port),
	})
	if err != nil {
		klog.Fatalf("Failed to create Consul client: %+v", err)
		return
	}

	r := consul.NewConsulRegister(consulClient, consul.WithCheck(&consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://billing-engine-http:%d/health", cfg.HTTPPort),
		Interval:                       "10s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "5m",
	}))

	// Init HTTP server
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", cfg.HTTPPort)),
		server.WithRegistry(r, &registry.Info{
			ServiceName: utils.GetAppName(),
			Addr:        hUtils.NewNetAddr("tcp", fmt.Sprintf("%s:%d", utils.GetLocalIP(), cfg.HTTPPort)),
			Weight:      10,
			Tags:        map[string]string{"env": cfg.Env},
		}),
		server.WithTraceLevel(stats.LevelDetailed),
		tracer,
	)
	h.Use(tracing.ServerMiddleware(traceCfg))
	httpHandler.RegisterRoutes(h, handler)

	h.Spin()
}
