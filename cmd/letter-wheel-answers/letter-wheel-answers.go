package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

//const filename = "./first_1000-3-to-9-letter-words.txt"
const filename = "./3-to-9-letter-words.txt"

func shouldWordBeIncluded(word string, mainLetter string, letterCounts map[string]int) bool {
	lettersSeen := make(map[string]int)

	for _, letter := range strings.Split(word, "") {
		lettersSeen[letter] += 1

		if lettersSeen[letter] > letterCounts[letter] {
			return false
		}
	}

	return lettersSeen[mainLetter] >= 1
}

func parseArgs() (string, map[string]int) {
	args := os.Args[1:]

	if len (args) != 9 {
		panic("Incorrect number of arguments. Expected 9 letters")
	}

	letterCounts := make(map[string]int)
	for _, letter := range args {
		letterCounts[letter] += 1
	}

	return args[0], letterCounts
}

func printOutput(mainLetter string, letterCounts map[string]int, words []string, time time.Duration) {
	var letters []string

	for letter, count := range letterCounts {
		if letter == mainLetter {
			if count > 1 {
				for i := 0; i < count - 1; i++ {
					letters = append(letters, letter)
				}
			}
			continue
		}

		for i := 0; i < count; i++ {
			letters = append(letters, letter)
		}
	}

	fmt.Printf("┏━━━━━━━━━━━┓\n")
	fmt.Printf("┃ %s   %s   %s ┃\n", letters[0], letters[1], letters[2])
	fmt.Printf("┃   ┏━━━┓   ┃\n")
	fmt.Printf("┃ %s ┃ %s ┃ %s ┃  Found %d words in %dms\n", letters[3], mainLetter, letters[4], len(words), time.Nanoseconds() / 1e6)
	fmt.Printf("┃   ┗━━━┛   ┃\n")
	fmt.Printf("┃ %s   %s   %s ┃\n", letters[5], letters[6], letters[7])
	fmt.Printf("┗━━━━━━━━━━━┛\n")

	nineLetterString := "Nine letter words: "
	nineLetterCount := 0
	for _, word := range words {
		if len(word) == 9 {
			if nineLetterCount > 0 {
				nineLetterString += ", "
			}
			nineLetterString += word
			nineLetterCount++
		}
	}
	if nineLetterCount > 0 {
		fmt.Println(nineLetterString)
	}

	eightLetterString := "Eight letter words: "
	eightLetterCount := 0
	for _, word := range words {
		if len(word) == 8 {
			if eightLetterCount > 0 {
				eightLetterString += ", "
			}
			eightLetterString += word
			eightLetterCount++
		}
	}
	if eightLetterCount > 0 {
		fmt.Println(eightLetterString)
	}

	fmt.Print("Words: ")

	for i, word := range words {
		fmt.Print(word)
		if i != len(words) - 1 {
			fmt.Print(", ")
		}
	}

	fmt.Println()
}

func main() {
	start := time.Now()
	var words []string
	mainLetter, letterCounts := parseArgs()

	fileHandle, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer closeIO(fileHandle)

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		word := scanner.Text()

		if shouldWordBeIncluded(word, mainLetter, letterCounts) {
			words = append(words, word)
		}
	}

	printOutput(mainLetter, letterCounts, words, time.Since(start))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func closeIO(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}