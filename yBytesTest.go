package yBytes

/**
Помощник
*/

// MakeEncodeInformationForTest декодированное слово Латиница. Последняя а - латинский символ
func MakeEncodeInformationForTest() *EncodeInformation {
	return &EncodeInformation{
		original:              []YByte{12, 1, 19, 9, 14, 9, 23, 97},
		rulePosNumbers:        0,
		rulePosSymbols:        0,
		rulePosCapitalSymbols: 1,
		rulePosSymbolsCyr:     127,
		rulePosSymbolsLat:     128,
		isNotConverting:       false,
		directionConverting:   0,
	}
}
