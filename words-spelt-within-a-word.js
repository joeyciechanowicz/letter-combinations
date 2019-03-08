#!/usr/bin/env node
const readline = require('readline')
const fs = require('fs');
const ProgressBar = require('progress');

/*
 * This script calculates which word contains the most other words within it, allowing for some letters to be skipped
 * i.e. then can spell the, he, hen & ten (notice that ten skips one letter)
 */

let count = 0;
const words = new Array(370100);

const wordsBasedOnStartingLetter = {
	'a': [],
	'b': [],
	'c': [],
	'd': [],
	'e': [],
	'f': [],
	'g': [],
	'h': [],
	'i': [],
	'j': [],
	'k': [],
	'l': [],
	'm': [],
	'n': [],
	'o': [],
	'p': [],
	'q': [],
	'r': [],
	's': [],
	't': [],
	'u': [],
	'v': [],
	'w': [],
	'x': [],
	'y': [],
	'z': [],
};

class WordPair {
	constructor(word, sortedLetters) {
		this.word = word;
		this.sortedLetters = sortedLetters;

		this.letterCounts = sortedLetters.reduce((counts, curr) => {
			if (counts[curr]) {
				counts[curr]++;
			} else {
				counts[curr] = 1;
			}
			return counts;
		}, {});
	}
}

function addWord(word, sortedWord) {
	const sortedLetters = sortedWord.split('');
	words[count] = new WordPair(word, sortedLetters);
	wordsBasedOnStartingLetter[sortedLetters[0]].push(words[count]);
	count++;
}

function checkWordCanBeSpeltFrom(checkAgainstCounts, wordToCheck) {
	for (let j = 0; j < checkAgainstCounts.length; j++) {
		if (!(wordToCheck.letterCounts[checkAgainstCounts[j][0]] <= checkAgainstCounts[j][1])) {
			return false;
		}
	}
	return true;
}

function findWords(wordPair) {
	const checkAgainstCounts = Object.entries(wordPair.letterCounts);
	const wordList = [];

	for (let j = 0; j < checkAgainstCounts.length; j++) {
		for (let i = 0; i < wordsBasedOnStartingLetter[checkAgainstCounts[j][0]].length; i++) {
			const word = wordsBasedOnStartingLetter[checkAgainstCounts[j][0]][i];
			if (checkWordCanBeSpeltFrom(checkAgainstCounts, word)) {
				wordList.push(word.word);
			}
		}
	}

	return wordList;
}

function findMostWords() {
	let maxWord;
	let maxWordsList = [];
	const bar = new ProgressBar('[:bar] :rate/wps :percent :etas', {total: words.length});

	for (let i = 0; i < words.length; i++) {
		bar.tick();

		const spellableWords = findWords(words[i])

		if (spellableWords.length > maxWordsList.length) {
			maxWord = words[i].word;
			maxWordsList = spellableWords;
		}
	}

	return {
		maxWord,
		maxWordsList
	};
}

let currentWord;
let haveWord = false;
const lineReader = readline.createInterface({
	input: fs.createReadStream('words_with_sorted_word.txt')
});

console.time('Time to add');
lineReader.on('line', (word) => {
	if (haveWord) {
		addWord(currentWord, word);
		haveWord = false;
	} else {
		currentWord = word;
		haveWord = true;
	}
});

lineReader.on('close', function () {
	console.timeEnd('Time to add');

	setImmediate(findMostWords);
});
