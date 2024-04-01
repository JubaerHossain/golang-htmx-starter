package core

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"strconv"

	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/cache"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/config"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/database"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/http"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/render"
	html "github.com/JubaerHossain/golang-htmx-starter/static"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type App struct {
	BuildVersion string
	Echo         *echo.Echo
	HttpPort     int
	PublicFS     fs.FS
	cache        cache.CacheService
	db           database.DB
	logger       *zap.Logger
}

func StartApp() (*App, error) {
	// Initialize environment variables
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	if err := config.LoadConfig(); err != nil {
		return nil, errors.New("failed to load configuration: " + err.Error())
	}

	// Start the database connection routine
	dbCh := make(chan database.DB)
	cacheCh := make(chan cache.CacheService)

	go func() {
		db, err := database.ConnectDB()
		if err != nil {
			logger.Error("failed to connect to database", zap.Error(err))
			dbCh <- nil
			return
		}
		if config.GlobalConfig.Migrate {
			if err := database.MigrateDB(db); err != nil {
				logger.Error("failed to migrate database", zap.Error(err))
				dbCh <- nil
				return
			}
		}
		dbCh <- db
	}()

	// Start the Redis cache initialization routine
	go func() {
		cacheService, err := cache.NewRedisCacheService(context.Background())
		if err != nil {
			logger.Error("failed to initialize cache service", zap.Error(err))
			cacheCh <- nil
			return
		}
		cacheCh <- cacheService
	}()

	// Wait for both database and cache routines to finish
	db := <-dbCh
	cacheService := <-cacheCh

	if db == nil || cacheService == nil {
		return nil, errors.New("failed to initialize database or cache service")
	}

	// Use default values if environment variables are not set
	httpPort, _ := strconv.Atoi(config.GlobalConfig.AppPort)

	app := &App{
		HttpPort:     httpPort,
		Echo:         echo.New(),
		PublicFS:     html.PublicFS,
		BuildVersion: config.GlobalConfig.AppEnv,
		cache:        cacheService,
		db:           db,
		logger:       logger,
	}

	// Initialize template renderer
	renderer, err := render.NewRenderer(app.PublicFS)
	if err != nil {
		return nil, errors.New("failed to create template renderer: " + err.Error())
	}

	// Pimp the echo instance
	e := app.Echo
	e.Renderer = renderer
	e.StaticFS("/", html.StaticFS)
	e.Use(http.SetCacheControl)

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.GlobalConfig.SessionSecret))))


	return app, nil
}


