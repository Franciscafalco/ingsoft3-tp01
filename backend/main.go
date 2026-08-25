package main

import (
	"log"
	"net/http"

	"backend/internal/db"
	"backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	gastosHandler := handlers.NewGastosHandler(database)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
