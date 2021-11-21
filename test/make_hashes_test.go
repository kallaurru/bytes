package test

import (
	slBytes "github.com/kallaurru/bytes"
	"testing"
)

func TestMakeHashesWord(t *testing.T) {
	ei := slBytes.MakeEncodeInformationForTest()
	hh, ht, h3 := slBytes.MakeHashesWord(ei.ViewEncoding())
	if hh == 0 || ht == 0 || h3 == 0 {
		t.Error("hh = ", hh, "ht = ", ht, "h3 = ", h3)
	}
	t.Log("hh = ", hh, "ht = ", ht, "h3 = ", h3)
}
