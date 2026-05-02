package main

import (
	"log"
	"os"

	"bouquet-app/internal/database"
	"bouquet-app/internal/handlers"
	"bouquet-app/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "bouquet-app/docs" // generated swagger docs
)

// @title           Bouquet App API
// @version         2.0
// @description     API untuk aplikasi pemesanan bouquet bunga dengan AI design generator
// @termsOfService  http://swagger.io/terms/

// @contact.name   Bouquet App Support
// @contact.email  support@bouquet-app.id

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Admin-Key

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Init PostgreSQL + auto-migrate + seed
	database.Connect()

	r := gin.Default()
	r.Use(middleware.CORS())

	// ── Swagger UI ──────────────────────────────────────────
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Public API ──────────────────────────────────────────
	api := r.Group("/api/v1")
	{
		// Bouquet types
		api.GET("/bouquet-types", handlers.GetBouquetTypes)

		// Flowers (dari DB)
		api.GET("/flowers", handlers.GetFlowers)

		// Katalog pre-made
		api.GET("/catalog", handlers.GetCatalog)

		// AI Agents
		api.POST("/agent/verify-selection", handlers.AgentVerifySelection)
		api.POST("/agent/generate-bouquet", handlers.AgentGenerateBouquet)
		api.GET("/agent/generate-status", handlers.GetGenerateStatus)

		// Orders
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders/:id", handlers.GetOrder)
		api.GET("/orders/:id/tracking", handlers.GetTracking)

		// Payment
		api.POST("/payment/token", handlers.CreatePaymentToken)
		api.POST("/payment/notification", handlers.PaymentNotification)
	}

	// ── Admin API (protected) ────────────────────────────────
	admin := r.Group("/api/v1/admin")
	admin.Use(handlers.AdminMiddleware())
	{
		// Dashboard stats
		admin.GET("/stats", handlers.AdminGetStats)

		// Order management
		admin.GET("/orders", handlers.AdminGetOrders)
		admin.PUT("/orders/:id", handlers.AdminUpdateOrder)

		// Flower management
		admin.GET("/flowers", handlers.AdminGetFlowers)
		admin.PUT("/flowers/:id", handlers.AdminUpdateFlower)

		// Catalog management
		admin.GET("/catalog", handlers.AdminGetCatalog)
		admin.POST("/catalog", handlers.AdminCreateCatalog)
		admin.PUT("/catalog/:id", handlers.AdminUpdateCatalog)
		admin.DELETE("/catalog/:id", handlers.AdminDeleteCatalog)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s", port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
