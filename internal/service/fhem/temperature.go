package fhem

import (
	"fmt"
)

var (
	// ScaleCelsius defines temperature measured in °C
	ScaleCelsius = scale("C")
	// ScaleFahrenheit defines temperature measured in °F
	ScaleFahrenheit = scale("F")
	// ScaleKelvin defines temperature measured in °K
	ScaleKelvin = scale("K")
)

type scale string

// Temperature stores a temperature together with the scale used
type Temperature struct {
	// The temperature in the given unit
	Value float64
	Scale scale
}

func (t *Temperature) String() string {
	return fmt.Sprintf("%.2f°%s", t.Value, t.Scale)
}
