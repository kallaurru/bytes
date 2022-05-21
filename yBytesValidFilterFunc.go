package yBytes

/**
Пока будем придерживаться концепции групповых функций, о возможности формирования руной слова.
*/

// ValidIsClassicWord - классическое слово
func ValidIsClassicWord(code rune, isFirst bool) bool {
	// isFirst - для поддержки интерфейса
	return FIsClassicWord(code)
}

// ValidIsNumericWord - классическая группа цифр
func ValidIsNumericWord(code rune, isFirst bool) bool {
	// isFirst - для поддержки интерфейса
	return FIsNumericWord(code)
}

// ValidIsAdvWord - максимально расширенный функционал по охвату возможных символов из которых состоит слово
func ValidIsAdvWord(code rune, isFirst bool) bool {
	if isFirst {
		return isLetter(code) || isNumber(code) || isTagBeginSymbol(code)
	}
	return isLetter(code) || isNumber(code) || isWordSpecSymbols(code) || isTagBeginSymbol(code)
}

// FIsClassicWord - классическое слово. Только буквы
func FIsClassicWord(code rune) bool {
	return isLetter(code)
}

// FIsNumericWord - числа
func FIsNumericWord(code rune) bool {
	return isNumber(code)
}

// FIsFormattedNumber - форматированные суммы
func FIsFormattedNumber(code rune) bool {
	// 32 - пробел,
	return isNumber(code) || code == 32 || isCurrencyFormatSymbols(code)
}
