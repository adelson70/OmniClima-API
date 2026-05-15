package main

import (
	"log"
	"os"

	"OmniClima/internal/platform/postgres"

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

	// db.AutoMigrate()

	r := gin.Default()

	public := r.Group("/")
	{
		public.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Bem-vindo à API OmniClima!"})
		})

		// Registro de rotas de módulos que são públicas
		// weatherHdl.RegisterPublicRoutes(public)
	}

	private := r.Group("/api")
	_ = private
	// private.Use(authMiddleware())
	{
		// Aqui você registra os módulos que exigem login do usuário
		// userHdl.RegisterRoutes(private)
		// wearHdl.RegisterRoutes(private)
		// sportsHdl.RegisterRoutes(private)
		// leisureHdl.RegisterRoutes(private)
	}

	webhook := r.Group("/webhook")
	_ = webhook
	// webhook.Use(authWebhookMiddleware())
	{
		// Registro do módulo de sensores
		// sensorHdl.RegisterRoutes(webhook)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
