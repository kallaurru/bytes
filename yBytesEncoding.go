package yBytes

import (
	"bufio"
	"bytes"
	slUtils "github.com/kallaurru/utils"
	"sync"
)

var latToCyrMap *map[rune]rune
var cyrToLatMap *map[rune]rune

func decodeSymbol(yb YByte, isCapital bool) rune {
	var diffUpLowerCase rune = 32
	var diffCyr rune = 1071

	// конвертируем исключения (ё ь ъ)
	switch yb {
	case 127:
		if isCapital {
			return 1025
		} else {
			return 1105
		}
	case 128:
		if isCapital {
			return 1066
		} else {
			return 1098
		}
	case 129:
		if isCapital {
			return 1068
		} else {
			return 1100
		}
	}
	if yb < 32 {
		// кириллический символ
		var symbol rune
		if yb > 26 {
			symbol = rune(yb) + diffCyr + 1
		} else {
			symbol = rune(yb) + diffCyr
		}
		if isCapital {
			symbol -= diffUpLowerCase
		}
		return symbol
	}
	// конвертируем цифры символы и маленькие латинские как есть
	if (yb >= 32 && yb <= 64) || (yb >= 96 && yb <= 126) {
		if isCapital {
			return rune(yb) - diffUpLowerCase
		}
		return rune(yb)
	}
	if yb >= 130 && yb < 255 {
		// доп информация которая не может быть корректно конвертирована в буквы
		return 1
	}
	// возможно что то упустил не штатная ситуация
	return 0
}

// encodeSymbol кодируем текущие символы слова для внутреннего использования системой
// все заглавные кодируем в прописные буквы
// кодировка:
//	- 0 оставляем как спец код.
//	- 1 - 31 укладываем кириллицу без ё ь ъ
//	- 32 - 126 оставляем как есть в ascii
//	- 127 - 129 вставляем ё ъ ь
//  - 130 - 255 - зарезервированы на будущее, возможно будем паковать доп инфу по слову
//  в данной функции не применяются
// Если возврат 0 значит конвертировался какой-то спец символ
func encodeSymbol(symbol rune) YByte {
	var diffUpLowerCase rune = 32
	var diffCyr rune = 1071

	// конвертируем исключения (ё ь ъ)
	switch symbol {
	// заглавная
	case 1025:
	case 1105:
		return 127
	// заглавная
	case 1066:
	case 1098:
		return 128
	// заглавная
	case 1068:
	case 1100:
		return 129
	}

	// конвертируем заглавные кириллические и латинские utf-8
	if (symbol >= 1040 && symbol <= 1071) || (symbol >= 65 && symbol <= 90) {
		symbol += diffUpLowerCase
	}
	// конвертируем маленькие кириллические utf-8
	if symbol >= 1072 && symbol <= 1103 {
		if symbol > 1098 {
			return YByte(symbol - diffCyr - 1)
		}
		return YByte(symbol - diffCyr)
	}
	// конвертируем цифры символы и маленькие латинские как есть
	if (symbol >= 32 && symbol <= 64) || (symbol >= 96 && symbol <= 126) {
		return YByte(symbol)
	}

	return 0
}

func processAddToStorage(
	chanEI <-chan *EncodeInformation,
	wg *sync.WaitGroup,
	ptrStorage *[]*EncodeInformation) {

	defer wg.Done()

	var counter uint64 = 0

	for ei := range chanEI {
		if ei == nil {
			continue
		}
		ei.AddFlowPosition(counter)
		ei.SetMode(true)
		*ptrStorage = append(*ptrStorage, ei)
		counter++
	}

}

func processScanParserMode(
	chanRune chan<- rune,
	wg *sync.WaitGroup,
	scan *bufio.Scanner,
	wordValidateFunc ValidationWordFunc) {

	var (
		flgIsFirst      = true
		flgIsProcessing = false
	)
	defer wg.Done()

	for scan.Scan() {
		// Получаем массив рун. Обычно это одна
		runes := bytes.Runes(scan.Bytes())
		for _, elem := range runes {
			isValid := wordValidateFunc(elem, flgIsFirst)
			if isValid {
				if flgIsFirst {
					slUtils.SwitcherRefOff(&flgIsFirst)
					if !flgIsProcessing {
						slUtils.SwitcherRefOn(&flgIsProcessing)
					}
				}
				chanRune <- elem
			} else if flgIsProcessing {
				slUtils.SwitcherRefOff(&flgIsProcessing)
				slUtils.SwitcherRefOn(&flgIsFirst)
				chanRune <- 0
			}
		}
	}

	close(chanRune)
}

func processScanFilterMode(
	chanRune chan<- rune,
	wg *sync.WaitGroup,
	scan *bufio.Scanner,
	wordFilterFunc FilterWordFunc) {

	var (
		flgIsFirst = true
		counter    = 0
	)
	defer wg.Done()

	for scan.Scan() {
		// Получаем массив рун. Обычно это одна
		runes := bytes.Runes(scan.Bytes())
		for _, elem := range runes {
			isValid := wordFilterFunc(elem)
			if isValid {
				if flgIsFirst {
					slUtils.SwitcherRefOff(&flgIsFirst)
				}
				chanRune <- elem
			} else {
				continue
			}
			if counter > maxLengthWord {
				break
			}
			counter++
		}
	}

	close(chanRune)
}
