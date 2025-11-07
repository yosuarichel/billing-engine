package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexServer "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"go.opentelemetry.io/otel/attribute"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	rpcHandler "github.com/yosuarichel/billing-engine/handler/rpc"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine/billingengineservice"
	"github.com/yosuarichel/billing-engine/pkg/config"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

func StartRPC(ctx context.Context, cfg *config.AppConfig, handler *rpcHandler.RpcHandler) {
	klog.Infof("Starting RPC Server on :%d ...", cfg.RPCPort)
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

	// Init HTTP health check server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.HealthCheckPort),
			Handler: mux,
		}
		klog.Infof("HTTP health check listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.CtxErrorf(ctx, "HTTP health check failed: %+v", err)
		}
	}()

	// Init Consul register
	r, err := consul.NewConsulRegister(
		fmt.Sprintf("%s:%d", cfg.Consul.Host, cfg.Consul.Port),
		consul.WithCheck(&consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://billing-engine-rpc:%d/health", cfg.HealthCheckPort),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "5m",
		}))
	if err != nil {
		klog.Fatalf("Failed to create Consul register: %+v", err)
	}

	// Init RPC server
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", cfg.RPCPort))
	svr := billingengineservice.NewServer(handler,
		kitexServer.WithRegistry(r),
		kitexServer.WithRegistryInfo(&registry.Info{
			ServiceName: appName,
			Weight:      1, // weights must be greater than 0 in consul,else received error and exit.
		}),
		kitexServer.WithServiceAddr(addr),
		kitexServer.WithSuite(tracing.NewServerSuite()),
		kitexServer.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: appName,
				Tags: map[string]string{
					"env": cfg.Env,
				},
			},
		),
	)
	if err := svr.Run(); err != nil {
		klog.CtxErrorf(ctx, "RPC service start error=%+v", err)
	}
}
