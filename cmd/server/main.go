package main

import (
	"log"
	"os"

	"OmniClima/internal/platform/middleware"
	"OmniClima/internal/platform/postgres"
	"OmniClima/internal/sensor"
	"OmniClima/internal/user"
	"OmniClima/internal/webhook"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo env não encontrado")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL não configurado")
	}

	db := postgres.Connection(dsn)
	_ = db

	db.AutoMigrate(
		&sensor.Sensor{},
		&sensor.SensorData{},
		&user.User{},
	)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	sensorRepo := sensor.NewRepository(db)
	sensorSvc := sensor.NewService(sensorRepo)
	sensorHdl := sensor.NewHandler(sensorSvc)
	webhookHdl := webhook.NewHandler(sensorSvc)

	public := r.Group("/")
	{
		public.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Bem-vindo à API OmniClima!"})
		})

		// Registro de rotas de módulos que são públicas
		// weatherHdl.RegisterPublicRoutes(public)
	}

	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware(db))
	{
		sensors := private.Group("/sensors")
		sensorHdl.RegisterRoutes(sensors)
		// userHdl.RegisterRoutes(private)
		// wearHdl.RegisterRoutes(private)
		// sportsHdl.RegisterRoutes(private)
		// leisureHdl.RegisterRoutes(private)
	}

	webhookR := r.Group("/webhook")
	webhookR.Use(middleware.AuthWebhook(db))
	{
		webhookHdl.RegisterRoutes(webhookR)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
