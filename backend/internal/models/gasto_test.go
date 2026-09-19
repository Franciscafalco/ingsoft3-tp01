package models

import (
	"testing"
	"time"
)

func TestValidar_Monto(t *testing.T) {
	casos := []struct {
		nombre      string
		monto       float64
		quiereError bool
	}{
		{"un monto positivo se acepta", 100, false},
		{"un monto en cero se rechaza", 0, true},
		{"un monto negativo se rechaza", -50, true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange: el resto de los datos, válidos, para aislar SOLO la regla del monto
			categoria := "comida"
			fecha := time.Now()

			// Act
			err := Validar(c.monto, categoria, fecha)

			// Assert
			huboError := err != nil
			if huboError != c.quiereError {
				t.Errorf("Validar(monto=%v) => error=%v, esperaba error=%v", c.monto, err, c.quiereError)
			}
		})
	}
}

func TestValidar_Fecha(t *testing.T) {
	casos := []struct {
		nombre      string
		fecha       time.Time
		quiereError bool
	}{
		{"una fecha pasada se acepta", time.Now().Add(-24 * time.Hour), false}, // dentro de 24 horas
		{"la fecha de hoy se acepta", time.Now(), false},
		{"una fecha futura se rechaza", time.Now().Add(24 * time.Hour), true}, // fuera de 24 horas
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange: monto y categoría válidos, para aislar SOLO la regla de la fecha
			monto := 100.0
			categoria := "comida"

			// Act
			err := Validar(monto, categoria, c.fecha)

			// Assert
			huboError := err != nil
			if huboError != c.quiereError {
				t.Errorf("Validar(fecha=%v) => error=%v, esperaba error=%v", c.fecha, err, c.quiereError)
			}
		})
	}
}

func TestValidar_Categoria(t *testing.T) {
	casos := []struct {
		nombre      string
		categoria   string
		quiereError bool
	}{
		{"una categoría de la lista se acepta", "comida", false},
		{"una categoría fuera de la lista se rechaza", "inexistente", true},
		{"una categoría vacía se rechaza", "", true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			monto := 100.0
			fecha := time.Now()

			// Act
			err := Validar(monto, c.categoria, fecha)

			// Assert
			huboError := err != nil
			if huboError != c.quiereError {
				t.Errorf("Validar(categoria=%q) => error=%v, esperaba error=%v", c.categoria, err, c.quiereError)
			}
		})
	}
}

func TestValidarTransicion(t *testing.T) { //sin Arrange porque los datos de entrada van directo en la tabla, no hay qué preparar
	casos := []struct {
		nombre      string
		actual      Estado
		nuevo       Estado
		quiereError bool
	}{
		{"pendiente a pagado se permite", EstadoPendiente, EstadoPagado, false},
		{"pagado a pendiente se rechaza", EstadoPagado, EstadoPendiente, true},
		{"pagado a pagado no cambia nada, se permite", EstadoPagado, EstadoPagado, false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Act
			err := ValidarTransicion(c.actual, c.nuevo)

			// Assert
			huboError := err != nil
			if huboError != c.quiereError {
				t.Errorf("ValidarTransicion(%v -> %v) => error=%v, esperaba error=%v", c.actual, c.nuevo, err, c.quiereError)
			}
		})
	}
}

func TestValidarEliminacion(t *testing.T) {
	casos := []struct {
		nombre      string
		actual      Estado
		quiereError bool
	}{
		{"un gasto pendiente se puede eliminar", EstadoPendiente, false},
		{"un gasto pagado no se puede eliminar", EstadoPagado, true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Act
			err := ValidarEliminacion(c.actual)

			// Assert
			huboError := err != nil
			if huboError != c.quiereError {
				t.Errorf("ValidarEliminacion(%v) => error=%v, esperaba error=%v", c.actual, err, c.quiereError)
			}
		})
	}
}
