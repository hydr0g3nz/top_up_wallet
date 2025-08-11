package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/wallet_topup_system/config"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/controller"
	"github.com/hydr0g3nz/wallet_topup_system/internal/adapter/repository/postgresql/repository"
	usecase "github.com/hydr0g3nz/wallet_topup_system/internal/application"
	"github.com/hydr0g3nz/wallet_topup_system/internal/infrastructure"
)

func main() {
	// Load configuration
	config := config.LoadFromEnv()
	// Setup logger
	logger, err := infrastructure.NewLogger(config.IsProduction())
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// // Connect to database
	db, err := infrastructure.ConnectDB(&config.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error()})
	}

	// Run migrations
	if err := infrastructure.MigrateDB(db); err != nil {
		logger.Fatal("Failed to run database migrations", map[string]interface{}{
			"error": err.Error()})
	}

	// Seed database with initial data
	if err := infrastructure.SeedDB(db); err != nil {
		logger.Fatal("Failed to seed database", map[string]interface{}{
			"error": err.Error()})
	}

	cache := infrastructure.NewRedisClient(config.Cache)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	// txRepo := repository.NewDBTransactionRepository(db)
	txManager := repository.NewTxManagerGorm(db)

	// Initialize use cases
	walletUsecase := usecase.NewWalletUsecase(userRepo, transactionRepo, walletRepo, cache, txManager, logger, *config)

	// Setup server
	server := infrastructure.NewFiber(infrastructure.ServerConfig{
		Address:      config.Server.Port,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
	})
	registerRoutes(server, walletUsecase)
	// Start server
	logger.Info("Starting server", map[string]interface{}{"port": config.Server.Port})

	if err := server.Listen(fmt.Sprintf(":%s", config.Server.Port)); err != nil {
		logger.Fatal("Failed to start server", map[string]interface{}{
			"error": err.Error()})
	}
}

// registerRoutes registers all API routes
func registerRoutes(
	app *fiber.App,
	walletUseCase usecase.WalletUsecase,

) {
	// Setup API routes
	api := app.Group("/api/v1")
	walletController := controller.NewWalletController(walletUseCase)
	walletController.RegisterRoutes(api)
}
