package yBytes

func isLetter(code rune) bool {
	return isLatinLetter(code) || isCyrillicLetter(code)
}

// Символы, которые могут включать слова
func isWordSpecSymbols(code rune) bool {
	// ` @ - _
	if code == 45 || code == 64 || code == 95 || code == 96 {
		return true
	}
	// % &
	if code > 36 && code < 39 {
		return true
	}
	return false
}

// Символы начала и конца международного формата номера телефона
func isInternationalPhoneNumberFormatSymbolBorder(code rune) bool {
	// +
	return code == 43 || isNumber(code)
}

// Символы международного формата номера телефона
func isInternationalPhoneNumberFormatSymbolBody(code rune) bool {
	return isNumber(code) || code == 32 || code == 40 || code == 41 || code == 45
}

// Форматированные суммы
func isCurrencyFormatSymbols(code rune) bool {
	// , .
	return code == 44 || code == 46
}

// Форматированные суммы
func isCurrencyFormatSymbolsWithSpace(code rune) bool {
	// пробел , .
	return code == 32 || code == 44 || code == 46
}

func isTagBeginSymbol(code rune) bool {
	// #
	return code == 35
}

func isNumber(code rune) bool {
	if code >= 48 && code <= 57 {
		return true
	}
	return false
}

func isLatinLetter(code rune) bool {
	if code >= 65 && code <= 90 {
		return true
	}

	if code >= 97 && code <= 122 {
		return true
	}

	return false
}

func isCyrillicLetter(code rune) bool {
	if code == 1025 || code == 1105 {
		return true
	}
	if code >= 1040 && code <= 1103 {
		return true
	}

	return false
}

func isUpperCyrillicLetter(code rune) bool {
	if code == 1025 {
		return true
	}
	if code >= 1040 && code <= 1071 {
		return true
	}

	return false
}

func isLowerCyrillicLetter(code rune) bool {
	if code == 1105 {
		return true
	}
	if code >= 1072 && code <= 1103 {
		return true
	}

	return false
}

func isUpperLatinLetter(code rune) bool {
	if code >= 65 && code <= 90 {
		return true
	}
	return false
}

func isLowerLatinLetter(code rune) bool {
	if code >= 97 && code <= 122 {
		return true
	}

	return false
}

func isCapitalSymbol(code rune) bool {
	return isUpperLatinLetter(code) || isUpperCyrillicLetter(code)
}

func isSymbol(code rune) bool {
	notSymbol := isNumber(code) || isLatinLetter(code) || isCyrillicLetter(code)
	return !notSymbol
}
