package service

import (
	"reflect"
	"testing"

	"backend/internal/models"
)

func TestCalcularResumen(t *testing.T) {
	casos := []struct {
		nombre   string
		gastos   []models.Gasto
		esperado Resumen
	}{
		{
			nombre:   "sin gastos, el total es cero",
			gastos:   []models.Gasto{},
			esperado: Resumen{TotalGeneral: 0, PorCategoria: map[string]float64{}},
		},
		{
			nombre: "un gasto, el total es ese monto",
			gastos: []models.Gasto{
				{Monto: 100, Categoria: "comida"},
			},
			esperado: Resumen{
				TotalGeneral: 100,
				PorCategoria: map[string]float64{"comida": 100},
			},
		},
		{
			nombre: "varios gastos en distintas categorías se suman por separado",
			gastos: []models.Gasto{
				{Monto: 100, Categoria: "comida"},
				{Monto: 50, Categoria: "transporte"},
				{Monto: 30, Categoria: "comida"},
			},
			esperado: Resumen{
				TotalGeneral: 180,
				PorCategoria: map[string]float64{"comida": 130, "transporte": 50},
			},
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Act
			resultado := CalcularResumen(c.gastos)

			// Assert
			if !reflect.DeepEqual(resultado, c.esperado) {
				t.Errorf("CalcularResumen(...) = %+v, esperaba %+v", resultado, c.esperado)
			}
		})
	}
}
