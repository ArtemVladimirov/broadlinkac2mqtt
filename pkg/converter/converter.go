package converter

func Temperature(inputUnit, outputUnit string, value float32) float32 {
	if inputUnit == "C" {
		if outputUnit == "C" {
			return value
		}

		if outputUnit == "F" {
			return float32(int((value * 9 / 5) + 32))
		}
	}

	if inputUnit == "F" {
		if outputUnit == "C" {
			return float32(int((value-32)*5/9*10)) / 10
		}

		if outputUnit == "F" {
			return value
		}
	}

	return 0
}
