package yBytes

// GrWordFuncFistControl функции контроля сборки слова. С контролем начального разрешенного символа
type GrWordFuncFistControl func(code rune, isFirst bool) bool

// GrWordFunc функции контроля сборки слова. Без контроля первого символа
type GrWordFunc func(code rune) bool
