package test

import (
	yBytes "github.com/kallaurru/bytes"
	"strings"
	"testing"
)

func TestEncodingProcess(t *testing.T) {
	reader := strings.NewReader("kaserg@mail, победа, достаток\n")
	eiList := yBytes.EncoderWords(reader, yBytes.ValidIsClassicWord)
	if len(eiList) != 4 {
		t.Error("count object not equal")
	}
}

func TestFilteringPhoneNumber(t *testing.T) {
	reader := strings.NewReader("+7 (926) 067 33 08")
	ei := yBytes.FilterWord(reader, yBytes.FIsNumericWord)
	if ei == nil {
		t.Error("ei is nil")
	}
}
