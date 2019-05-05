package int_tree

import (
	"github.com/joeyciechanowicz/letter-combinations/pkg/reader"
	"sort"
)

func toAlphabetIndex(letter rune) int {
	return int(letter) - int(toRune("a"))
}

func toRune(letter string) rune {
	return []rune(letter)[0]
}

var RuneToLetters = map[rune]int {
	toRune("a"): toAlphabetIndex(toRune("a")),
	toRune("b"): toAlphabetIndex(toRune("b")),
	toRune("c"): toAlphabetIndex(toRune("c")),
	toRune("d"): toAlphabetIndex(toRune("d")),
	toRune("e"): toAlphabetIndex(toRune("e")),
	toRune("f"): toAlphabetIndex(toRune("f")),
	toRune("g"): toAlphabetIndex(toRune("g")),
	toRune("h"): toAlphabetIndex(toRune("h")),
	toRune("i"): toAlphabetIndex(toRune("i")),
	toRune("j"): toAlphabetIndex(toRune("j")),
	toRune("k"): toAlphabetIndex(toRune("k")),
	toRune("l"): toAlphabetIndex(toRune("l")),
	toRune("m"): toAlphabetIndex(toRune("m")),
	toRune("n"): toAlphabetIndex(toRune("n")),
	toRune("o"): toAlphabetIndex(toRune("o")),
	toRune("p"): toAlphabetIndex(toRune("p")),
	toRune("q"): toAlphabetIndex(toRune("q")),
	toRune("r"): toAlphabetIndex(toRune("r")),
	toRune("s"): toAlphabetIndex(toRune("s")),
	toRune("t"): toAlphabetIndex(toRune("t")),
	toRune("u"): toAlphabetIndex(toRune("u")),
	toRune("v"): toAlphabetIndex(toRune("v")),
	toRune("w"): toAlphabetIndex(toRune("w")),
	toRune("x"): toAlphabetIndex(toRune("x")),
	toRune("y"): toAlphabetIndex(toRune("y")),
	toRune("z"): toAlphabetIndex(toRune("z")),
}

var Alphabet = [26]rune {
	toRune("a"),
	toRune("b"),
	toRune("c"),
	toRune("d"),
	toRune("e"),
	toRune("f"),
	toRune("g"),
	toRune("h"),
	toRune("i"),
	toRune("j"),
	toRune("k"),
	toRune("l"),
	toRune("m"),
	toRune("n"),
	toRune("o"),
	toRune("p"),
	toRune("q"),
	toRune("r"),
	toRune("s"),
	toRune("t"),
	toRune("u"),
	toRune("v"),
	toRune("w"),
	toRune("x"),
	toRune("y"),
	toRune("z"),
}

type LetterCount struct {
	Letter int
	Count byte
}

type WordDetails struct {
	Word             string
	SortedLetterCounts []LetterCount
}

type WordDetailsSlice []WordDetails

type Node struct {
	Children map[int]*Node
	Words    []*WordDetails
}

type runeSlice []rune
func (p runeSlice) Len() int           { return len(p) }
func (p runeSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p runeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func NewWordDetails(word string) WordDetails {
	var details = WordDetails{
		word,
		[]LetterCount{},
	}

	sortedLetters := []rune(word)
	sort.Sort(runeSlice(sortedLetters))

	letterCounts := make(map[int]byte)

	details.SortedLetterCounts = append(details.SortedLetterCounts, LetterCount{toAlphabetIndex(sortedLetters[0]), 0})

	for _, char := range sortedLetters {
		alphaIndex := toAlphabetIndex(char)
		letterCounts[alphaIndex] += 1

		if details.SortedLetterCounts[len(details.SortedLetterCounts)-1].Letter != alphaIndex {
			details.SortedLetterCounts = append(details.SortedLetterCounts, LetterCount{alphaIndex, 0})
		}
	}

	for i, runeCount := range details.SortedLetterCounts {
		details.SortedLetterCounts[i].Count = letterCounts[runeCount.Letter]
	}

	return details
}

/**
Creates a trie of WordDetails where int is used instead of runes
The int is the index of a letter in Alphabet
 */
func CreateIntDictionaryTree(filename string) (Node, []WordDetails){
	nodeCount := 0

	var words []WordDetails
	var trie = Node{
		make(map[int]*Node),
		make([]*WordDetails, 0),
	}

	reader.ReadFile(filename, func(word string) {
		var details WordDetails
		details = NewWordDetails(word)

		var head *Node
		head = &trie

		for _, runeCount := range details.SortedLetterCounts {
			if _, ok := head.Children[runeCount.Letter]; !ok {
				nodeCount++

				head.Children[runeCount.Letter] = &Node{
					make(map[int]*Node),
					[]*WordDetails{},
				}
			}

			head = head.Children[runeCount.Letter]
		}

		words = append(words, details)
		head.Words = append(head.Words, &details)
	})

	//fmt.Println("Trie nodes: ", nodeCount)
	return trie, words
}
