package test

import (
	yBytes "github.com/kallaurru/bytes"
	"log"
	"strings"
	"testing"
)

func TestEncodingProcess(t *testing.T) {
	reader := strings.NewReader("kaserg@mail, победа, достаток, счастье, удача, радость, сорок пять лет")
	eiList := yBytes.EncoderWords(reader, yBytes.ValidIsClassicWord)
	if len(eiList) != 10 {
		t.Error("count object not equal")
	}
	for _, ei := range eiList {
		log.Println(ei.ViewClassic())
	}
}

func TestFilteringPhoneNumber(t *testing.T) {
	reader := strings.NewReader("+7 (926) 067-33-08")
	ei := yBytes.FilterWord(reader, yBytes.FIsNumericWord)
	reader = strings.NewReader("+7 (926) 067-33-08")
	eiPhoneFormat := yBytes.FilterWord(reader, yBytes.FIsPhoneNumberFormattedInter)
	log.Println(ei.ViewClassic())
	log.Println(eiPhoneFormat.ViewClassic())
	if ei == nil {
		t.Error("ei is nil")
	}
}
