package yBytes

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"sync"
)

// EncodeLine декодируем текстовую строку. Может быть, просто слово, с невалидным символом в середине будет декодировано
// как 2 слова и более.
func EncodeLine(line string, wordValidateFunc GrWordFuncFistControl) []*EncodeInformation {
	reader := strings.NewReader(line)
	scan := bufio.NewScanner(reader)
	scan.Split(bufio.ScanRunes)

	// размеры
	size := 8
	// канал основного цикла сборки слова
	chanRune := make(chan rune, size)
	// канал для возврата информации по собранному слову
	chanWordInfo := make(chan *EncodeInformation, size)
	wg := &sync.WaitGroup{}
	storage := make([]*EncodeInformation, 0, size)
	wg.Add(1)
	go EncodeFlowRunes(chanRune, chanWordInfo, wg, wordValidateFunc, false)
	wg.Add(1)
	go func(chIn chan<- *EncodeInformation, wg *sync.WaitGroup) {
		defer wg.Done()
		for ei := range chanWordInfo {
			storage = append(storage, ei)
		}
	}(chanWordInfo, wg)

	for scan.Scan() {
		// Получаем массив рун. Обычно это одна
		runes := bytes.Runes(scan.Bytes())
		for _, elem := range runes {
			chanRune <- elem
		}
	}

	if err := scan.Err(); err != nil {
		close(chanRune)
	}
	close(chanRune)
	// ожидаем завершения горутин
	wg.Wait()

	return storage
}

func EIBuilderPipe(
	channelIn <-chan rune,
	channelOut chan<- *EncodeInformation,
	wg *sync.WaitGroup) {

	var (
		ei *EncodeInformation
	)
	defer wg.Done()

	for runeCode := range channelIn {
		if runeCode == 0 {
			channelOut <- ei
			continue
		}

	}

}

// EncodeFlowRunes - кодируем слово "на лету"
func EncodeFlowRunes(
	channelIn <-chan rune,
	channelOut chan<- *EncodeInformation,
	wg *sync.WaitGroup,
	wordFunc GrWordFuncFistControl,
	chanOutExternalControl bool) {
	var (
		// флаг признак того, что начался процесс сборки и конвертации слова
		flgProcessing = false
		// флаг прочтения первой буквы. Нужен для проверочной функции сборки слова
		flgIsFirst        = true
		encodeInformation *EncodeInformation
		// текущая позиция в слове. Диапазон от 0 до 63
		position int = 0
	)
	defer wg.Done()
	for runeCode := range channelIn {
		isApproved := wordFunc(runeCode, flgIsFirst)

		if isApproved {
			// первым меняем флаг первого символа
			if flgIsFirst {
				flgIsFirst = false
			}

			// момент когда начат процесс сканирования слова
			if flgProcessing == false {
				flgProcessing = true
				encodeInformation = new(EncodeInformation)
				encodeInformation.SetMode(true)
			}

			// основной цикл сборки слова и кодирования
			if flgProcessing {
				if isCyrillicLetter(runeCode) {
					encodeInformation.rulePosSymbolsCyr |= 1 << position
				}
				if isLatinLetter(runeCode) {
					encodeInformation.rulePosSymbolsLat |= 1 << position
				}
				if isNumber(runeCode) {
					encodeInformation.rulePosNumbers |= 1 << position
				}
				if isCapitalSymbol(runeCode) {
					encodeInformation.rulePosCapitalSymbols |= 1 << position
				}
				if isSymbol(runeCode) {
					encodeInformation.rulePosSymbols |= 1 << position
				}
				yByte := encodeSymbol(runeCode)
				if yByte == 0 {
					// не критичная ситуация конвертации символа
					encodeInformation.errList = append(encodeInformation.errList, fmt.Errorf("not converted symbol in pos - %d", position))
				}
				encodeInformation.original = append(encodeInformation.original, yByte)
				position++
				// стопорим на 63 если было не стандартное слово
				if position > maxLengthWord {
					position = maxLengthWord
					encodeInformation.errList = append(encodeInformation.errList, fmt.Errorf("word len mehr that 63 original"))
				}
			}

		} else {
			// пришел не разрешенный символ.

			// момент когда происходил процесс сборки
			if flgProcessing {
				flgProcessing = false
				flgIsFirst = true
				// обнуляем позиции для сборки нового слова
				position = 0
				// отправляем информацию по декодированию наружу
				channelOut <- encodeInformation
			}

		}

	}
	// момент если шел процесс сборки слова и буфер закончился на последнем символе буфера
	if flgProcessing {
		flgProcessing = false
		flgIsFirst = true
		// обнуляем позиции для сборки нового слова
		position = 0
		// отправляем информацию по декодированию наружу
		channelOut <- encodeInformation
	}
	if chanOutExternalControl == false {
		close(channelOut)
	}

}

// EncodeSymbol интерфейс для внутренней функции конвертации отдельного символа
func EncodeSymbol(symbol rune) YByte {
	return encodeSymbol(symbol)
}

// DecodeSymbol интерфейс для внутренней функции декодирования отдельного символа
func DecodeSymbol(yb YByte, isCapital bool) rune {
	return decodeSymbol(yb, isCapital)
}

func DecodeWord(symbols []YByte, options uint8) string {
	if len(symbols) == 0 {
		return ""
	}
	var (
		flgFirst        = true
		flgFirstCapital = options&optFirstCapital > 0
		flgIsAbbrev     = options&optAbbrev > 0
		r               rune
	)
	out := make([]rune, 0, len(symbols))
	// если стоит признак аббревиатуры, то первый символ всегда заглавная
	flgFirstCapital = flgFirstCapital || flgIsAbbrev

	for _, yb := range symbols {
		if !flgFirst {
			r = DecodeSymbol(yb, flgIsAbbrev)
		} else {
			flgFirst = false
			r = DecodeSymbol(yb, flgFirstCapital)
		}

		out = append(out, r)
	}

	return string(out)
}
