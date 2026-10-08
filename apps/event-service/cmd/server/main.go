package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/event-ticketing/event-service/internal/config"
	"github.com/event-ticketing/event-service/internal/handler"
	"github.com/event-ticketing/event-service/internal/middleware"
	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/repository"
	"github.com/event-ticketing/event-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Structured JSON logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	cfg := config.Load()

	slog.Info("starting event-service",
		"port", cfg.Server.Port,
		"env", cfg.Server.Env,
		"db_host", cfg.Database.Host,
		"db_name", cfg.Database.DBName,
	)

	// Database connection with retry
	db, err := connectDB(cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("connected to database")

	// Auto-migrate
	if err := db.AutoMigrate(
		&model.Venue{},
		&model.Event{},
		&model.TicketType{},
		&model.Seat{},
	); err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	slog.Info("database migration completed")

	// Repositories
	eventRepo := repository.NewEventRepository(db)
	venueRepo := repository.NewVenueRepository(db)
	ticketTypeRepo := repository.NewTicketTypeRepository(db)
	seatRepo := repository.NewSeatRepository(db)

	// Services
	eventSvc := service.NewEventService(eventRepo, venueRepo)
	venueSvc := service.NewVenueService(venueRepo)
	ticketTypeSvc := service.NewTicketTypeService(ticketTypeRepo)
	seatSvc := service.NewSeatService(seatRepo, ticketTypeRepo)

	// Handlers
	eventHandler := handler.NewEventHandler(eventSvc)
	venueHandler := handler.NewVenueHandler(venueSvc)
	ticketTypeHandler := handler.NewTicketTypeHandler(ticketTypeSvc)
	seatHandler := handler.NewSeatHandler(seatSvc)
	healthHandler := handler.NewHealthHandler(db)

	// Gin setup
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger("event-service"))

	// Health endpoints
	healthHandler.RegisterRoutes(r)

	// API v1 routes
	v1 := r.Group("/api/v1")
	eventHandler.RegisterRoutes(v1)
	venueHandler.RegisterRoutes(v1)
	ticketTypeHandler.RegisterRoutes(v1)
	seatHandler.RegisterRoutes(v1)

	// Graceful shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("event-service listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down event-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	slog.Info("event-service stopped")
}

func connectDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxOpenConns(25)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)
			return db, nil
		}

		slog.Warn("waiting for database...",
			"attempt", i+1,
			"error", err,
		)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after 30 attempts: %w", err)
}
