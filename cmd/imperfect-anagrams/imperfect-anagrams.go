package main

import (
	"fmt"
	"math"
	"time"

	"github.com/joeyciechanowicz/letter-combinations/pkg/dictionary-tree"
)

type wordAnagramsPair struct {
	word          string
	anagramsCount int
}

func isWordAnagram(mainWord *dictionary_tree.WordDetails, wordToCheck *dictionary_tree.WordDetails) bool {
	if len(wordToCheck.Word) > len(mainWord.Word) {
		return false
	}

	for _, char := range wordToCheck.SortedSet {
		if wordToCheck.LetterCounts[char] > mainWord.LetterCounts[char] {
			return false
		}
	}

	return true
}

func searchSet(letters []rune, start int, head *dictionary_tree.Node, currentWord *dictionary_tree.WordDetails, anagramsCount *int) {
	if len(head.Words) > 0 {
		for i := 0; i < len(head.Words); i++ {
			if isWordAnagram(currentWord, head.Words[i]) {
				*anagramsCount += 1
			}
		}
	}

	for i := start; i < len(letters); i++ {
		if _, ok := head.Children[letters[i]]; ok {
			searchSet(letters, i+1, head.Children[letters[i]], currentWord, anagramsCount)
		}
	}
}

func findAnagrams(trie dictionary_tree.Node, wordChan <-chan dictionary_tree.WordDetails, anagramsForWord chan<- wordAnagramsPair) {
	for {
		currentWord, ok := <-wordChan
		if !ok {
			return
		}

		anagramsCount := 0
		searchSet(currentWord.SortedSet, 0, &trie, &currentWord, &anagramsCount)
		pair := wordAnagramsPair{currentWord.Word, anagramsCount}

		anagramsForWord <- pair
	}
}

func trackMostAnagrams(finished chan<- bool, numWords int, anagramsForWord <-chan wordAnagramsPair) {
	maxAnagramCount := 0
	maxWord := ""

	start := time.Now()
	ticks := 0
	total := float64(numWords)

	for {
		pair := <-anagramsForWord

		if pair.anagramsCount > maxAnagramCount {
			maxWord = pair.word
			maxAnagramCount = pair.anagramsCount
		}

		ticks += 1

		if ticks%500 == 0 {
			curr := float64(ticks)
			ratio := math.Min(math.Max(curr/total, 0), 1)
			percent := int32(math.Floor(ratio * 100))
			elapsed := float64(time.Since(start).Seconds())
			rate := curr / elapsed

			fmt.Printf("\r%d%% %d/wps", percent, int32(rate))
		}

		if ticks == numWords {
			fmt.Printf("\n\nLongest word: %s with %d imperfect-anagrams\n", maxWord, maxAnagramCount)
			finished <- true
			return
		}
	}
}

func findWordWithMostAnagrams(trie dictionary_tree.Node, words []dictionary_tree.WordDetails) {
	const numCpus = 4

	finished := make(chan bool)
	anagramsForWord := make(chan wordAnagramsPair, numCpus)
	wordChan := make(chan dictionary_tree.WordDetails, numCpus)

	go trackMostAnagrams(finished, len(words), anagramsForWord)

	for i := 0; i < numCpus; i++ {
		go findAnagrams(trie, wordChan, anagramsForWord)
	}

	go func() {
		for i := 0; i < len(words); i++ {
			wordChan <- words[i]
		}

		close(wordChan)
	}()

	<-finished

	close(anagramsForWord)
}

func main() {
	var trie, words = dictionary_tree.CreateDictionaryTree("./words_alpha.txt")
	//var trie, words = dictionary_tree.CreateDictionaryTree("./first_2000_words.txt")

	findWordWithMostAnagrams(trie, words)
}
