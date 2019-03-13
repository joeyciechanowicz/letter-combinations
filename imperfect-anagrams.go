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

type wordAnagramsPair struct {
	word     string
	anagrams []WordDetails
}

func findMostAnagramsForSlice(ticks chan<- int, anagrams chan<- wordAnagramsPair, wordSlice []WordDetails) {
	var maxWord string
	var maxAnagrams []WordDetails

	pingFrequency := len(wordSlice) / 100

	for i, currentWord := range wordSlice {
		if currentWord.canSkip {
			continue
		}

		if i % pingFrequency == 0 {
			ticks <- pingFrequency
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

	if len(words) % pingFrequency != 0 {
		ticks <- len(words) % pingFrequency
	}
	anagrams <- wordAnagramsPair{maxWord, maxAnagrams}
}

func printProgress(ticks <-chan int) {
	start := time.Now()
	totalTicks := 0
	total := float64(len(words))

	for {
		amount := <-ticks
		if amount == -1 {
			return
		}

		totalTicks += amount
		curr := float64(totalTicks)

		ratio := math.Min(math.Max(curr/total, 0), 1)
		percent := int32(math.Floor(ratio * 100))
		elapsed := float64(time.Since(start).Seconds())
		eta := elapsed * (total/curr - 1)
		rate := curr / elapsed

		fmt.Printf("\r%d%% %d/wps %f seconds    ", percent, int32(rate), eta)
		//default:
		//	time.Sleep(100 * time.Millisecond)
		//}
	}
}

func findMostAnagrams() (string, []WordDetails) {
	const numCpus = 4
	chunkSize := len(words) / numCpus

	ticks := make(chan int, 50)
	defer close(ticks)
	anagrams := make(chan wordAnagramsPair)
	defer close(anagrams)

	go printProgress(ticks)

	for i := 0; i < numCpus; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == numCpus-1 {
			end = len(words)
		}

		wordSlice := words[start:end]

		go findMostAnagramsForSlice(ticks, anagrams, wordSlice)
	}

	fmt.Println()
	receivedCount := 0
	var maxAnagramsPair wordAnagramsPair

	for {
		select {
		case pair := <-anagrams:
			if len(pair.anagrams) > len(maxAnagramsPair.anagrams) {
				maxAnagramsPair = pair
			}
			receivedCount++

			if receivedCount == numCpus {
				ticks <- -1
				return maxAnagramsPair.word, maxAnagramsPair.anagrams
			}

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
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

	fmt.Printf("\rLongest word: %s with %d anagrams    ", word, len(anagrams))
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
