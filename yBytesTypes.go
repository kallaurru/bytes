package yBytes

type YByte = byte

// GrWordFuncFistControl функции контроля сборки слова. С контролем начального разрешенного символа
type GrWordFuncFistControl func(code rune, isFirst bool) bool

// GrWordFunc функции контроля сборки слова. Без контроля первого слова
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

// ViewEncoding показываем в виде универсальных байтов
func (ei *EncodeInformation) ViewEncoding() []YByte {
	ei.prepare()
	out := ei.original
	if len(ei.converted) > 0 && ei.isNotConverting == false {
		out = ei.converted
	}

	return out
}

// ViewEncodingOriginal показываем вариант слова без конвертирования
func (ei *EncodeInformation) ViewEncodingOriginal() []YByte {
	return ei.original
}

// ViewOriginal показываем как было в оригинале
func (ei *EncodeInformation) ViewOriginal() string {
	ei.prepare()
	cacheRune := make([]rune, len(ei.original), len(ei.original))
	for pos, yByte := range ei.original {
		mask := uint64(1 << pos)
		isCapital := ei.rulePosCapitalSymbols&mask > 0
		runeSymbol := decodeSymbol(yByte, isCapital)
		cacheRune[pos] = runeSymbol
	}
	return string(cacheRune)
}

// ViewClassic показываем прописными буквами приведенными к одному типу символов если это возможно
func (ei *EncodeInformation) ViewClassic() string {
	ei.prepare()
	if len(ei.runes) == 0 {
		ei.makeCache()
	}
	return string(ei.runes)
}

// Формируем кэш рун. Все буквы строчные
func (ei *EncodeInformation) makeCache() {
	ei.prepare()
	ei.runes = make([]rune, len(ei.original), len(ei.original))
	symbolList := ei.original

	if ei.directionConverting > 0 && ei.isNotConverting == false {
		symbolList = ei.converted
	}

	for idx, yByte := range symbolList {
		runeElement := decodeSymbol(yByte, false)
		ei.runes[idx] = runeElement
	}
}

// подготовка под возможную конвертацию
func (ei *EncodeInformation) prepare() {
	if ei.directionConverting == 0 && (ei.rulePosSymbolsCyr == 0 || ei.rulePosSymbolsLat == 0) {
		// подготовка под конвертацию не нужна
		return
	}
	if ei.converted != nil {
		// уже была конвертация
		return
	}
	countCyr := GetCountOnBytes64(ei.rulePosSymbolsCyr)
	countLat := GetCountOnBytes64(ei.rulePosSymbolsLat)
	if countCyr >= countLat {
		// конвертируем латиницу в кириллицу
		ei.directionConverting = 1
	} else {
		ei.directionConverting = 2
	}
	ei.converted = make([]YByte, len(ei.original))
	copy(ei.converted, ei.original)
	ei.convert()
}

// Пробуем конвертировать символы
func (ei *EncodeInformation) convert() {
	var (
		convertingRule uint64
	)
	if ei.directionConverting == 1 {
		// конвертируем латиницу в кириллицу
		convertingRule = ei.rulePosSymbolsLat
	} else if ei.directionConverting == 2 {
		convertingRule = ei.rulePosSymbolsCyr
	}
	for pos, yByte := range ei.original {
		mask := uint64(1 << pos)
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
