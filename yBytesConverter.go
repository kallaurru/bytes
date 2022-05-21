package yBytes

import (
	"fmt"
	"strconv"
	"strings"
)

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

//ConvertYBytes конвертируем в строку для сохранения в redis
func ConvertYBytes(bytes []YByte) string {
	var out string
	for _, b := range bytes {
		out = fmt.Sprintf("%s,%s", out, strconv.Itoa(int(b)))
	}
	return strings.Trim(out, ",")
}

//ConvertToYBytes конвертируем в строку для сохранения в redis
func ConvertToYBytes(in string) []YByte {
	els := strings.Split(in, ",")
	yBytes := make([]YByte, 0, 2)
	for _, n := range els {
		yBytes = append(yBytes, ConvertSymbolAsStringToYByte(n))
	}

	return yBytes
}

//ConvertSymbolAsStringToYByte редис вернет YByte символ как строку. Конвертируем в число и тип YByte.
func ConvertSymbolAsStringToYByte(n string) YByte {
	tmp, err := strconv.Atoi(n)
	if err != nil {
		return 0
	}
	return YByte(tmp)
}
