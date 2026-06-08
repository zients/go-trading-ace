package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"
	"trading-ace/config"
	"trading-ace/controllers"
	"trading-ace/helpers"
	"trading-ace/logger"
	"trading-ace/middlewares"
	"trading-ace/repositories"
	"trading-ace/routes"
	"trading-ace/services"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

const startupTimeout = 10 * time.Second

func NewDB(config *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func NewRedis(config *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func NewGinServer(config *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Timeout(config.Server.RequestTimeout()))
	return r
}

func NewAppContext(lc fx.Lifecycle) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			cancel()
			return nil
		},
	})

	return ctx
}

func SetupServer(
	lc fx.Lifecycle,
	appCtx context.Context,
	r *gin.Engine,
	logger logger.ILogger,
	config *config.Config,
	ethereumService services.IEthereumService,
	homeRoutes routes.IHomeRoutes,
	campaignRoutes routes.ICampaignRoutes,
) {
	homeRoutes.RegisterHomeRoutes()
	campaignRoutes.RegisterCampaignRoutes()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: r,
	}

	ethereumCtx, cancelEthereum := context.WithCancel(appCtx)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", server.Addr)
			if err != nil {
				cancelEthereum()
				return fmt.Errorf("failed to listen on %s: %w", server.Addr, err)
			}

			go func() {
				if err := ethereumService.SubscribeEthereumSwap(ethereumCtx); err != nil {
					logger.Error("Ethereum swap subscription stopped: %v", err)
				}
			}()

			go func() {
				if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server stopped unexpectedly: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancelEthereum()
			return server.Shutdown(ctx)
		},
	})
}

func main() {
	app := fx.New(
		fx.StartTimeout(startupTimeout),
		fx.StopTimeout(startupTimeout),
		fx.Provide(

			// Base
			NewGinServer,
			NewDB,
			NewRedis,
			config.LoadConfig,
			logger.NewLogrusLogger,

			// Controllers
			controllers.NewHomeController,
			controllers.NewCampaignController,

			// Repositories
			repositories.NewTaskRepository,
			repositories.NewTaskHistoryRepository,

			// Routes
			routes.NewHomeRoutes,
			routes.NewCampaignRoutes,

			// Services
			services.NewCampaignService,
			services.NewEthereumService,

			// Helper
			helpers.NewRedisHelper,
			NewAppContext,
		),
		fx.Invoke(SetupServer),
	)

	app.Run()
}
