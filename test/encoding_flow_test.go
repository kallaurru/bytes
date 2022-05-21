package test

import (
	"bufio"
	"bytes"
	yBytes "github.com/kallaurru/bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
)

// тестируем кодирование слов поступающих из сканера в потоке
func TestFlowEncodingProcess(t *testing.T) {
	// входящий текст который требуется разобрать на слова.
	// В слове латиница последняя буква латинская
	reader := strings.NewReader("kaserg@mail, победа, достаток\n")
	scan := bufio.NewScanner(reader)
	scan.Split(bufio.ScanRunes)

	// каналы.
	channelSize := 8
	mapSize := 16
	// канал основного цикла сборки слова
	chanGeneral := make(chan rune, channelSize)
	// канал для возврата информации по собранному слову
	chanWordInfo := make(chan *yBytes.EncodeInformation, channelSize)
	wg := &sync.WaitGroup{}
	storageMutex := &sync.RWMutex{}

	storage := make(map[string]*yBytes.EncodeInformation, mapSize)

	// запускаем функцию сборки слова в потоке
	wg.Add(1)
	go yBytes.EncodeFlowRunes(chanGeneral, chanWordInfo, wg, yBytes.ValidIsAdvWord, false)

	// запускаем функцию основного потока для складирования слов
	wg.Add(1)
	go func(channelIn <-chan *yBytes.EncodeInformation, wg *sync.WaitGroup, mutex *sync.RWMutex) {
		defer wg.Done()
		// читаем поступающую информацию из канала
		for ei := range channelIn {
			// достаем из информации по кодированию
			word := ei.ViewClassic()
			mutex.Lock()
			storage[word] = ei
			mutex.Unlock()
		}

	}(chanWordInfo, wg, storageMutex)
	// начинаем основной цикл сканирования
	for scan.Scan() {
		// получаем массив рун. Обычно это одна
		runes := bytes.Runes(scan.Bytes())
		for _, elem := range runes {
			chanGeneral <- elem
		}
	}

	if err := scan.Err(); err != nil {
		close(chanGeneral)
		t.Error(err)
	}
	close(chanGeneral)
	// ожидаем завершения горутин
	wg.Wait()
	for word, ei := range storage {
		hh, ht, h3 := yBytes.MakeHashesWord(ei.ViewEncoding())
		t.Log("word - ", word)
		t.Log("hh - ", hh, "ht - ", ht, "h3 - ", h3)
		t.Log(" ----------------- ----------------- ------------ ")
	}
	t.Log("map has", len(storage), "elements")
	assert.Equal(t, 3, len(storage))
}

/**
Копии функций для тестирования в одном потоке
*/

func TestConvertEncodingInformation(t *testing.T) {
	// делаем на примере слова Латиница. Последняя "а" должна быть латинской
	ei := yBytes.MakeEncodeInformationForTest()
	word := ei.ViewClassic()
	assert.Equal(t, word, "латиница", "не корректная конвертация")
}
