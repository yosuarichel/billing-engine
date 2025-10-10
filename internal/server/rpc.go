package server

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexServer "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"go.opentelemetry.io/otel/attribute"

	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	rpcHandler "github.com/yosuarichel/billing-engine/handler/rpc"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine/billingengineservice"
	"github.com/yosuarichel/billing-engine/pkg/config"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

func StartRPC(ctx context.Context, cfg *config.AppConfig, handler *rpcHandler.RpcHandler) {
	klog.Infof("Starting RPC Server on :%d ...", cfg.RPCPort)
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(utils.GetAppName()),
		provider.WithExportEndpoint("otel-collector:4317"),
		provider.WithInsecure(),
		provider.WithResourceAttribute(attribute.String("env", cfg.Env)),
	)
	defer func() {
		if err := p.Shutdown(ctx); err != nil {
			klog.CtxErrorf(ctx, "failed to shutdown OTel provider: %+v", err)
		}
	}()

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", cfg.RPCPort))
	svr := billingengineservice.NewServer(handler,
		kitexServer.WithServiceAddr(addr),
		kitexServer.WithSuite(tracing.NewServerSuite()),
		kitexServer.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: utils.GetAppName(),
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
