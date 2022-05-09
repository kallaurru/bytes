package yBytes

/**
Пока будем придерживаться концепции групповых функций, о возможности формирования руной слова.
*/

// GrIsClassicWord - классическое слово
func GrIsClassicWord(code rune, isFirst bool) bool {
	// isFirst - для поддержки интерфейса
	return isLetter(code)
}

// GrIsNumeric - классическая группа цифр
func GrIsNumeric(code rune, isFirst bool) bool {
	// isFirst - для поддержки интерфейса
	return isNumber(code)
}

// GrIsAdvancedWords - максимально расширенный функционал по охвату возможных символов из которых состоит слово
func GrIsAdvancedWords(code rune, isFirst bool) bool {
	if isFirst {
		return isLetter(code) || isNumber(code) || isTagBeginSymbol(code)
	}
	return isLetter(code) || isNumber(code) || isWordSpecSymbols(code) || isTagBeginSymbol(code)
}
