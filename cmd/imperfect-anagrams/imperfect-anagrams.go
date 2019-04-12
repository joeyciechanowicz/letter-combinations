package main

import (
	"fmt"
	"math"
	"time"

	"github.com/joeyciechanowicz/letter-combinations/pkg/dictionary-tree"
)

type wordAnagramsPair struct {
	word     string
	anagrams []dictionary_tree.WordDetails
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

func searchSet(letters []rune, start int, head *dictionary_tree.Node, currentWord *dictionary_tree.WordDetails, wordList *dictionary_tree.WordDetailsSlice) {
	if len(head.Words) > 0 {
		for i := 0; i < len(head.Words); i++ {
			if isWordAnagram(currentWord, head.Words[i]) {
				*wordList = append(*wordList, *head.Words[i])
			}
		}
	}

	for i := start; i < len(letters); i++ {
		if _, ok := head.Children[letters[i]]; ok {
			searchSet(letters, i+1, head.Children[letters[i]], currentWord, wordList)
		}
	}
}

func findAnagrams(trie dictionary_tree.Node, wordChan <-chan dictionary_tree.WordDetails, anagramsForWord chan<- wordAnagramsPair) {
	for {
		currentWord, ok := <-wordChan
		if !ok {
			return
		}

		var anagrams dictionary_tree.WordDetailsSlice
		searchSet(currentWord.SortedSet, 0, &trie, &currentWord, &anagrams)
		pair := wordAnagramsPair{currentWord.Word, anagrams}

		anagramsForWord <- pair
	}
}

func trackMostAnagrams(finished chan<- bool, numWords int, anagramsForWord <-chan wordAnagramsPair) {
	var maxAnagrams []dictionary_tree.WordDetails
	maxWord := ""

	start := time.Now()
	ticks := 0
	total := float64(numWords)

	for {
		pair := <-anagramsForWord

		if len(pair.anagrams) > len(maxAnagrams) {
			maxWord = pair.word
			maxAnagrams = pair.anagrams
		}

		ticks += 1

		if ticks%100 == 0 {
			curr := float64(ticks)
			ratio := math.Min(math.Max(curr/total, 0), 1)
			percent := int32(math.Floor(ratio * 100))
			elapsed := float64(time.Since(start).Seconds())
			rate := curr / elapsed

			fmt.Printf("\r%d%% %d/wps", percent, int32(rate))
		}

		if ticks == numWords {
			fmt.Printf("\n\nLongest word: %s with %d imperfect-anagrams\n", maxWord, len(maxAnagrams))
			finished <- true
			return
		}
	}
}

func findWordWithMostAnagrams(trie dictionary_tree.Node, words []dictionary_tree.WordDetails) {
	const numCpus = 4

	var wordChans [numCpus]chan dictionary_tree.WordDetails

	finished := make(chan bool)
	anagramsForWord := make(chan wordAnagramsPair, numCpus)

	go trackMostAnagrams(finished, len(words), anagramsForWord)

	for i := 0; i < numCpus; i++ {
		func(index int) {
			wordChans[index] = make(chan dictionary_tree.WordDetails, 4)
			go findAnagrams(trie, wordChans[index], anagramsForWord)
		}(i)
	}

	go func() {
		chanIndex := 0
		for i := 0; i < len(words); i++ {
			wordChans[chanIndex] <- words[i]

			chanIndex = (chanIndex + 1) % numCpus
		}

		for i := 0; i < numCpus; i++ {
			close(wordChans[i])
		}
	}()

	<-finished

	close(anagramsForWord)
}

func main() {
	var trie, words = dictionary_tree.CreateDictionaryTree("./words_alpha.txt")
	//var trie, words = dictionary_tree.CreateDictionaryTree("./first_2000_words.txt")

	findWordWithMostAnagrams(trie, words)
}
