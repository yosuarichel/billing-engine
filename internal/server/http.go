package server

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"

	httpHandler "github.com/yosuarichel/billing-engine/handler/http"
	"github.com/yosuarichel/billing-engine/pkg/config"
)

func StartHTTP(ctx context.Context, cfg *config.AppConfig, handler *httpHandler.HttpHandler) {
	// cfg := config.GetAppCfg()
	klog.Infof("Starting HTTP Server on :%d ...", cfg.HTTPPort)

	h := server.Default(server.WithHostPorts(fmt.Sprintf(":%d", cfg.HTTPPort)))
	httpHandler.RegisterRoutes(h, handler)

	h.Spin()
}
