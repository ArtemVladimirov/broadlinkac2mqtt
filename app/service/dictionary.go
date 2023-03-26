package service

var (
	verticalFixationStatuses = map[int]string{
		0b00000001: "top",
		0b00000010: "middle1",
		0b00000011: "middle2",
		0b00000100: "middle3",
		0b00000101: "bottom",
		0b00000110: "swing",
		0b00000111: "auto",
	}

	verticalFixationStatusesInvert = map[string]int{
		"top":     0b00000001,
		"middle1": 0b00000010,
		"middle2": 0b00000011,
		"middle3": 0b00000100,
		"bottom":  0b00000101,
		"swing":   0b00000110,
		"auto":    0b00000111,
	}

	//horizontalFixationStatuses = map[int]string{
	//	2: "LEFT_FIX",
	//	1: "LEFT_FLAP",
	//	7: "LEFT_RIGHT_FIX",
	//	0: "LEFT_RIGHT_FLAP",
	//	6: "RIGHT_FIX",
	//	5: "RIGHT_FLAP",
	//	0: "ON",
	//	1: "OFF",
	//}

	fanStatuses = map[int]string{
		0b00000011: "low",
		0b00000010: "medium",
		0b00000001: "high",
		0b00000101: "auto",
		0b00000000: "none",
	}

	fanStatusesInvert = map[string]int{
		"low":    0b00000011,
		"medium": 0b00000010,
		"high":   0b00000001,
		"auto":   0b00000101,
		"none":   0b00000000,
	}

	modeStatuses = map[int]string{
		0b00000001: "cool",
		0b00000010: "dry",
		0b00000100: "heat",
		0b00000000: "auto",
		0b00000110: "fan_only",
	}

	modeStatusesInvert = map[string]int{
		"cool":     0b00000001,
		"dry":      0b00000010,
		"heat":     0b00000100,
		"auto":     0b00000000,
		"fan_only": 0b00000110,
	}

	onOffStatusesInvert = map[string]int{
		"ON":  1,
		"OFF": 0,
	}
)
