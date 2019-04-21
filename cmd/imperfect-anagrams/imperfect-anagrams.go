package main

import (
	"flag"
	"fmt"
	"github.com/joeyciechanowicz/letter-combinations/pkg/rune-tree"
	"github.com/joeyciechanowicz/letter-combinations/pkg/stats"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

type wordAnagramsPair struct {
	word          string
	anagramsCount int
}

func isWord1AnagramOfWord2(word1 *rune_tree.WordDetails, word2 *rune_tree.WordDetails) bool {
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
func searchSet(letters []rune_tree.RuneCount, start int, head *rune_tree.Node, currentWord *rune_tree.WordDetails, anagramsCount *int) {
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
func findAnagrams(trie *rune_tree.Node, wordChan <-chan rune_tree.WordDetails, rateIncrements chan<- bool, maxAnagram chan<- wordAnagramsPair) {
	var maxAnagramCount int
	var maxWord string

	for {
		currentWord, ok := <-wordChan
		if !ok {
			maxAnagram <- wordAnagramsPair{maxWord, maxAnagramCount}
			return
		}

		anagramsCount := 0
		searchSet(currentWord.SortedRuneCounts, 0, trie, &currentWord, &anagramsCount)

		if anagramsCount > maxAnagramCount {
			maxWord = currentWord.Word
			maxAnagramCount = anagramsCount
		}

		rateIncrements <- true
	}
}

func walkTrie(node *rune_tree.Node, wordChan chan<- rune_tree.WordDetails) {
	if len(node.Children) == 0 {
		for _, word := range node.Words {
			wordChan <- *word
		}
	} else {
		for _, childNode := range node.Children {
			walkTrie(childNode, wordChan)
		}
	}
}


func findWordWithMostAnagrams(trie rune_tree.Node, words []rune_tree.WordDetails) {
	const numCpus = 8

	finished := make(chan bool)
	maxAnagrams := make(chan wordAnagramsPair)
	wordChan := make(chan rune_tree.WordDetails, numCpus)
	statUpdates := make(chan bool, numCpus)

	go stats.PrintRate(finished, statUpdates)

	for i := 0; i < numCpus; i++ {
		go findAnagrams(&trie, wordChan, statUpdates, maxAnagrams)
	}

	go func() {
		walkTrie(&trie, wordChan)

		close(wordChan)
	}()

	var maxAnagramCount int
	var maxWord string
	for i := 0; i < numCpus; i++ {
		select {
			case pair := <- maxAnagrams:
				if pair.anagramsCount > maxAnagramCount {
					maxWord = pair.word
					maxAnagramCount = pair.anagramsCount
				}
		}
	}

	finished <- true

	close(finished)
	close(maxAnagrams)
	close(statUpdates)

	fmt.Printf("\n\nLongest word: %s with %d imperfect-anagrams\n", maxWord, maxAnagramCount)
}

var cpuprofile = "cpu.prof"
//var cpuprofile = ""

//var memprofile = "mem.prof"
var memprofile = ""

func main() {
	//var trie, words = dictionary_tree.CreateRuneDictionaryTree("./words_alpha.txt")
	var trie, words = rune_tree.CreateRuneDictionaryTree("./words_no-names-or-places.txt")
	//var trie, words = dictionary_tree.CreateRuneDictionaryTree("./first_2000_words.txt")

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
