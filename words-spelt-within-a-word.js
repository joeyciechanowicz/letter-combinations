#!/usr/bin/env node
const readline = require('readline');
const fs = require('fs');
const ProgressBar = require('progress');
const Combinatorics = require('js-combinatorics');


/*
 * This script calculates which word contains the most other words within it, allowing for some letters to be skipped
 * i.e. then can spell the, he, hen & ten (notice that ten skips one letter)
 */

const words = [];

const tri = {
	pathTotal: 0,
	children: {},
	words: [],
	currentPath: ''
};

class WordDetails {
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

		this.sortedSet = sortedLetters.reduce((set, curr) => {
			if (set[set.length - 1] !== curr) {
				set.push(curr);
			}
			return set;
		}, [sortedLetters[0]]);
	}
}

function addWord(word, sortedWord) {
	const sortedLetters = sortedWord.split('');
	const wordDetails = new WordDetails(word, sortedLetters);
	words.push(wordDetails);

	let head = tri;
	for (let j = 0; j < wordDetails.sortedSet.length; j++) {
		const letter = wordDetails.sortedSet[j];

		if (!head.children[letter]) {
			head.children[letter] = {
				children: {},
				words: []
			};
		}

		head = head.children[letter];
	}

	head.words.push(wordDetails);
}

function getCombinations(array, size, start, initialStuff, wordList, wordDetails) {
	if (initialStuff.length >= size) {
		getWordsForCombination(wordDetails, initialStuff).forEach(x => wordList.push(x));
	} else {
		for (let i = start; i < array.length; ++i) {
			getCombinations(array, size, i + 1, initialStuff.concat(array[i]), wordList, wordDetails);
		}
	}
}

/**
 *
 * @param checkAgainstCounts
 * @param wordToCheck
 * @returns {boolean}
 */
function checkWordCanBeSpeltFrom(mainWord, wordToCheck, combination) {
	for (let char in wordToCheck.letterCounts) {
		if (wordToCheck.letterCounts[char] > mainWord.letterCounts[char]) {
			return false;
		}
	}
	return true;
}

function getWordsForCombination(mainWord, combination) {
	let head = tri;
	for (let j = 0; j < combination.length; j++) {
		const letter = combination[j];

		if (!head.children[letter]) {
			return [];
		}

		head = head.children[letter];
	}

	return head.words.filter(word => checkWordCanBeSpeltFrom(mainWord, word, combination) === true);
}

/**
 *
 * @param {WordDetails} wordDetails
 * @returns {Array}
 */
function findWords(wordDetails) {
	let wordList = [];

	for (let i = 1; i < wordDetails.sortedSet.length; i++) {
		getCombinations(wordDetails.sortedSet, i, 0, [], wordList, wordDetails);
	}

	return wordList;
}

function findMostWords() {
	let maxWord;
	let maxWordsList = [];

	const bar = new ProgressBar('[:bar] :rate/wps :percent', {total: words.length});

	for (let i = 0; i < words.length; i++) {
		bar.tick();

		const spellableWords = findWords(words[i]);

		if (spellableWords.length > maxWordsList.length) {
			maxWord = words[i].word;
			maxWordsList = spellableWords;
		}
	}

	console.log({
		maxWord,
		count: maxWordsList.length,
		words: maxWordsList.map(x => x.word).join(', ')
	});
}

let currentWord;
let haveWord = false;
const lineReader = readline.createInterface({
	input: fs.createReadStream('words_with_sorted_word.txt')
	// input: fs.createReadStream('smaller_words_with_sorted_word.txt')
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
