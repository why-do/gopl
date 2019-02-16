package unitconv

func CTOF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func IToM(i Inch) Meter {
	return Meter(i * 2.54 / 100)
}

func MToI(m Meter) Inch {
	return Inch(m * 100 * 0.394)
}

func PToK(p Pound) Kilogram {
	return Kilogram(p * 0.453592)
}

func KToP(k Kilogram) Pound {
	return Pound(k * 2.20462)
}