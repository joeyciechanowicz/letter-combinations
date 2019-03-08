#!/usr/bin/env node
const readline = require('readline')
const fs = require('fs');

/*
 * This script calculates what group of letters can spell the most unique words, using each letter once.
 *
 * i.e. e,n,h,t can spell then, hen, he, ent
 */

const tree = {
	pathTotal: 0,
	children: {},
	words: [],
	currentPath: ''
};

let currentMax = 0;
let currentMaxPointer = null;
let totalWords = 0;

function addWord(word) {
	const sortedLetters = word.split('').sort();

	let head = tree;
	for (let j = 0; j < sortedLetters.length; j++) {
		const letter = sortedLetters[j];

		if (!head.children[letter]) {
			head.children[letter] = {
				pathTotal: 0,
				children: {},
				words: [],
				currentPath: sortedLetters.slice(0, j + 1)
			};
		}

		head.children[letter].pathTotal++;
		head = head.children[letter];
	}

	head.words.push(word);
	if (head.words.length > currentMax) {
		currentMax = head.words.length;
		currentMaxPointer = head;
	}
	totalWords++;
}

function processWordFile(wordFile) {
	const lineReader = readline.createInterface({
		input: fs.createReadStream(wordFile)
	});

	lineReader.on('line', function (word) {
		addWord(word.toLowerCase());
	});

	lineReader.on('close', function () {
		console.log(`Total words ${totalWords.toLocaleString()}\n`);
		console.log(`Highest combination of words come from the letters ${currentMaxPointer.currentPath.join(', ')}`);
		console.log(`Words (${currentMaxPointer.words.length}): ${currentMaxPointer.words.join(', ')}`);
	});
}

processWordFile('words_alpha.txt');
// processWordFile('words.txt');
