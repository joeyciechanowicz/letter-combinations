package main

import (
	"flag"
	"fmt"
	"github.com/joeyciechanowicz/letter-combinations/pkg/int-tree"
	"github.com/joeyciechanowicz/letter-combinations/pkg/stats"
	"log"
	"os"
	"runtime/pprof"
)

const WHEEL_SIZE = 9
const TOTAL_WHEELS = 52451256

type wheelSolutions struct {
	wheel     Wheel
	wordCount int
}

type Wheel struct {
	MainLetter   int
	LetterCounts []int_tree.LetterCount
}

func canWordBeSpeltFromWheel(word []int_tree.LetterCount, wheel Wheel) bool {
	seenMainLetter := false

	// Iterate the main words runes, setting an index for the other word
	i := -1

	for j := 0; j < len(word); j++ {
		letterAndCount := word[j]

		if letterAndCount.Letter == wheel.MainLetter {
			seenMainLetter = true
		}

		// move the other words index along until we find a letter that matches
		// returning false if we reach the end or the rune counts are incorrect
		for {
			i++
			if i == len(wheel.LetterCounts) {
				return false
			}

			if letterAndCount.Letter == wheel.LetterCounts[i].Letter {
				if letterAndCount.Count <= wheel.LetterCounts[i].Count {
					break
				} else {
					return false
				}
			}
		}
	}

	return seenMainLetter
}

func findWordsForWheel(set []int, head *int_tree.Node, start int, currentWheel Wheel, wordCount *int) {
	if len(head.Words) > 0 {
		for i := 0; i < len(head.Words); i++ {
			if canWordBeSpeltFromWheel(head.Words[i].SortedLetterCounts, currentWheel) {
				*wordCount += 1
			}
		}
	}

	for i := start; i < len(set); i++ {
		if _, ok := head.Children[set[i]]; ok {
			findWordsForWheel(set, head.Children[set[i]], i+1, currentWheel, wordCount)
		}
	}
}

func countLetters(wheel [WHEEL_SIZE - 1]int) [26]byte {
	var letterCounts [26]byte

	// Abuse that we know exactly the range of our data (26 letters in the alphabet)
	for i := 0; i < WHEEL_SIZE-1; i++ {
		letterCounts[wheel[i]]++
	}

	return letterCounts
}

/**
Takes the 8 surrounding wheel runes off a channel, iterates the centre 26 letters and finds the word-count for each wheel
 */
func findWords(trie *int_tree.Node, wheelChan <-chan [WHEEL_SIZE - 1]int, stats chan<- bool, maxWheelChan chan<- wheelSolutions) {
	var maxWordCount int
	var maxWheel Wheel

	for {
		currentWheel, ok := <-wheelChan
		if !ok {
			maxWheelChan <- wheelSolutions{maxWheel, maxWordCount}
			return
		}

		letterCounts := countLetters(currentWheel)

		for mainLetter := 0; mainLetter < 26; mainLetter++ {
			var compressedLetterCounts []int_tree.LetterCount

			letterCounts[mainLetter]++

			for letter, count := range letterCounts {
				if count > 0 {
					compressedLetterCounts = append(compressedLetterCounts, int_tree.LetterCount{
						Letter: letter,
						Count:  count,
					})
				}
			}

			wordCount := 0
			wheel := Wheel{mainLetter, compressedLetterCounts}
			findWordsForWheel(currentWheel[:], trie, 0, wheel, &wordCount)

			if wordCount > maxWordCount {
				maxWheel = wheel
				maxWordCount = wordCount
			}

			letterCounts[mainLetter]--

			stats <- true
		}
	}
}

func combinationRepetitionUtil(wheelChan chan<- [WHEEL_SIZE - 1]int, chosen [WHEEL_SIZE - 1]int, index, r, start, end int) {
	// Since index has become r, current combination is complete
	if index == r {
		wheelChan <- chosen
		return
	}

	// One by one choose all elements (without considering if the element is already chosen or not)
	for i := start; i <= end; i++ {
		chosen[index] = i
		combinationRepetitionUtil(wheelChan, chosen, index+1, r, i, end)
	}
}

