package main

import (
	"flag"
	"github.com/joeyciechanowicz/letter-combinations/pkg/dictionary-tree"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)


var cpuprofile = "cpu.prof"
//var cpuprofile = ""

//var memprofile = "mem.prof"
var memprofile = ""

func main() {
	//var trie, words = dictionary_tree.CreateDictionaryTree("./words_alpha.txt")
	//var trie, words = dictionary_tree.CreateDictionaryTree("./words_no-names-or-places.txt")
	//var trie, words = dictionary_tree.CreateDictionaryTree("./first_2000_words.txt")

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
