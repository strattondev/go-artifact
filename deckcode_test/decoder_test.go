package deckcode_test

import (
	"github.com/strattonw/go-artifact/deckcode"
	"reflect"
	"testing"
)

func TestBlueRedExample(t *testing.T) {
	d, _ := deckcode.ParseDeck("ADCJQUQI30zuwEYg2ABeF1Bu94BmWIBTEkLtAKlAZakAYmHh0JsdWUvUmVkIEV4YW1wbGU_")

	if d.Name != "Blue/Red Example" {
		t.Fatalf("Expected \"Blue/Red Example\", actual \"%s\"", d.Name)
	}

	 var expectedHeroes = map[int]int {
	 	4003: 1,
	 	10006: 1,
	 	10030: 1,
	 	10033: 3,
	 	10065: 2,
	 }

	 var actualHeroes = make(map[int]int)

	 for _, h := range d.Heroes {
	 	actualHeroes[h.Id] = h.Turn
	 }

	 if !reflect.DeepEqual(expectedHeroes, actualHeroes) {
	 	t.Fatalf("Expected heroes %+v, actual %+v", expectedHeroes, d.Heroes)
	 }

	 var expectedCards = map[int]int {
		 3000: 2,
		 3001: 2,
		 10132: 3,
		 10157: 3,
		 10191: 2,
		 10203: 2,
		 10212: 2,
		 10223: 1,
		 10307: 3,
		 10344: 3,
		 10366: 3,
		 10402: 3,
		 10411: 3,
		 10418: 3,
		 10425: 3,
	 }

	 var actualCards = make(map[int]int)

	 for _, c := range d.Cards {
	 	actualCards[c.Id] = c.Count
	 }

	 if !reflect.DeepEqual(expectedCards, actualCards) {
		 t.Fatalf("Expected cards %+v, actual %+v", expectedHeroes, d.Heroes)
	 }
}
