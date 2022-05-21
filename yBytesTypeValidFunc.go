package yBytes

// ValidationWordFunc функции контроля сборки слова. С контролем начального разрешенного символа
type ValidationWordFunc func(code rune, isFirst bool) bool

// FilterWordFunc функции контроля сборки слова. Без контроля первого символа
type FilterWordFunc func(code rune) bool
