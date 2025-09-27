package server

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	kitexServer "github.com/cloudwego/kitex/server"

	rpcHandler "github.com/yosuarichel/billing-engine/handler/rpc"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine/billingengineservice"
	"github.com/yosuarichel/billing-engine/pkg/config"
)

func StartRPC(ctx context.Context, cfg *config.AppConfig, handler *rpcHandler.RpcHandler) {
	klog.Infof("Starting RPC Server on :%d ...", cfg.RPCPort)

	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", cfg.RPCPort))
	svr := billingengineservice.NewServer(handler,
		kitexServer.WithServiceAddr(addr),
	)

	if err := svr.Run(); err != nil {
		klog.CtxErrorf(ctx, "RPC service start error=%+v", err)
	}
}
