package yBytes

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
