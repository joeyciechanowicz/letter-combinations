package dictionary_tree

import (
	"github.com/joeyciechanowicz/letter-combinations/pkg/reader"
	"sort"
)

type WordDetails struct {
	Word         string
	SortedSet    []rune
	LetterCounts map[rune]byte
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
		[]rune{},
		make(map[rune]byte),
	}

	sortedLetters := []rune(word)
	sort.Sort(runeSlice(sortedLetters))

	details.SortedSet = append(details.SortedSet, sortedLetters[0])

	for _, char := range sortedLetters {
		details.LetterCounts[char] += 1

		if details.SortedSet[len(details.SortedSet)-1] != char {
			details.SortedSet = append(details.SortedSet, char)
		}
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

		for _, letter := range details.SortedSet {
			if _, ok := head.Children[letter]; !ok {
				nodeCount++

				head.Children[letter] = &Node{
					make(map[rune]*Node),
					[]*WordDetails{},
				}
			}

			head = head.Children[letter]
		}

		words = append(words, details)
		head.Words = append(head.Words, &details)
	})

	//fmt.Println("Trie nodes: ", nodeCount)
	return trie, words
}
