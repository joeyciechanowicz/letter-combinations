package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

type WordDetails struct {
	word         string
	sortedSet    []string
	letterCounts map[string]uint8
	canSkip      bool
}

type WordDetailsSlice []WordDetails

type Node struct {
	children map[string]*Node
	words    []*WordDetails
}

type Trie struct {
	root *Node
	size int
}

var words []WordDetails
var trie = Trie{
	root: &Node{children: make(map[string]*Node)},
	size: 0,
}

func newWordDetails(word string) WordDetails {
	var details = WordDetails{
		word,
		[]string{},
		make(map[string]uint8),
		false,
	}

	sortedLetters := strings.Split(word, "")
	sort.Strings(sortedLetters)

	details.sortedSet = append(details.sortedSet, sortedLetters[0])

	for _, char := range sortedLetters {
		details.letterCounts[char] += 1

		if details.sortedSet[len(details.sortedSet)-1] != char {
			details.sortedSet = append(details.sortedSet, char)
		}
	}

	return details
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func isWordAnagram(mainWord *WordDetails, wordToCheck *WordDetails) bool {
	if len(wordToCheck.word) > len(mainWord.word) {
		return false
	}

	for _, char := range wordToCheck.sortedSet {
		if wordToCheck.letterCounts[char] > mainWord.letterCounts[char] {
			return false
		}
	}

	return true
}

func searchSet(letters []string, start int, head *Node, currentWord *WordDetails, wordList *WordDetailsSlice) {
	if len(head.words) > 0 {
		for i := 0; i < len(head.words); i++ {
			if isWordAnagram(currentWord, head.words[i]) {
				*wordList = append(*wordList, *head.words[i])
			}
		}
	}

	for i := start; i < len(letters); i++ {
		if _, ok := head.children[letters[i]]; ok {
			searchSet(letters, i+1, head.children[letters[i]], currentWord, wordList)
		}
	}
}
//
//type wordAnagramsPair struct {
//	word string
//	anagrams []WordDetails
//}
//
//func findMostAnagramsForSlice(ticks chan<- int, ) {
//
//}

func findMostAnagrams() (string, []WordDetails) {
	start := time.Now()
	defer timeTrack(start, "Search for anagram")
	var maxWord string
	var maxAnagrams []WordDetails

	total := float64(len(words))

	fmt.Println()

	for i, currentWord := range words {
		if currentWord.canSkip {
			continue
		}

		if i % 1000 == 0 {
			curr := float64(i)

			ratio := math.Min(math.Max(curr/total, 0), 1)
			percent := int32(math.Floor(ratio * 100))
			//incomplete, complete, completeLength;
			elapsed := float64(time.Since(start).Seconds())
			eta := elapsed * (total / curr - 1)
			rate := curr / elapsed

			fmt.Printf("\r %d%% %d/wps %f seconds", percent, int32(rate), eta)
		}

		var anagrams WordDetailsSlice
		searchSet(currentWord.sortedSet, 0, trie.root, &currentWord, &anagrams)

		// todo refactor searchSet to be iterative instead of recursive
		//word, anagrams := findAnagrams(currentWord)

		if len(anagrams) > len(maxAnagrams) {
			maxWord = currentWord.word
			maxAnagrams = anagrams
		}
	}

	return maxWord, maxAnagrams
}

func buildTrie(filename string) {
	defer timeTrack(time.Now(), "Build trie and word details")

	nodeCount := 0

	fileHandle, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer Close(fileHandle)

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		var details WordDetails
		details = newWordDetails(scanner.Text())

		var head *Node
		head = trie.root

		for _, letter := range details.sortedSet {
			if _, ok := head.children[letter]; !ok {
				nodeCount++

				head.children[letter] = &Node{
					make(map[string]*Node),
					[]*WordDetails{},
				}
			}

			head = head.children[letter]
		}

		if len(head.children) > 0 {
			details.canSkip = true
		}

		words = append(words, details)
		head.words = append(head.words, &details)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Trie nodes: ", nodeCount)
}

func main() {
	buildTrie("./words_alpha.txt")
	//buildTrie("./first_2000_words.txt")
	word, anagrams := findMostAnagrams()

	fmt.Println("Longest word: ", word, " with ", len(anagrams), " anagrams")
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
