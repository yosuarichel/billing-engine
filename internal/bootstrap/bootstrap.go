package bootstrap

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/yosuarichel/billing-engine/pkg/config"
	"github.com/yosuarichel/billing-engine/pkg/infra/db"
	"github.com/yosuarichel/billing-engine/pkg/infra/external"
	"github.com/yosuarichel/billing-engine/pkg/infra/redis"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

type AppDeps struct {
	Cfg          *config.AppConfig
	Ctx          context.Context
	OTelProvider provider.OtelProvider
}

func Init() *AppDeps {
	ctx := context.Background()

	// Load the config
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		klog.CtxFatalf(ctx, "failed to load config: %v", err)
	}
	klog.CtxInfof(ctx, "Config loaded successfully: %+v", cfg)

	// Init DB
	db.MustInitDB(ctx)

	// Init Redis
	redis.MustInitRedis(ctx)

	// Init sonyflake
	utils.InitSonyflakeCluster()

	// Init external clients
	external.InitBillingCustomerClient()

	return &AppDeps{
		Cfg: cfg,
		Ctx: ctx,
	}
}