// Recursively calculates all combinations of the alphabet for the last 8 chars of a wheel.
// We then iterate the alphabet and append the combination 26 times
func combinationRepetition(wheelChan chan [WHEEL_SIZE - 1]int) {
	var chosen [WHEEL_SIZE - 1]int

	if testMode {
		combinationRepetitionUtil(wheelChan, chosen, 0, WHEEL_SIZE-2, 0, 14)
	} else {
		combinationRepetitionUtil(wheelChan, chosen, 0, WHEEL_SIZE-2, 0, 26-1)
	}
}

func printOutput(wheel Wheel, numFound int) {
	var letters []string
	for _, letterCounts := range wheel.LetterCounts {
		if letterCounts.Letter == wheel.MainLetter {
			for j := 0; j < int(letterCounts.Count - 1); j++ {
				letters = append(letters, string(int_tree.Alphabet[letterCounts.Letter]))
			}
		} else {
			for j := 0; j < int(letterCounts.Count); j++ {
				letters = append(letters, string(int_tree.Alphabet[letterCounts.Letter]))
			}
		}
	}

	fmt.Printf("\n┏━━━━━━━━━━━┓\n")
	fmt.Printf("┃ %s   %s   %s ┃\n", letters[0], letters[1], letters[2])
	fmt.Printf("┃   ┏━━━┓   ┃\n")
	fmt.Printf("┃ %s ┃ %s ┃ %s ┃  Found %d\n", letters[7], string(int_tree.Alphabet[wheel.MainLetter]), letters[3], numFound)
	fmt.Printf("┃   ┗━━━┛   ┃\n")
	fmt.Printf("┃ %s   %s   %s ┃\n", letters[6], letters[5], letters[4])
	fmt.Printf("┗━━━━━━━━━━━┛\n")
}

func findBestLetterWheel(trie int_tree.Node, details []int_tree.WordDetails) (Wheel, int) {
	const NUM_CPUS = 8

	finished := make(chan bool)
	maxWheelChan := make(chan wheelSolutions)
	outerWheelChan := make(chan [WHEEL_SIZE - 1]int, NUM_CPUS)
	statUpdates := make(chan bool, NUM_CPUS)

	go stats.PrintProgress(finished, statUpdates, TOTAL_WHEELS)

	for i := 0; i < NUM_CPUS; i++ {
		go findWords(&trie, outerWheelChan, statUpdates, maxWheelChan)
	}

	go func() {
		combinationRepetition(outerWheelChan)
		close(outerWheelChan)
	}()

	var maxWheel Wheel
	var maxWordCount int
	for i := 0; i < NUM_CPUS; i++ {
		select {
		case pair := <-maxWheelChan:
			if pair.wordCount > maxWordCount {
				maxWheel = pair.wheel
				maxWordCount = pair.wordCount
			}
		}
	}

	finished <- true

	close(finished)
	close(maxWheelChan)
	close(statUpdates)

	return maxWheel, maxWordCount
}

func findWordsForWheelClarification(wheel Wheel, words []int_tree.WordDetails) int{
	count := 0

	for _, word := range words {
		if canWordBeSpeltFromWheel(word.SortedLetterCounts, wheel) {
			count++
		}
	}

	return count
}

var testMode = true

func main() {
	var trie int_tree.Node
	var words []int_tree.WordDetails

	if testMode {
		trie, words = int_tree.CreateIntDictionaryTree("./first_1000-3-to-9-letter-words.txt")
	} else {
		trie, words = int_tree.CreateIntDictionaryTree("./3-to-9-letter-words.txt")
	}

	wheel, wordCount := findBestLetterWheel(trie, words)
	clarificationWordCount := findWordsForWheelClarification(wheel, words)

	printOutput(wheel, wordCount)
	fmt.Printf("Clarification count found %d words\n", clarificationWordCount)

	func(cpuProfile string) {
		flag.Parse()
		if cpuProfile != "" {
			f, err := os.Create(cpuProfile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			defer f.Close()
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
		}
	}("")
	//}("cpu.prof")
}
