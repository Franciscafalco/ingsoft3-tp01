package service

import "backend/internal/models"

type Resumen struct {
	TotalGeneral  float64            `json:"totalGeneral"`
	PorCategoria  map[string]float64 `json:"porCategoria"`
}

// CalcularResumen suma el total general y el total por categoría de una
// lista de gastos. Es una función pura (sin acceso a base de datos) para
// que se pueda testear con casos típicos, cero gastos y una sola categoría.
func CalcularResumen(gastos []models.Gasto) Resumen {
	resumen := Resumen{PorCategoria: map[string]float64{}}
	for _, g := range gastos {
		resumen.TotalGeneral += g.Monto
		resumen.PorCategoria[g.Categoria] += g.Monto
	}
	return resumen
}
