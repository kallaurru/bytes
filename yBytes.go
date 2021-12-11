package yBytes

/* Основной интерфейс */

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"sync"
)

const maxLengthWord = 63

// EncodeLine декодируем текстовую строку. Может быть просто слово, с невалидным символом в середине будет декодировано
// как 2 слова и более
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
		// получаем массив рун. Обычно это одна
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
	return ""
}

// ConvertingSymbols конвертируем символы, если все хорошо возвращаем true
func ConvertingSymbols(runeView *[]rune, symbols *[]YByte, toCyr bool) bool {
	for idx, symbol := range *runeView {
		var converted rune = 0
		if toCyr && isLatinLetter(symbol) {
			converted = convertLatToCyr(symbol)
		} else if !toCyr && isCyrillicLetter(symbol) {
			converted = convertCyrToLat(symbol)
		}
		if converted == 0 {
			// символ не конвертируемый прерываем цикл
			return false
		}
		(*symbols)[idx] = encodeSymbol(converted)
	}

	return true
}

func GetCountOnBytes64(rule uint64) int {
	if rule == 0 {
		return 0
	}

	out := 0
	mask := uint64(0xffffffff)
	lowBytes := rule & mask
	upBytes := (rule >> 32) & mask

	if lowBytes > 0 {
		out += GetCountOnBytes32(uint32(lowBytes))
	}
	if upBytes > 0 {
		out += GetCountOnBytes32(uint32(upBytes))
	}
	return out
}

func GetCountOnBytes32(rule uint32) int {
	if rule == 0 {
		return 0
	}

	out := 0
	lowBytes := rule & 0xffff
	upBytes := (rule >> 16) & 0xffff

	if lowBytes > 0 {
		out += GetCountOnBytes16(uint16(lowBytes))
	}
	if upBytes > 0 {
		out += GetCountOnBytes16(uint16(upBytes))
	}
	return out
}

func GetCountOnBytes16(rule uint16) int {
	if rule == 0 {
		return 0
	}
	out := 0
	lowByte := rule & 0xff
	upByte := (rule >> 8) & 0xff

	if lowByte > 0 {
		out += GetCountOnByte8(byte(lowByte))
	}

	if upByte > 0 {
		out += GetCountOnByte8(byte(upByte))
	}

	return out
}

func GetCountOnByte8(b byte) int {
	if b == 0 {
		return 0
	}

	out := 0
	for i := 0; i < 8; i++ {
		mask := 1 << i
		if byte(mask)&b > 0 {
			out++
		}
	}
	return out
}

// GetPositionsOnBytesUI - находим позиции включенных битов в числовых без знаковых линейках
func GetPositionsOnBytesUI(rule interface{}) []int {
	var i, length uint8
	var tmpRule uint64
	var pos int

	switch rule.(type) {
	case uint64:
		length = 8
		tmpRule = rule.(uint64)
	case uint32:
		length = 4
		tmp := rule.(uint32)
		tmpRule = uint64(tmp)
	case uint16:
		length = 2
		tmp := rule.(uint16)
		tmpRule = uint64(tmp)
	case uint8:
		length = 1
		tmp := rule.(uint8)
		tmpRule = uint64(tmp)
	default:
		// не верный тип
		return nil
	}
	if tmpRule == 0 {
		// пустая линейка
		return nil
	}
	out := make([]int, 0, 2)

	// цикл по байтам
	for i = 0; i < length; i++ {
		shift := i * 8
		part := (tmpRule >> shift) & 0xff
		if part == 0 {
			continue
		}
		// цикл по битам
		for bitI := 0; bitI < 8; bitI++ {
			mask := 1 << bitI
			if part&uint64(mask) > 0 {
				pos = (int(i) * 8) + bitI
				out = append(out, pos)
			}
		}
	}

	return out
}
