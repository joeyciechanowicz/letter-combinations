package int_tree

import (
	"github.com/joeyciechanowicz/letter-combinations/pkg/reader"
	"sort"
)

func ToAlphabetIndex(letter rune) int {
	return int(letter) - int(ToRune("a"))
}

func ToRune(letter string) rune {
	return []rune(letter)[0]
}

var RuneToLetters = map[rune]int {
	ToRune("a"): ToAlphabetIndex(ToRune("a")),
	ToRune("b"): ToAlphabetIndex(ToRune("b")),
	ToRune("c"): ToAlphabetIndex(ToRune("c")),
	ToRune("d"): ToAlphabetIndex(ToRune("d")),
	ToRune("e"): ToAlphabetIndex(ToRune("e")),
	ToRune("f"): ToAlphabetIndex(ToRune("f")),
	ToRune("g"): ToAlphabetIndex(ToRune("g")),
	ToRune("h"): ToAlphabetIndex(ToRune("h")),
	ToRune("i"): ToAlphabetIndex(ToRune("i")),
	ToRune("j"): ToAlphabetIndex(ToRune("j")),
	ToRune("k"): ToAlphabetIndex(ToRune("k")),
	ToRune("l"): ToAlphabetIndex(ToRune("l")),
	ToRune("m"): ToAlphabetIndex(ToRune("m")),
	ToRune("n"): ToAlphabetIndex(ToRune("n")),
	ToRune("o"): ToAlphabetIndex(ToRune("o")),
	ToRune("p"): ToAlphabetIndex(ToRune("p")),
	ToRune("q"): ToAlphabetIndex(ToRune("q")),
	ToRune("r"): ToAlphabetIndex(ToRune("r")),
	ToRune("s"): ToAlphabetIndex(ToRune("s")),
	ToRune("t"): ToAlphabetIndex(ToRune("t")),
	ToRune("u"): ToAlphabetIndex(ToRune("u")),
	ToRune("v"): ToAlphabetIndex(ToRune("v")),
	ToRune("w"): ToAlphabetIndex(ToRune("w")),
	ToRune("x"): ToAlphabetIndex(ToRune("x")),
	ToRune("y"): ToAlphabetIndex(ToRune("y")),
	ToRune("z"): ToAlphabetIndex(ToRune("z")),
}

var Alphabet = [26]rune {
	ToRune("a"),
	ToRune("b"),
	ToRune("c"),
	ToRune("d"),
	ToRune("e"),
	ToRune("f"),
	ToRune("g"),
	ToRune("h"),
	ToRune("i"),
	ToRune("j"),
	ToRune("k"),
	ToRune("l"),
	ToRune("m"),
	ToRune("n"),
	ToRune("o"),
	ToRune("p"),
	ToRune("q"),
	ToRune("r"),
	ToRune("s"),
	ToRune("t"),
	ToRune("u"),
	ToRune("v"),
	ToRune("w"),
	ToRune("x"),
	ToRune("y"),
	ToRune("z"),
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

	details.SortedLetterCounts = append(details.SortedLetterCounts, LetterCount{ToAlphabetIndex(sortedLetters[0]), 0})

	for _, char := range sortedLetters {
		alphaIndex := ToAlphabetIndex(char)
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
