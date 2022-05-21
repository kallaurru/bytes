package yBytes

const (
	// Первая буква - заглавная
	optFirstCapital uint8 = 0x01 // Первая буква - заглавная.
	optAbbrev       uint8 = 0x02 // Это аббревиатура (все буквы заглавные)

	maxLengthWord   = 63 // Ограничитель на длину слова в кодировщике
	chanSizeDefault = 8  // Глубина канала по умолчанию
)
