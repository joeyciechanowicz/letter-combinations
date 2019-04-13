package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/joeyciechanowicz/letter-combinations/pkg/dictionary-tree"
)

type wordAnagramsPair struct {
	word          string
	anagramsCount int
}

func isWord1AnagramOfWord2(word1 *dictionary_tree.WordDetails, word2 *dictionary_tree.WordDetails) bool {
	if len(word1.Word) > len(word2.Word) {
		return false
	}

	// Iterate the main words runes, setting an index for the other word
	i := -1
	for j := 0; j < len(word1.SortedRuneCounts); j++ {
		runeAndCount := word1.SortedRuneCounts[j]

		// move the other words index along until we find a letter that matches
		// returning false if we reach the end or the rune counts are incorrect
		for {
			i++
			if i == len(word2.SortedRuneCounts) {
				return false
			}

			if runeAndCount.Letter == word2.SortedRuneCounts[i].Letter {
				if runeAndCount.Count <= word2.SortedRuneCounts[i].Count {
					break
				} else {
					return false
				}
			}
		}
	}

	return true
}

/**
Recursively searches the trie for any words that can be spelt using the given set of letter

Works by taking the trie head and set of letters and trying to walk the trie using that set.
However we iterate the letters each time to allow us to search and find words that are anagrams
i.e. then (e,h,n,t) can spell hen (e,h,n) and net (e,n,t). By iterating AND recursing we check all branches of the trie
that could contain anagrams.
 */
func searchSet(letters []dictionary_tree.RuneCount, start int, head *dictionary_tree.Node, currentWord *dictionary_tree.WordDetails, anagramsCount *int) {
	if len(head.Words) > 0 {
		for i := 0; i < len(head.Words); i++ {
			if isWord1AnagramOfWord2(head.Words[i], currentWord) {
				*anagramsCount += 1
			}
		}
	}

	for i := start; i < len(letters); i++ {
		if _, ok := head.Children[letters[i].Letter]; ok {
			searchSet(letters, i+1, head.Children[letters[i].Letter], currentWord, anagramsCount)
		}
	}
}

/**
Takes words off a channel and finds all the anagrams for that word
 */
func findAnagrams(trie dictionary_tree.Node, wordChan <-chan dictionary_tree.WordDetails, anagramsForWord chan<- wordAnagramsPair) {
	for {
		currentWord, ok := <-wordChan
		if !ok {
			return
		}

		anagramsCount := 0
		searchSet(currentWord.SortedRuneCounts, 0, &trie, &currentWord, &anagramsCount)
		pair := wordAnagramsPair{currentWord.Word, anagramsCount}

		anagramsForWord <- pair
	}
}

/**
Takes a channel of word/anagrams pair and tracks which word has the most anagrams
also prints run stats continuously
 */
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
	const numCpus = 8

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

var cpuprofile = "cpu.prof"
//var cpuprofile = ""

//var memprofile = "mem.prof"
var memprofile = ""

func main() {
	//var trie, words = dictionary_tree.CreateDictionaryTree("./words_alpha.txt")
	var trie, words = dictionary_tree.CreateDictionaryTree("./words_no-names-or-places.txt")
	//var trie, words = dictionary_tree.CreateDictionaryTree("./first_2000_words.txt")

	flag.Parse()
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	findWordWithMostAnagrams(trie, words)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
