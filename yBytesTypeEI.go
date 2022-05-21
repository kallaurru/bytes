package yBytes

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math"
	"strconv"
)

// EncodeInformation информация по декодированному слову
// все свойства которые начинаются с rule, являются линейками позиционирования нужных символов
// !!! Все позиционирование в линейках начинается с ноля
// по умолчанию считается что нормальное слово  !!!! не будет больше 64 символов
type EncodeInformation struct {
	// Храним список полученных ошибок в результате конвертации
	errList []error
	// Представление кодированного слова в виде массива универсальной кодировки
	original []YByte
	// Представление конвертированного слова
	converted []YByte
	// Обработано символов. Важно для режима фильтрации
	processingSymbols int
	// Представление в виде рун. Кэш слова для быстрого перевода в строку. Хранится в прописных символах
	runes []rune
	// Позиции цифр в слове
	rulePosNumbers uint64
	// Позиции допустимых не буквенно-цифровых символов
	rulePosSymbols uint64
	// Позиции прописных букв в слове. Для работы с аббревиатурами
	rulePosCapitalSymbols uint64
	// Позиции букв кириллицы в слове
	rulePosSymbolsCyr uint64
	// Позиции букв латиницы в слове
	rulePosSymbolsLat uint64
	// Информация создавалась в режиме разборки слов, false - в режиме фильтрации
	isParsingMode bool
	// False - существующие в слове символы оппозитной кодировки не могут быть конвертированы.
	isNotConverting bool
	// Направление конвертирования -0 не потребовалось, 1 - из латиницы в кириллицу, 2 - из кириллицы в латиницу
	directionConverting uint8
	// Позиция в потоке, для позиционирования в деревьях или графах
	posInFlow uint64
	// Для унификации и построения индексов. Используем uuid v.4
	uuid string
}

/** Признаки слова */

func (ei *EncodeInformation) IsClassicWord() bool {
	return ei.rulePosNumbers == 0 && ei.rulePosSymbols == 0
}

func (ei *EncodeInformation) IsCyrillicWord() bool {
	return ei.rulePosSymbolsLat == 0 && ei.rulePosSymbolsCyr > 0
}

func (ei *EncodeInformation) IsLatWord() bool {
	return !ei.IsCyrillicWord()
}

func (ei *EncodeInformation) IsNumber() bool {
	lenW := len(ei.original)
	max := uint64(math.Pow(2, float64(lenW)))

	return ei.rulePosNumbers == max-1
}

func (ei *EncodeInformation) SetMode(isParsing bool) {
	ei.isParsingMode = isParsing
}

func (ei *EncodeInformation) AddUuid() {
	ei.uuid = uuid.NewV4().String()
}

func (ei *EncodeInformation) AddFlowPosition(pos uint64) {
	ei.posInFlow = pos
}

func (ei *EncodeInformation) IsParsingMode() bool {
	return ei.isParsingMode
}

func (ei *EncodeInformation) ProcessingSymbols() int {
	return ei.processingSymbols
}

func (ei *EncodeInformation) Len() int {
	return len(ei.original)
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

// Подготовка под возможную конвертацию
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

func (ei *EncodeInformation) PrepareToRedis() map[string]string {
	m := make(map[string]string)
	m["uuid"] = ei.uuid
	m["original"] = ConvertYBytes(ei.original)
	m["converted"] = ConvertYBytes(ei.converted)
	m["pos_in_flow"] = fmt.Sprintf("%v", ei.posInFlow)
	m["processing_symbols"] = fmt.Sprintf("%d", ei.processingSymbols)
	m["rule_pos_numbers"] = fmt.Sprintf("%v", ei.rulePosNumbers)
	m["rule_pos_symbols"] = fmt.Sprintf("%v", ei.rulePosSymbols)
	m["rule_pos_capital_symbols"] = fmt.Sprintf("%v", ei.rulePosCapitalSymbols)
	m["rule_pos_symbols_cyr"] = fmt.Sprintf("%v", ei.rulePosSymbolsCyr)
	m["rule_pos_symbols_lat"] = fmt.Sprintf("%v", ei.rulePosSymbolsLat)
	m["direction_converting"] = fmt.Sprintf("%v", ei.directionConverting)
	if ei.isNotConverting {
		m["is_not_converting"] = "true"
	} else {
		m["is_not_converting"] = "false"
	}
	if ei.isParsingMode {
		m["is_parsing_mode"] = "true"
	} else {
		m["is_parsing_mode"] = "false"
	}

	return m
}

func (ei *EncodeInformation) UpdateFromRedis(m map[string]string) error {
	if val, ok := m["uuid"]; ok {
		ei.uuid = val
	}
	if val, ok := m["original"]; ok {
		ei.original = ConvertToYBytes(val)
	}
	if val, ok := m["converted"]; ok {
		ei.converted = ConvertToYBytes(val)
	}
	if val, ok := m["processing_symbols"]; ok {
		intVal, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		ei.processingSymbols = int(intVal)
	}
	if val, ok := m["rule_pos_numbers"]; ok {
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		ei.rulePosNumbers = intVal
	}

	if val, ok := m["pos_in_flow"]; ok {
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		ei.posInFlow = intVal
	}

	if val, ok := m["rule_pos_symbols"]; ok {
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		ei.rulePosSymbols = intVal
	}

	if val, ok := m["rule_pos_capital_symbols"]; ok {
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		ei.rulePosCapitalSymbols = intVal
	}

	if val, ok := m["rule_pos_symbols_cyr"]; ok {
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		ei.rulePosSymbolsCyr = intVal
	}

	if val, ok := m["rule_pos_symbols_lat"]; ok {
		intVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		ei.rulePosSymbolsCyr = intVal
	}

	if val, ok := m["direction_converting"]; ok {
		if val == "true" {
			ei.isNotConverting = true
		} else {
			ei.isNotConverting = false
		}
	}

	if val, ok := m["parsing_mode"]; ok {
		if val == "true" {
			ei.isParsingMode = true
		} else {
			ei.isParsingMode = false
		}
	}

	return nil
}
