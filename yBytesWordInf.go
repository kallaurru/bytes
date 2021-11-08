package yBytes

/**
Пока будем придерживаться концепции групповых функций, о возможности формирования руной слова.
*/

// GrIsClassicWord - классическое слово
func GrIsClassicWord(code rune) bool {
	return isLetter(code)
}

// GrIsNumeric - классическая группа цифр
func GrIsNumeric(code rune) bool {
	return isNumber(code)
}

// GrIsAdvancedWords - максимально расширенный функционал по охвату возможных символов из которых состоит слово
func GrIsAdvancedWords(code rune, isFirst bool) bool {
	if isFirst {
		return isLetter(code) || isNumber(code) || isTagBeginSymbol(code)
	}
	return isLetter(code) || isNumber(code) || isWordSpecSymbols(code) || isTagBeginSymbol(code)
}

/**
Сейчас до конца не проработаны эти функции.
Замораживаем
*/

// ValidateWordClassic классическое слово, только буквы
func ValidateWordClassic(word *string, isQuickMode bool) bool {
	if isQuickMode {
		return bringToCorrectFormQuick(word, vswClassicLetterWord)
	}
	return bringToCorrectFormFull(word, vswClassicLetterWord)
}

// ValidateWordAdvanced расширенный набор символов, может найти слова через дефис, email
func ValidateWordAdvanced(word *string, isQuickMode bool) bool {
	if isQuickMode {
		return bringToCorrectFormQuick(word, vswAdvancedWord)
	}
	return bringToCorrectFormFull(word, vswAdvancedWord)
}

// ValidateClassicNumericWord классические числа без форматирования
func ValidateClassicNumericWord(word *string, isQuickMode bool) bool {
	if isQuickMode {
		return bringToCorrectFormQuick(word, vswClassicNumericWord)
	}
	return bringToCorrectFormFull(word, vswClassicNumericWord)
}

// ValidateFormatSum форматированная сумма с разделением разрядов и десятых
func ValidateFormatSum(word *string, isQuickMode bool) bool {
	if isQuickMode {
		return bringToCorrectFormQuick(word, vswNumericWithSymbols)
	}
	return bringToCorrectFormFull(word, vswNumericWithSymbols)
}

// ValidatePhoneFormat форматированные номера телефонов
func ValidatePhoneFormat(word *string, isQuickMode bool) bool {
	if isQuickMode {
		return bringToCorrectFormQuick(word, vswInternationalPhoneNumberFormat)
	}
	return bringToCorrectFormFull(word, vswInternationalPhoneNumberFormat)
}
