package handlers

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

// GastosRepository es lo que el handler necesita de donde sea que
// se guarden los gastos — no le importa si es Postgres, memoria, o un mock.
type GastosRepository interface {
	Listar(categoria string) ([]models.Gasto, error)
	ObtenerPorID(id uint) (models.Gasto, error)
	Crear(gasto *models.Gasto) error
	Guardar(gasto *models.Gasto) error
	Eliminar(gasto *models.Gasto) error
}
type GormGastosRepository struct {
	DB *gorm.DB
}

func NewGormGastosRepository(db *gorm.DB) *GormGastosRepository {
	return &GormGastosRepository{DB: db}
}

func (r *GormGastosRepository) Listar(categoria string) ([]models.Gasto, error) {
	var gastos []models.Gasto
	query := r.DB.Order("fecha desc")
	if categoria != "" {
		query = query.Where("categoria = ?", categoria)
	}
	err := query.Find(&gastos).Error
	return gastos, err
}

func (r *GormGastosRepository) ObtenerPorID(id uint) (models.Gasto, error) {
	var gasto models.Gasto
	err := r.DB.First(&gasto, id).Error
	return gasto, err
}

func (r *GormGastosRepository) Crear(gasto *models.Gasto) error {
	return r.DB.Create(gasto).Error
}

func (r *GormGastosRepository) Guardar(gasto *models.Gasto) error {
	return r.DB.Save(gasto).Error
}

func (r *GormGastosRepository) Eliminar(gasto *models.Gasto) error {
	return r.DB.Delete(gasto).Error
}
