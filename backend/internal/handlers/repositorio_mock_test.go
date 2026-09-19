package handlers

import (
	"backend/internal/models"

	"github.com/stretchr/testify/mock"
)

// RepositorioMock implementa GastosRepository, pero en vez de tocar una base
// de verdad, registra qué le pidieron y devuelve lo que el test le programó.
type RepositorioMock struct {
	mock.Mock
}

func (m *RepositorioMock) Listar(categoria string) ([]models.Gasto, error) {
	args := m.Called(categoria)
	return args.Get(0).([]models.Gasto), args.Error(1)
}

func (m *RepositorioMock) ObtenerPorID(id uint) (models.Gasto, error) {
	args := m.Called(id)
	return args.Get(0).(models.Gasto), args.Error(1)
}

func (m *RepositorioMock) Crear(gasto *models.Gasto) error {
	args := m.Called(gasto)
	return args.Error(0)
}

func (m *RepositorioMock) Guardar(gasto *models.Gasto) error {
	args := m.Called(gasto)
	return args.Error(0)
}

func (m *RepositorioMock) Eliminar(gasto *models.Gasto) error {
	args := m.Called(gasto)
	return args.Error(0)
}
