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
	Repo GastosRepository
}

func NewGastosHandler(repo GastosRepository) *GastosHandler {
	return &GastosHandler{Repo: repo}
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
	gastos, err := h.Repo.Listar(c.Query("categoria"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gastos)
}

func (h *GastosHandler) Resumen(c *gin.Context) {
	gastos, err := h.Repo.Listar("")
	if err != nil {
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
	if err := h.Repo.Crear(&gasto); err != nil {
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

	gasto, err := h.Repo.ObtenerPorID(uint(id))
	if err != nil {
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

	if err := h.Repo.Guardar(&gasto); err != nil {
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

	gasto, err := h.Repo.ObtenerPorID(uint(id))
	if err != nil {
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

	if err := h.Repo.Eliminar(&gasto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
