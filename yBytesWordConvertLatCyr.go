package yBytes

// convertCyrToLat конвертируем кириллический символ в латинский по данным визуализации
// по простому - буквы одинаково выглядят
// если 0 - значит нет возможности конвертировать
func convertCyrToLat(code rune) rune {
	if cyrToLatMap == nil {
		cyrToLatMap = makeMapLookTheSameSymbols(true)
	}
	return convertSymbolLookTheSame(code, cyrToLatMap)
}

// convertLatToCyr конвертируем латински символ в кириллический по данным визуализации
// по простому - буквы одинаково выглядят
func convertLatToCyr(code rune) rune {
	if latToCyrMap == nil {
		latToCyrMap = makeMapLookTheSameSymbols(false)
	}
	return convertSymbolLookTheSame(code, latToCyrMap)
}

// convertSymbolLookTheSame конвертируем символ по загруженной карте
func convertSymbolLookTheSame(code rune, codesMap *map[rune]rune) rune {
	value, ok := (*codesMap)[code]
	if !ok {
		return 0
	}
	return value
}

// makeMapLookTheSameSymbols toCyr - направление с латиницы в кириллицу, если false - наоборот
func makeMapLookTheSameSymbols(toCyr bool) *map[rune]rune {
	out := map[rune]rune{}
	// буква А
	if toCyr {
		out[65] = 1040
	} else {
		out[1040] = 65
	}
	// буква В
	if toCyr {
		out[66] = 1042
	} else {
		out[1042] = 66
	}
	// буква С
	if toCyr {
		out[67] = 1057
	} else {
		out[1057] = 67
	}
	// буква Е
	if toCyr {
		out[69] = 1045
	} else {
		out[1045] = 69
	}
	// буква H
	if toCyr {
		out[72] = 1053
	} else {
		out[1053] = 72
	}
	// буква К
	if toCyr {
		out[75] = 1050
	} else {
		out[1050] = 75
	}
	// буква М
	if toCyr {
		out[77] = 1052
	} else {
		out[1052] = 77
	}
	// буква О
	if toCyr {
		out[79] = 1054
	} else {
		out[1054] = 79
	}
	// буква Р
	if toCyr {
		out[80] = 1056
	} else {
		out[1056] = 80
	}
	// буква Т
	if toCyr {
		out[84] = 1058
	} else {
		out[1058] = 84
	}
	// буква Х
	if toCyr {
		out[88] = 1061
	} else {
		out[1061] = 88
	}
	// буква а
	if toCyr {
		out[97] = 1072
	} else {
		out[1072] = 97
	}
	// буква с
	if toCyr {
		out[99] = 1089
	} else {
		out[1089] = 99
	}
	// буква е
	if toCyr {
		out[101] = 1077
	} else {
		out[1077] = 101
	}
	// буква о
	if toCyr {
		out[111] = 1086
	} else {
		out[1086] = 111
	}
	// буква р
	if toCyr {
		out[112] = 1088
	} else {
		out[1088] = 112
	}
	// буква х
	if toCyr {
		out[120] = 1093
	} else {
		out[1093] = 120
	}

	return &out
}
