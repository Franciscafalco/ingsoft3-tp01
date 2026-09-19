package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestGastosHandler_Eliminar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("gasto pendiente: se llama Eliminar en el repositorio", func(t *testing.T) {
		// Arrange
		repo := new(RepositorioMock)
		gastoPendiente := models.Gasto{ID: 1, Estado: models.EstadoPendiente}
		repo.On("ObtenerPorID", uint(1)).Return(gastoPendiente, nil)
		repo.On("Eliminar", mock.Anything).Return(nil)

		handler := NewGastosHandler(repo)
		router := gin.New()
		router.DELETE("/gastos/:id", handler.Eliminar)

		req := httptest.NewRequest(http.MethodDelete, "/gastos/1", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert: no miramos un valor devuelto — miramos la INTERACCIÓN
		repo.AssertCalled(t, "Eliminar", mock.Anything)
		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, esperaba %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("gasto pagado: NUNCA se llama Eliminar", func(t *testing.T) {
		// Arrange
		repo := new(RepositorioMock)
		gastoPagado := models.Gasto{ID: 2, Estado: models.EstadoPagado}
		repo.On("ObtenerPorID", uint(2)).Return(gastoPagado, nil)

		handler := NewGastosHandler(repo)
		router := gin.New()
		router.DELETE("/gastos/:id", handler.Eliminar)

		req := httptest.NewRequest(http.MethodDelete, "/gastos/2", nil)
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		repo.AssertNotCalled(t, "Eliminar", mock.Anything)
		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, esperaba %d", w.Code, http.StatusConflict)
		}
	})
}

func TestGastosHandler_Actualizar_RechazaTransicionInvalida(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	repo := new(RepositorioMock)
	gastoPagado := models.Gasto{ID: 3, Monto: 100, Categoria: "comida", Fecha: time.Now(), Estado: models.EstadoPagado}
	repo.On("ObtenerPorID", uint(3)).Return(gastoPagado, nil)

	handler := NewGastosHandler(repo)
	router := gin.New()
	router.PUT("/gastos/:id", handler.Actualizar)

	body := `{"monto":100,"categoria":"comida","fecha":"2026-01-01T00:00:00Z","estado":"pendiente"}`
	req := httptest.NewRequest(http.MethodPut, "/gastos/3", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	repo.AssertNotCalled(t, "Guardar", mock.Anything)
	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, esperaba %d", w.Code, http.StatusConflict)
	}
}
func TestGastosHandler_Crear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("gasto válido: se crea y se llama al repositorio", func(t *testing.T) {
		// Arrange
		repo := new(RepositorioMock)
		repo.On("Crear", mock.Anything).Return(nil)

		handler := NewGastosHandler(repo)
		router := gin.New()
		router.POST("/gastos", handler.Crear)

		body := `{"monto":100,"categoria":"comida","descripcion":"Super","fecha":"2020-01-01T00:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/gastos", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		repo.AssertCalled(t, "Crear", mock.Anything)
		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, esperaba %d", w.Code, http.StatusCreated)
		}
	})

	t.Run("monto inválido: se rechaza sin tocar el repositorio", func(t *testing.T) {
		// Arrange
		repo := new(RepositorioMock)

		handler := NewGastosHandler(repo)
		router := gin.New()
		router.POST("/gastos", handler.Crear)

		body := `{"monto":0,"categoria":"comida","descripcion":"Super","fecha":"2020-01-01T00:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/gastos", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		router.ServeHTTP(w, req)

		// Assert
		repo.AssertNotCalled(t, "Crear", mock.Anything)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, esperaba %d", w.Code, http.StatusBadRequest)
		}
	})
}
func TestGastosHandler_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	repo := new(RepositorioMock)
	gastos := []models.Gasto{
		{ID: 1, Monto: 100, Categoria: "comida"},
		{ID: 2, Monto: 50, Categoria: "transporte"},
	}
	repo.On("Listar", "").Return(gastos, nil)

	handler := NewGastosHandler(repo)
	router := gin.New()
	router.GET("/gastos", handler.Listar)

	req := httptest.NewRequest(http.MethodGet, "/gastos", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	repo.AssertCalled(t, "Listar", "")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, esperaba %d", w.Code, http.StatusOK)
	}
}
