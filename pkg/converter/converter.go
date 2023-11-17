package converter

func Temperature(inputUnit, outputUnit string, value float32) float32 {
	if inputUnit == "C" {
		if outputUnit == "C" {
			return value
		}

		if outputUnit == "F" {
			return (value * 9 / 5) + 32
		}
	}

	if inputUnit == "F" {
		if outputUnit == "C" {
			return (value - 32) * 5 / 9
		}

		if outputUnit == "F" {
			return value
		}
	}

	return 0

}
