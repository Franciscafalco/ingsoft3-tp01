package models

import (
	"errors"
	"time"
)

type Estado string

const (
	EstadoPendiente Estado = "pendiente"
	EstadoPagado    Estado = "pagado"
)

var CategoriasValidas = []string{
	"comida",
	"transporte",
	"vivienda",
	"entretenimiento",
	"salud",
	"otros",
}

var (
	ErrMontoInvalido      = errors.New("el monto debe ser mayor a 0")
	ErrFechaFutura        = errors.New("la fecha no puede ser futura")
	ErrCategoriaInvalida  = errors.New("categoría inválida")
	ErrTransicionInvalida = errors.New("un gasto pagado no puede volver a pendiente")
	ErrEliminarPagado     = errors.New("no se puede eliminar un gasto pagado")
)

type Gasto struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Monto       float64   `json:"monto"`
	Categoria   string    `json:"categoria"`
	Descripcion string    `json:"descripcion"`
	Fecha       time.Time `json:"fecha"`
	Estado      Estado    `gorm:"default:pendiente" json:"estado"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func categoriaValida(categoria string) bool {
	for _, c := range CategoriasValidas {
		if c == categoria {
			return true
		}
	}
	return false
}

// Validar aplica las reglas de negocio comunes a alta y edición de un gasto.
func Validar(monto float64, categoria string, fecha time.Time) error {
	if monto <= 0 {
		return ErrMontoInvalido
	}
	if !categoriaValida(categoria) {
		return ErrCategoriaInvalida
	}
	if fecha.After(time.Now()) {
		return ErrFechaFutura
	}
	return nil
}

// ValidarTransicion aplica la regla de transición de estado: pendiente -> pagado
// está permitido; pagado -> pendiente no.
func ValidarTransicion(actual, nuevo Estado) error {
	if actual == EstadoPagado && nuevo == EstadoPendiente {
		return ErrTransicionInvalida
	}
	return nil
}

// ValidarEliminacion aplica la restricción: no se puede eliminar un gasto pagado.
func ValidarEliminacion(actual Estado) error {
	if actual == EstadoPagado {
		return ErrEliminarPagado
	}
	return nil
}

// TODO: endpoint de gasto
