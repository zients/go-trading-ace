package main

import (
	"context"
	"database/sql"
	"fmt"
	"trading-ace/config"
	"trading-ace/controllers"
	"trading-ace/helpers"
	"trading-ace/logger"
	"trading-ace/routes"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var ctx = context.Background()

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

	if err := db.Ping(); err != nil {
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

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func NewGinServer() *gin.Engine {
	r := gin.Default()
	return r
}

func SetupServer(
	r *gin.Engine,
	config *config.Config,
	homeRoutes routes.IHomeRoutes,
) {
	homeRoutes.RegisterHomeRoutes()

	r.Run(fmt.Sprintf(":%d", config.Server.Port))
}

func main() {
	app := fx.New(
		fx.Provide(

			// Base
			NewGinServer,
			NewDB,
			NewRedis,
			config.LoadConfig,
			logger.NewLogrusLogger,

			// Controllers
			controllers.NewHomeController,

			// Repositories

			// Routes
			routes.NewHomeRoutes,

			// Helper
			helpers.NewRedisHelper,
		),
		fx.Invoke(SetupServer),
	)

	app.Start(ctx)
}
