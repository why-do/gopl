package unitconv

import "fmt"

type Celsius float64
type Fahrenheit float64

func (c Celsius) String() string {
	return fmt.Sprintf("%g℃ ", c)
}

func (f Fahrenheit) String() string {
	return fmt.Sprintf("%g℉ ", f)
}

type Inch float64
type Meter float64

func (i Inch) String() string {
	return fmt.Sprintf("%gin", i)
}

func (m Meter) String() string {
	return fmt.Sprintf("%gm", m)
}

type Pound float64
type Kilogram float64

func (p Pound) String() string {
	return fmt.Sprintf("%glb", p)
}

func (k Kilogram) String() string {
	return fmt.Sprintf("%gkg", k)
}