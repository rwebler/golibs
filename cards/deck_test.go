package main

import "testing"
import "os"

func TestNewDeck(t *testing.T) {
	d := newDeck()
	assertions(d, t)
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	testfile := "_testdeck"
	os.Remove(testfile)
	d := newDeck()
	d.saveToFile(testfile)
	d2 := newDeckFromFile(testfile)

	assertions(d2, t)

	os.Remove(testfile)
}

func assertions(d deck, t *testing.T) {
	xl := 52
	if len(d) != xl {
		t.Errorf("Expected deck length of %d but got %d", xl, len(d))
	}

	xf := "Ace of Clubs"
	if d[0] != xf {
		t.Errorf("Expected first card to be %s but got %s", xf, d[0])
	}

	xl2 := "King of Spades"
	if d[len(d)-1] != xl2 {
		t.Errorf("Expected last card to be %s but got %s", xl2, d[len(d)-1])
	}
}
