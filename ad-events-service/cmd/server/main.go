package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"ad-events-service/internal/handler"
	"ad-events-service/internal/repository"
	"ad-events-service/internal/service"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment variables")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	logger.Info("Successfully connected to the database!")

	campaignRepo := repository.NewRepository(pool)
	bannerRepo := repository.NewBannerRepository(pool)
	statsRepo := repository.NewStatsRepository(pool)
	eventRepo := repository.NewEventRepository(pool)

	campaignService := service.NewCampaignService(campaignRepo)
	bannerService := service.NewBannerService(bannerRepo, campaignRepo)
	statsService := service.NewStatsService(statsRepo, campaignRepo)
	eventService := service.NewEventService(eventRepo, campaignRepo, bannerRepo)
	adService := service.NewAdService(bannerRepo, campaignRepo)

	campaignHandler := handler.NewCampaignHandler(campaignService, logger)
	bannerHandler := handler.NewBannerHandler(bannerService, logger)
	statsHandler := handler.NewStatsHandler(statsService)
	eventHandler := handler.NewEventHandler(eventService)
	adHandler := handler.NewAdHandler(adService)

	router := gin.Default()
	router.GET("/campaigns", campaignHandler.GetAllCampaigns)
	router.GET("/campaigns/:id", campaignHandler.GetCampaignByID)
	router.GET("/campaigns/name", campaignHandler.GetCampaignByName)
	router.POST("/campaigns", campaignHandler.CreateCampaign)
	router.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)
	router.DELETE("/campaigns/:id", campaignHandler.DeleteCampaign)
	router.PATCH("/campaigns/:id", campaignHandler.PatchCampaign)

	router.GET("/campaigns/:id/banners", bannerHandler.GetAllBannersByCampaignId)
	router.GET("/banners/:id", bannerHandler.GetBannerByID)
	router.POST("/campaigns/:id/banners", bannerHandler.CreateBanner)
	router.PUT("/banners/:id", bannerHandler.UpdateBanner)
	router.DELETE("/banners/:id", bannerHandler.DeleteBanner)
	router.PATCH("/banners/:id", bannerHandler.PatchBanner)

	router.GET("/banners/:id/stats", statsHandler.GetBannerStatsByID)
	router.GET("/campaigns/:id/stats", statsHandler.GetCampaignStatsByID)
	router.GET("/campaigns/:id/stats/daily", statsHandler.GetDailyStats)

	router.POST("/banners/:id/impression", eventHandler.TrackImpression)
	router.POST("/banners/:id/click", eventHandler.TrackClick)

	router.GET("/display", adHandler.GetBannerForDisplay)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
