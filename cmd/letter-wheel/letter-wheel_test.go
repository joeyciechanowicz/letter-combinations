package main

import (
	"os"
	"testing"

	"github.com/joeyciechanowicz/letter-combinations/pkg/int-tree"
)

var word = int_tree.NewWordDetails("hello")

/*
a a a
o e h
l l e
 */
var wheelWord = int_tree.NewWordDetails("aaaeehllo")
var mainLetter = int_tree.ToAlphabetIndex(int_tree.ToRune("e"))
var wheel = Wheel{mainLetter, wheelWord.SortedLetterCounts}
var rawWheel = [WHEEL_SIZE - 1]int{0, 0, 0, 4, 7, 11, 11, 14}

var trie int_tree.Node

func TestMain(m *testing.M) {
	trie, _ = int_tree.CreateIntDictionaryTree("../../3-to-9-letter-words.txt")

	code := m.Run()
	os.Exit(code)
}

/*
canWordBeSpeltFromWheel
 */

func TestCanWordBeSpeltFromWheel(t *testing.T) {
	canSpell := canWordBeSpeltFromWheel(word.SortedLetterCounts, wheel)

	if !canSpell {
		t.Errorf("Word could not be spelt.")
	}
}

func BenchmarkCanWordBeSpeltFromWheel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		canWordBeSpeltFromWheel(word.SortedLetterCounts, wheel)
	}
}

/*
findWordsForWheel
 */
func TestFindWordsForWheel(t *testing.T) {
	var wheelCount = 0
	findWordsForWheel(&trie, 0, wheel, &wheelCount)

	if wheelCount != 15 {
		t.Errorf("Count was incorrect, got: %d, want: %d.", wheelCount, 15)
	}
}

func BenchmarkFindWordsForWheel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var wheelCount = 0
		findWordsForWheel(&trie, 0, wheel, &wheelCount)
	}
}

/*
findWords
 */
func TestFindWords(t *testing.T) {
	wheelChan := make(chan [WHEEL_SIZE - 1]int)
	stats := make(chan bool)
	maxWheelChan := make(chan wordCountForWheel)

	// Dump stats
	go func() {
		for {
			<-stats
		}
	}()

	go findWords(&trie, wheelChan, stats, maxWheelChan)

	wheelChan <- rawWheel
	close(wheelChan)

	result := <-maxWheelChan

	if result.wordsCount != 67 {
		t.Errorf("Word count was incorrect, got: %d, want: %d.", result.wordsCount, 67)
	}

	if int_tree.Alphabet[result.wheel.MainLetter] != rune("s"[0]) {
		t.Errorf("Main letter was incorrect, got: %s, want: %s.", string(int_tree.Alphabet[result.wheel.MainLetter]), "s")
	}
}

func BenchmarkFindWords(b *testing.B) {
	b.ReportAllocs()

	wheelChan := make(chan [WHEEL_SIZE - 1]int)
	stats := make(chan bool)
	maxWheelChan := make(chan wordCountForWheel)

	go findWords(&trie, wheelChan, stats, maxWheelChan)

	for i := 0; i < b.N; i++ {
		wheelChan <- rawWheel

		for j := 0; j < 26; j++ {
			<-stats
		}
	}

	close(wheelChan)
}