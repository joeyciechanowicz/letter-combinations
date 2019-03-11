#!/usr/bin/env node
const readline = require('readline');
const fs = require('fs');
const ProgressBar = require('progress');


/*
 * This script calculates which word contains the most other words within it, allowing for some letters to be skipped
 * i.e. then can spell the, he, hen & ten (notice that ten skips one letter)
 */

const words = [];

const trie = {
	pathTotal: 0,
	children: {},
	words: [],
	currentPath: ''
};

let trieNodesCount = 0;
let incorrectChecks = 0;
let correctChecks = 0;

class WordDetails {
	constructor(word) {
		this.word = word;
		const sortedLetters = word.split('').sort();

		this.letterCounts = {};
		this.sortedSet = [sortedLetters[0]];

		for (let i = 0; i < sortedLetters.length; i++) {
			const curr = sortedLetters[i];

			if (this.letterCounts[curr]) {
				this.letterCounts[curr]++;
			} else {
				this.letterCounts[curr] = 1;
			}

			if (this.sortedSet[this.sortedSet.length - 1] !== curr) {
				this.sortedSet.push(curr);
			}
		}
	}
}

/**
 *
 * @param mainWord
 * @param wordToCheck
 * @returns {boolean}
 */
function isWordAnagram(mainWord, wordToCheck) {
	if (wordToCheck.word.length > mainWord.word.length) {
		return false;
	}

	for (let i = 0; i < wordToCheck.sortedSet.length; i++) {
		const char = wordToCheck.sortedSet[i];
		if (wordToCheck.letterCounts[char] > mainWord.letterCounts[char]) {
			incorrectChecks++;
			return false;
		}
	}

	correctChecks++;
	return true;
}

function searchSet(letters, start, head, currentWord, wordList) {
	if (head.words.length > 0) {
		for (let i = 0; i < head.words.length; i++) {
			if (isWordAnagram(currentWord, head.words[i])) {
				wordList.push(head.words[i]);
			}
		}
	}

	for (let i = start; i < letters.length; i++) {
		if (head.children[letters[i]]) {
			searchSet(letters, i + 1, head.children[letters[i]], currentWord, wordList);
		}
	}
}

/**
 *
 * @param {WordDetails} currentWord
 * @returns {Array}
 */
function findAnagrams(currentWord) {
	let wordList = [];

	searchSet(currentWord.sortedSet, 0, trie, currentWord, wordList);

	return wordList;
}

function findMostWords() {
	let maxWord;
	let maxWordsList = [];

	const bar = new ProgressBar('[:bar] :rate/wps :percent m=:average', {total: words.length});

	console.time('Time to find words');

	let skipped = 0;

	for (let i = 0; i < words.length; i++) {
		bar.tick(1, {});

		const currentWord = words[i];
		if (currentWord.canSkip) {
			skipped++;
			continue;
		}

		const anagrams = findAnagrams(currentWord);

		if (anagrams.length > maxWordsList.length) {
			maxWord = words[i].word;
			maxWordsList = anagrams;
		}
	}

	console.timeEnd('Time to find words');

	console.log(`Longest word: ${maxWord} with ${maxWordsList.length} words. Incorrect checks: ${incorrectChecks}. Correct checks: ${correctChecks}`);

	if (maxWordsList.length > 100) {
		console.log('Words have been written to anagrams.txt');
		fs.writeFileSync('anagrams.txt', maxWordsList.map(x => x.word).join(', '));
	}
}

function addWord(word) {
	const wordDetails = new WordDetails(word);
	words.push(wordDetails);

	let head = trie;
	for (let j = 0; j < wordDetails.sortedSet.length; j++) {
		const letter = wordDetails.sortedSet[j];

		if (!head.children[letter]) {
			trieNodesCount++;
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

const lineReader = readline.createInterface({
	input: fs.createReadStream('2000_random_words_alpha.txt')
	// input: fs.createReadStream('first_2000_words.txt')
	// input: fs.createReadStream('words_alpha.txt')
});

console.time('Time to add');
lineReader.on('line', (word) => {
	addWord(word);
});

lineReader.on('close', function () {
	console.timeEnd('Time to add');
	console.log(`${(process.memoryUsage().heapUsed / 1024 / 1024).toFixed(1)} mb. ${trieNodesCount.toLocaleString()} nodes in trie`);

	setImmediate(findMostWords);
});
