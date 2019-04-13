package dictionary_tree

import (
	"github.com/joeyciechanowicz/letter-combinations/pkg/reader"
	"sort"
)

type RuneCount struct {
	Letter rune
	Count byte
}

type WordDetails struct {
	Word             string
	SortedRuneCounts []RuneCount
}

type WordDetailsSlice []WordDetails

type Node struct {
	Children map[rune]*Node
	Words    []*WordDetails
}

type runeSlice []rune

func (p runeSlice) Len() int           { return len(p) }
func (p runeSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p runeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func newWordDetails(word string) WordDetails {
	var details = WordDetails{
		word,
		[]RuneCount{},
	}

	// Construct a set of letter and their associated counts
	// afterwards we finish constructing our rune-count array
	sortedLetters := []rune(word)
	sort.Sort(runeSlice(sortedLetters))

	letterCounts := make(map[rune]byte)

	details.SortedRuneCounts = append(details.SortedRuneCounts, RuneCount{sortedLetters[0], 0})

	for _, char := range sortedLetters {
		letterCounts[char] += 1

		if details.SortedRuneCounts[len(details.SortedRuneCounts)-1].Letter != char {
			details.SortedRuneCounts = append(details.SortedRuneCounts, RuneCount{char, 0})
		}
	}

	for i, runeCount := range details.SortedRuneCounts {
		details.SortedRuneCounts[i].Count = letterCounts[runeCount.Letter]
	}

	return details
}


func CreateDictionaryTree(filename string) (Node, []WordDetails){
	nodeCount := 0

	var words []WordDetails
	var trie = Node{
		make(map[rune]*Node),
		make([]*WordDetails, 0),
	}

	reader.ReadFile(filename, func(word string) {
		var details WordDetails
		details = newWordDetails(word)

		var head *Node
		head = &trie

		for _, runeCount := range details.SortedRuneCounts {
			if _, ok := head.Children[runeCount.Letter]; !ok {
				nodeCount++

				head.Children[runeCount.Letter] = &Node{
					make(map[rune]*Node),
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
