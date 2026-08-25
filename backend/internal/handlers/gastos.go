package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GastosHandler struct {
	DB *gorm.DB
}

func NewGastosHandler(db *gorm.DB) *GastosHandler {
	return &GastosHandler{DB: db}
}

type gastoRequest struct {
	Monto       float64   `json:"monto"`
	Categoria   string    `json:"categoria"`
	Descripcion string    `json:"descripcion"`
	Fecha       time.Time `json:"fecha"`
}

type gastoUpdateRequest struct {
	gastoRequest
	Estado models.Estado `json:"estado"`
}

func (h *GastosHandler) Listar(c *gin.Context) {
	var gastos []models.Gasto
	query := h.DB.Order("fecha desc")
	if categoria := c.Query("categoria"); categoria != "" {
		query = query.Where("categoria = ?", categoria)
	}
	if err := query.Find(&gastos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gastos)
}

func (h *GastosHandler) Resumen(c *gin.Context) {
	var gastos []models.Gasto
	if err := h.DB.Find(&gastos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.CalcularResumen(gastos))
}

func (h *GastosHandler) Crear(c *gin.Context) {
	var req gastoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.Validar(req.Monto, req.Categoria, req.Fecha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gasto := models.Gasto{
		Monto:       req.Monto,
		Categoria:   req.Categoria,
		Descripcion: req.Descripcion,
		Fecha:       req.Fecha,
		Estado:      models.EstadoPendiente,
	}
	if err := h.DB.Create(&gasto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gasto)
}

func (h *GastosHandler) Actualizar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var gasto models.Gasto
	if err := h.DB.First(&gasto, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "gasto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req gastoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.Validar(req.Monto, req.Categoria, req.Fecha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.ValidarTransicion(gasto.Estado, req.Estado); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	gasto.Monto = req.Monto
	gasto.Categoria = req.Categoria
	gasto.Descripcion = req.Descripcion
	gasto.Fecha = req.Fecha
	gasto.Estado = req.Estado

	if err := h.DB.Save(&gasto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gasto)
}

func (h *GastosHandler) Eliminar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var gasto models.Gasto
	if err := h.DB.First(&gasto, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "gasto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidarEliminacion(gasto.Estado); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Delete(&gasto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
