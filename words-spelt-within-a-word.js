#!/usr/bin/env node
const readline = require('readline');
const fs = require('fs');
const ProgressBar = require('progress');


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

/**
 *
 * @param mainWord
 * @param wordToCheck
 * @returns {boolean}
 */
function checkWordCanBeSpeltFrom(mainWord, wordToCheck) {
	for (let char in wordToCheck.letterCounts) {
		// noinspection JSUnfilteredForInLoop
		if (wordToCheck.letterCounts[char] > mainWord.letterCounts[char]) {
			return false;
		}
	}
	return true;
}

function addWordsForCombination(mainWord, combination, wordList) {
	let head = tri;
	for (let j = 0; j < combination.length; j++) {
		const letter = combination[j];

		if (!head.children[letter]) {
			return [];
		}

		head = head.children[letter];
	}

	for (let i = 0; i < head.words.length; i++) {
		if (checkWordCanBeSpeltFrom(mainWord, head.words[i])) {
			wordList.push(head.words[i]);
		}
	}
}


function getCombinations(array, size, start, initialStuff, combinations) {
	if (initialStuff.length >= size) {
		combinations.push(initialStuff);
	} else {
		for (let i = start; i < array.length; ++i) {
			getCombinations(array, size, i + 1, initialStuff.concat(array[i]), combinations);
		}
	}
}

function getAllCombinations(sortedSet) {
	const combinations = [];
	for (let i = 1; i < sortedSet.length; i++) {
		getCombinations(sortedSet, i, 0, [], combinations);
	}
	return combinations;
}

/**
 *
 * @param {WordDetails} wordDetails
 * @returns {Array}
 */
function findWords(wordDetails) {
	let wordList = [];

	const combinations = getAllCombinations(wordDetails.sortedSet);

	for (let i = 0; i < combinations.length; i++) {
		addWordsForCombination(wordDetails, combinations[i], wordList);
	}

	return wordList;
}

function findMostWords() {
	let maxWord;
	let maxWordsList = [];

	const bar = new ProgressBar('[:bar] :rate/wps :percent', {total: words.length});

	console.time('Time to find words');

	let skipped = 0;

	for (let i = 0; i < words.length; i++) {
		bar.tick();

		const word = words[i];
		if (word.canSkip) {
			skipped++;
			continue;
		}

		const spellableWords = findWords(word);

		if (spellableWords.length > maxWordsList.length) {
			maxWord = words[i].word;
			maxWordsList = spellableWords;
		}
	}

	console.timeEnd('Time to find words');

	console.log(`
		Longest word: ${maxWord} with ${maxWordsList.length} words. (Skipped ${skipped})
		${maxWordsList.map(x => x.word).join(', ')}
	`);
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

	if (Object.keys(head.children).length > 0) {
		wordDetails.canSkip = true;
	}
	head.words.push(wordDetails);
}

let currentWord;
let haveWord = false;
const lineReader = readline.createInterface({
	// input: fs.createReadStream('reversed_words_with_sorted_word.txt')
	input: fs.createReadStream('reversed_smaller_words_with_sorted_word.txt')
});

console.time('Time to add');
lineReader.on('line', (word) => {
	if (haveWord) {
		addWord(word, currentWord);
		haveWord = false;
	} else {
		currentWord = word;
		haveWord = true;
	}
});

lineReader.on('close', function () {
	console.timeEnd('Time to add');
	console.log(process.memoryUsage().heapUsed / 1024 / 1024 + ' mb');

	setTimeout(findMostWords, 0);
});
