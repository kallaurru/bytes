package yBytes

// типы функций которые проверяют символы из которых может состоять слово
// каждая функция содержит определенный список допустимых символов
// именование начинаем с префикса vsw (valid symbol word)
type funcValidateWordSymbols func(symbol rune, isBorder bool) bool

/** Функции валидации символов слова */

func vswClassicLetterWord(symbol rune, isBorder bool) bool {
	// для общего интерфейса
	isBorder = true
	return isLetter(symbol) && isBorder
}

func vswClassicNumericWord(symbol rune, isBorder bool) bool {
	// для общего интерфейса
	isBorder = true
	return isNumber(symbol) && isBorder
}

func vswAdvancedWord(symbol rune, isBorder bool) bool {
	if isBorder {
		return isLetter(symbol) || isNumber(symbol) || (symbol == 37)
	}
	return isLetter(symbol) || isWordSpecSymbols(symbol) || isNumber(symbol)
}

func vswNumericWithSymbols(symbol rune, isBorder bool) bool {
	if isBorder {
		// пропускаем знак доллара и процента
		return (symbol == 36) || (symbol == 37) || isNumber(symbol)
	}
	return isNumber(symbol) || isCurrencyFormatSymbols(symbol)
}

func vswInternationalPhoneNumberFormat(symbol rune, isBorder bool) bool {
	if isBorder {
		return isInternationalPhoneNumberFormatSymbolBorder(symbol)
	}
	return isInternationalPhoneNumberFormatSymbolBody(symbol)
}

// общий функционал определения валидности входящего в слово символа
func isValidSymbolWord(symbol rune, isBorder bool, validateFunc funcValidateWordSymbols) bool {
	return validateFunc(symbol, isBorder)
}

func validateWordLen(word string) bool {
	lenWord := len(word)
	if lenWord <= 0 || lenWord > 1024 {
		return false
	}
	return true
}

// cutHeaderWord обрезаем лишние символы с начала слова.
// Возвращаем true если ничего не изменилось
func cutHeaderWord(runeView *[]rune, isValidFirst, isValidLast bool, validatingFunc funcValidateWordSymbols) bool {
	// вариант когда первый символ слова валиден. Ничего не делаем
	if isValidFirst {
		return true
	}

	lenRuneView := len(*runeView)
	firsPosNewWord := 0

	if lenRuneView < 3 {
		if !isValidLast {
			*runeView = []rune{}
			return false
		}
	}

	lastPositionCycle := lenRuneView - 1
	if !isValidLast {
		lastPositionCycle--
	}

	for i := 1; i <= lastPositionCycle; i++ {
		if validatingFunc((*runeView)[i], true) {
			firsPosNewWord = i
			break
		}
	}
	if firsPosNewWord == 0 {
		// значит все символы не валидны
		*runeView = []rune{}
		return false
	}
	// скорректированный массив
	*runeView = (*runeView)[firsPosNewWord:]
	return false
}

// полностью не валидное слово или короткая строка обрабатываются в обрезке головы слова
func cutTailWord(runeView *[]rune, isValidLast bool, validatingFunc funcValidateWordSymbols) bool {
	// вариант когда последний символ слова валиден. Ничего не делаем
	if isValidLast {
		return true
	}

	lenRuneView := len(*runeView)
	// если дошло до этой позиции последний символ не валиден
	lastPosNewWord := lenRuneView - 1

	for i := lastPosNewWord; i > 0; i-- {
		idx := lastPosNewWord - 1
		if validatingFunc((*runeView)[idx], true) {
			lastPosNewWord = i
			break
		}
	}
	// скорректированный массив
	*runeView = (*runeView)[:lastPosNewWord]
	return false
}

// проходим по всему слову и удаляем не валидные символы
func bringToCorrectFormFull(word *string, validatingFunc funcValidateWordSymbols) bool {
	if !validateWordLen(*word) {
		return false
	}
	runeView := []rune(*word)
	lenRuneView := len(runeView)
	runeViewNew := make([]rune, 0, lenRuneView)
	for idx, code := range runeView {
		isBorder := (idx == 0) || (idx == lenRuneView-1)
		if validatingFunc(code, isBorder) {
			runeViewNew = append(runeViewNew, code)
		}
	}
	*word = string(runeViewNew)

	return len(runeViewNew) != 0
}

// общий функционал интерфейсных функций валидации слова. Быстрый алгоритм, отрезаются только не подходящие символы
// только на границах слова.
// return false оказалось пустым в результате преобразования
func bringToCorrectFormQuick(word *string, validatingFunc funcValidateWordSymbols) bool {
	if !validateWordLen(*word) {
		return false
	}
	runeView := []rune(*word)
	lenRuneView := len(runeView)
	// проверим на валидность все символы слова. Если первый и последний символы окажутся верными значит слово правильное
	isValidFirst := isValidSymbolWord(runeView[0], true, validatingFunc)
	isValidLast := isValidSymbolWord(runeView[lenRuneView-1], true, validatingFunc)

	if !isValidFirst || !isValidLast {
		cutHeaderWord(&runeView, isValidFirst, isValidLast, validatingFunc)
		if len(runeView) == 0 {
			// не валидное слово с отсчетом сначала
			*word = ""
			return false
		}
		cutTailWord(&runeView, isValidLast, validatingFunc)
		if len(runeView) == 0 {
			*word = ""
			return false
		}
		*word = string(runeView)
	}
	return true
}
