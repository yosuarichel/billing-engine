package server

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"go.opentelemetry.io/otel/attribute"

	"github.com/hertz-contrib/obs-opentelemetry/tracing"
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

	// Init HTTP server
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", cfg.HTTPPort)),
		server.WithTraceLevel(stats.LevelDetailed),
		tracer,
	)
	h.Use(tracing.ServerMiddleware(traceCfg))
	httpHandler.RegisterRoutes(h, handler)

	h.Spin()
}
