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
			hh |= uint32(code) << (shift * 8)
		}
		if ht > 0 {
			tmp := ht >> 8
			ht = uint32(code) << 24
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

	// формируем третий хэш если слово больше 8 букв
	if lenW > 8 {
		// из длины вычесть 4 буква ht и еще одну позицию, что встать на нужный индекс
		firstPoint := lenW - 4 - 1
		lastPoint := uint32(4)
		switch lenW {
		default:
			lastPoint = firstPoint - 2
		case 9:
			lastPoint = firstPoint
		case 10:
			lastPoint = firstPoint - 1
		}

		for i := lastPoint; i <= firstPoint; i++ {
			code := yBytes[i]
			packH3(&h3, uint32(code))
		}
	}
	// добавляем длину слова
	packSpecialH3(&h3, lenW, false, false)

	return hh, ht, h3
}

func GetAnotherVariants(h3 uint32) []uint32 {
	out := make([]uint32, 0, 2)

	// если пришло слово с включенным шестым или седьмым битом
	// значит нужно сформировать хэш без включенных бит (такое слово есть всегда)
	if h3&0xc0 != 0 {
		out = append(out, h3&0x3f)
	}

	// добавляем с включенным битом 6
	if h3&0x40 == 0 {
		out = append(out, h3|0x40)
	}
	// добавляем с включенным битом 8
	if h3&0x80 == 0 {
		out = append(out, h3|0x80)
	}

	return out
}
