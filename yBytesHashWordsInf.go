package yBytes

// MakeHashesWord - формируем три хэша по слову.
// hh - хэш заголовка, первые 4 буквы
// ht - хэш хвоста слова последние 4 буквы
// h3 - хэш доп информации три буквы с хвоста после ht информация о длине и ударениях в одинаково
// пишущихся словах
func MakeHashesWord(yBytes []YByte) (uint32, uint32, uint32) {
	var hh, ht, h3, lenW uint32

	for idx, code := range yBytes {
		if idx < 4 {
			shift := 3 - idx
			hh |= uint32(code) << shift
		}
		if ht > 0 {
			tmp := ht >> 8
			ht |= uint32(code) << 24
			ht |= tmp
		} else {
			ht |= uint32(code) << 24
		}
		if idx > 7 {
			// формируем третий хэш транзитным путем
			packH3(&h3, uint32(code))
		}
		lenW++
	}

	// добавляем длину слова
	packSpecialH3(&h3, lenW, false, false)

	return hh, ht, h3
}
