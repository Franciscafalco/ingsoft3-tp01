package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"backend/internal/db"
	"backend/internal/handlers"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// Sentry: el DSN llega por el entorno. Sin SENTRY_DSN (tests, QA, compose local) queda apagado.
	// Release usa el commit que la imagen ya trae horneado (COMMIT_SHA): cada error dice qué commit lo produjo.
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Release:          os.Getenv("COMMIT_SHA"),
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Printf("sentry no se pudo iniciar: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	repo := handlers.NewGormGastosRepository(database)
	gastosHandler := handlers.NewGastosHandler(repo)
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	// /health acepta GET y HEAD: el plan gratuito de UptimeRobot pregunta con HEAD.
	router.Match([]string{http.MethodGet, http.MethodHead}, "/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "commit": os.Getenv("COMMIT_SHA")})
	})

	api := router.Group("/api")
	{
		api.GET("/gastos", gastosHandler.Listar)
		api.POST("/gastos", gastosHandler.Crear)
		api.PUT("/gastos/:id", gastosHandler.Actualizar)
		api.DELETE("/gastos/:id", gastosHandler.Eliminar)
		api.GET("/gastos/resumen", gastosHandler.Resumen)
	}

	router.Run(":8080")
}
