#!/usr/bin/env node
const readline = require('readline')
const fs = require('fs');

const toDashedLayers = require('./lib/to-dashed-layers');

/*
 * This script calculates which word contains the most other words within it, allowing for some letters to be skipped
 * i.e. then can spell the, he, hen & ten (notice that ten skips one letter)
 */

let count = 0;
const words = new Array(370100);

class WordPair {
	constructor(word, sortedWord) {
		this.word = word;
		this.sortedWord = sortedWord;

		this.letterCounts = sortedWord.split('').reduce((counts, curr) => {
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
	words[count] = new WordPair(word, sortedWord);
	count++;
}

function checkWordCanBeSpeltFrom(mainWord, wordToCheck) {
	for (let j = 0; j < mainWord.length; j++) {
		if (!(wordToCheck.letterCounts[mainWord[j][0]] <= mainWord[j][1])) {
			return false;
		}
	}
	return true;
}

function findWords(wordPair) {
	const checkAgainstCounts = Object.entries(wordPair.letterCounts);
	const wordList = [];

	for (let i = 0; i < words.length; i++) {
		const word = words[i];
		if (checkWordCanBeSpeltFrom(checkAgainstCounts, word)) {
			wordList.push(word.word);
		}
	}

	return wordList;
}

function findMostWords() {
	let maxWord;
	let maxWordsList = [];

	for (let i = 0; i < words.length; i++) {
		console.time(`Iteration ${i}`);
		const spellableWords = findWords(words[i])

		if (spellableWords.length > maxWordsList.length) {
			maxWord = words[i].word;
			maxWordsList = spellableWords;
		}

		console.timeEnd(`Iteration ${i}`);
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
