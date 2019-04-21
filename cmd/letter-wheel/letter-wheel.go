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
	wheel     [WHEEL_SIZE]int
	wordCount int
}

func canWordBeSpeltFromWheel(word []int_tree.LetterCount, wheel []int_tree.LetterCount) bool {
	// Iterate the main words runes, setting an index for the other word
	i := -1
	for j := 0; j < len(word); j++ {
		letterAndCount := word[j]

		// move the other words index along until we find a letter that matches
		// returning false if we reach the end or the rune counts are incorrect
		for {
			i++
			if i == len(wheel) {
				return false
			}

			if letterAndCount.Letter == wheel[i].Letter {
				if letterAndCount.Count <= wheel[i].Count {
					break
				} else {
					return false
				}
			}
		}
	}

	return true
}

func searchSet(set []int, head *int_tree.Node, start int, currentWheel []int_tree.LetterCount, wordCount *int) {
	if len(head.Words) > 0 {
		for i := 0; i < len(head.Words); i++ {
			if canWordBeSpeltFromWheel(head.Words[i].SortedLetterCounts, currentWheel) {
				*wordCount += 1
			}
		}
	}

	for i := start; i < len(set); i++ {
		if _, ok := head.Children[set[i]]; ok {
			searchSet(set, head.Children[set[i]], i+1, currentWheel, wordCount)
		}
	}
}

func wheelToLetterCount(wheel [WHEEL_SIZE]int) []int_tree.LetterCount {
	var letterCounts []int_tree.LetterCount

	letterCounts = append(letterCounts, int_tree.LetterCount{
		Letter: wheel[0],
		Count:  1,
	})

	prevLetter := wheel[0]

	for i := 1; i < WHEEL_SIZE; i++ {
		currLetter := wheel[i]
		if currLetter == prevLetter {
			letterCounts[len(letterCounts)-1].Count++
		} else {
			letterCounts = append(letterCounts, int_tree.LetterCount{
				Letter: currLetter,
				Count:  1,
			})
			prevLetter = currLetter
		}
	}

	return letterCounts
}

/**
Takes wheel off a channel and finds all the words for it
 */
func findWords(trie *int_tree.Node, wheelChan <-chan [WHEEL_SIZE]int, stats chan<- bool, maxWheelChan chan<- wheelSolutions) {
	var maxWordCount int
	var maxWheel [WHEEL_SIZE]int

	for {
		currentWheel, ok := <-wheelChan
		if !ok {
			maxWheelChan <- wheelSolutions{maxWheel, maxWordCount}
			return
		}

		letterCounts := wheelToLetterCount(currentWheel)

		wordCount := 0
		searchSet(currentWheel[:], trie, 0, letterCounts, &wordCount)

		if wordCount > maxWordCount {
			maxWheel = currentWheel
			maxWordCount = wordCount
		}

		stats <- true
	}
}

func combinationRepetitionUtil(wheelChan chan<- [WHEEL_SIZE]int, chosen [WHEEL_SIZE]int, index, r, start, end int) {
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

// The main function that prints all combinations of size r
// in arr[] of size n with reactions. This function mainly
// uses CombinationRepetitionUtil()
func combinationRepetition(wheelChan chan [WHEEL_SIZE]int) {
	var chosen [WHEEL_SIZE]int
	combinationRepetitionUtil(wheelChan, chosen, 0, WHEEL_SIZE, 0, 26-1)
	//combinationRepetitionUtil(wheelChan, chosen, 0, WHEEL_SIZE, 0, 14)
}

func printOutput(lettersAsIndicies [WHEEL_SIZE]int, numFound int) {
	var letters [WHEEL_SIZE]string
	for i, letterIndicie := range lettersAsIndicies {
		letters[i] = string(int_tree.Alphabet[letterIndicie])
	}

	fmt.Printf("\n┏━━━━━━━━━━━┓\n")
	fmt.Printf("┃ %s   %s   %s ┃\n", letters[1], letters[2], letters[3])
	fmt.Printf("┃   ┏━━━┓   ┃\n")
	fmt.Printf("┃ %s ┃ %s ┃ %s ┃  Found %d\n", letters[4], letters[0], letters[5], numFound)
	fmt.Printf("┃   ┗━━━┛   ┃\n")
	fmt.Printf("┃ %s   %s   %s ┃\n", letters[6], letters[7], letters[8])
	fmt.Printf("┗━━━━━━━━━━━┛\n")
}

func findBestLetterWheel(trie int_tree.Node, details []int_tree.WordDetails) {
	const NUM_CPUS = 8

	finished := make(chan bool)
	maxWheelChan := make(chan wheelSolutions)
	wheelChan := make(chan [WHEEL_SIZE]int, NUM_CPUS)
	statUpdates := make(chan bool, NUM_CPUS)

	go stats.PrintProgress(finished, statUpdates, TOTAL_WHEELS)

	for i := 0; i < NUM_CPUS; i++ {
		go findWords(&trie, wheelChan, statUpdates, maxWheelChan)
	}

	go func() {
		combinationRepetition(wheelChan)
		close(wheelChan)
	}()

	var maxWheel [WHEEL_SIZE]int
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

	printOutput(maxWheel, maxWordCount)
}

func main() {
	trie, words := int_tree.CreateIntDictionaryTree("./3-to-9-letter-words.txt")
	//trie, words := int_tree.CreateIntDictionaryTree("./first_1000-3-to-9-letter-words.txt")

	findBestLetterWheel(trie, words)

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
