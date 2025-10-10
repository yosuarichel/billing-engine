package db

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	db  *gorm.DB
	cfg = config.GetAppCfg()
)

func MustInitDB(ctx context.Context) {
	init := func(configDB *config.DBConfig) *gorm.DB {
		var dsn string
		var dialector gorm.Dialector

		switch configDB.Driver {
		case "postgres":
			dsn = fmt.Sprintf(
				"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
				configDB.User, configDB.Password, configDB.Host, configDB.Port, configDB.DBName, configDB.SSLMode,
			)
			dialector = postgres.Open(dsn)

		case "mysql":
			// charset, parseTime, and loc are common options
			dsn = fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				configDB.User, configDB.Password, configDB.Host, configDB.Port, configDB.DBName,
			)
			dialector = mysql.Open(dsn)

		default:
			klog.CtxErrorf(ctx, "unsupported driver: %s", configDB.Driver)
			return nil
		}

		gormDB, err := gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			klog.CtxErrorf(ctx, "failed to connect database: %w", err)
			return nil
		}

		sqlDB, err := gormDB.DB()
		if err != nil {
			panic(err)
		}

		sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.DB.MaxConnLifetime) * time.Minute)
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.DB.MaxConnIdleTime) * time.Minute)

		if err := sqlDB.Ping(); err != nil {
			panic(err)
		}

		klog.CtxInfof(ctx, "✅ Connected to Database at %s", dsn)
		return gormDB
	}

	db = init(&cfg.DB)

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
