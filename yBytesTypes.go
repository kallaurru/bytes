package yBytes

type YByte = byte

type GrWordFuncFistControl func(code rune, isFirst bool) bool
type GrWordFunc func(code rune) bool

// EncodeInformation информация по декодированному слову
// все свойства которые начинаются с rule, являются линейками позиционирования нужных символов
// !!! Все позиционирование в линейках начинается с ноля
// по умолчанию считается что нормальное слово  !!!! не будет больше 64 символов
type EncodeInformation struct {
	// храним список полученных ошибок в результате конвертации
	errList []error
	// представление кодированного слова в виде массива универсальной кодировки
	original []YByte
	// представление конвертированного слова
	converted []YByte
	// представление в виде рун. Кэш слова для быстрого перевода в строку. Хранится в прописных символах
	runes []rune
	// позиции цифр в слове
	rulePosNumbers uint64
	// позиции допустимых не буквенно-цифровых символов
	rulePosSymbols uint64
	// позиции прописных букв в слове. Для работы с аббревиатурами
	rulePosCapitalSymbols uint64
	// позиции букв кириллицы в слове
	rulePosSymbolsCyr uint64
	// позиции букв латиницы в слове
	rulePosSymbolsLat uint64
	// false - существующие в слове символы оппозитной кодировки не могут быть конвертированы
	isNotConverting bool
	// направление конвертирования -0 не потребовалось, 1 - из латиницы в кириллицу, 2 - из кириллицы в латиницу
	directionConverting uint8
}

// показываем в виде универсальных байтов
func (ei *EncodeInformation) viewEncoding() []YByte {
	return ei.original
}

// показываем как было в оригинале
func (ei *EncodeInformation) viewOriginal() string {
	cacheRune := make([]rune, 0, len(ei.original))
	for pos, yByte := range ei.original {
		mask := 1 << pos
		isCapital := ei.rulePosCapitalSymbols&mask > 0
		runeSymbol := decodeSymbol(yByte, isCapital)
		cacheRune[pos] = runeSymbol
	}
	return string(cacheRune)
}

// показываем прописными буквами приведенными к одному типу символов если это возможно
func (ei *EncodeInformation) viewClassic() string {
	if len(ei.runes) == 0 {
		ei.makeCache()
	}
	return string(ei.runes)
}

// Формируем кэш рун. Все буквы строчные
func (ei *EncodeInformation) makeCache() {
	ei.runes = make([]rune, 0, len(ei.original))
	symbolList := ei.original

	if ei.directionConverting > 0 && ei.isNotConverting == false {
		symbolList = ei.converted
	}

	for idx, yByte := range symbolList {
		runeElement := decodeSymbol(yByte, false)
		ei.runes[idx] = runeElement
	}
}

// Пробуем конвертировать символы
func (ei *EncodeInformation) convert() {
	var (
		convertingRule uint64
	)
	countCyr := GetCountOnBytes64(ei.rulePosSymbolsCyr)
	countLat := GetCountOnBytes64(ei.rulePosSymbolsLat)
	if countLat == 0 || countCyr == 0 {
		// не нуждаемся в конвертации
		return
	}
	if countCyr >= countLat {
		// конвертируем латиницу в кириллицу
		ei.directionConverting = 1
		convertingRule = ei.rulePosSymbolsLat
	} else {
		ei.directionConverting = 2
		convertingRule = ei.rulePosSymbolsCyr
	}
	ei.converted = ei.original
	for pos, yByte := range ei.original {
		mask := 1 << pos
		hasSymbol := convertingRule&mask > 0
		isCapital := ei.rulePosCapitalSymbols&mask > 0
		if hasSymbol {
			runeSymbol := decodeSymbol(yByte, isCapital)
			if ei.directionConverting == 1 {
				runeSymbol = convertLatToCyr(runeSymbol)
			} else if ei.directionConverting == 2 {
				runeSymbol = convertCyrToLat(runeSymbol)
			}
			if runeSymbol == 0 {
				// конвертация не удалась
				ei.isNotConverting = true
			} else {
				ei.converted[pos] = encodeSymbol(runeSymbol)
			}
		}
	}

}
