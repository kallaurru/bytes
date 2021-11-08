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
	symbols []YByte
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
	return ei.symbols
}

// показываем как было в оригинале
func (ei *EncodeInformation) viewOriginal() string {
	for idx, yByte := range ei.symbols {

	}
	return ""
}

// показываем прописными буквами приведенными к одному типу символов если это возможно
func (ei *EncodeInformation) viewClassic() string {
	for idx, yByte := range ei.symbols {

	}
	return ""
}
