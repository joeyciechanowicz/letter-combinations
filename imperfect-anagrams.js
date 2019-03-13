#!/usr/bin/env node --experimental-worker
const readline = require('readline');
const fs = require('fs');
const ProgressBar = require('progress');
const os = require('os');
const { Worker } = require('worker_threads');


/*
 * This script calculates which word contains the most other words within it, allowing for some letters to be skipped
 * i.e. then can spell the, he, hen & ten (notice that ten skips one letter)
 */

const words = [];

const trie = {
	children: {},
	words: []
};

let trieNodesCount = 0;

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

function run() {
	console.time('Worker Spawning');

	let maxWord;
	let maxWordsList = [];
	let workersReceivedFrom = 0;

	const numCpus = os.cpus().length;
	const chunkSize = Math.floor(words.length / numCpus);
	const reportFrequency =  Math.floor(chunkSize / 100);

	const bar = new ProgressBar('[:bar] :rate/wps :percent', {total: words.length});

	for (let i = 0; i < numCpus; i++) {
		const start = i * chunkSize;
		const end = i < numCpus - 1 ? (i + 1) * chunkSize : words.length;

		const worker = new Worker('./anagram-worker.js', {workerData: {trie, words, start, end}});
		worker.on('message', (data) => {
			if (data.tick >= 0) {
				bar.tick(data.tick, {});
				return;
			}

			if (data.anagrams) {
				workersReceivedFrom++;
				if (data.anagrams.length > maxWordsList.length) {
					maxWord = data.maxWord;
					maxWordsList = data.anagrams;
				}

				if (workersReceivedFrom === numCpus) {
					console.log(`The word "${maxWord}" has ${maxWordsList.length} imperfect anagrams`);
					if (maxWordsList.length > 100) {
						console.log('Words have been written to anagrams.txt');
						fs.writeFileSync('anagrams.txt', maxWordsList.join(', '));
					} else {
						console.log(maxWordsList.join(','));
					}
				}

				return;
			}

			console.error(`Unhandled message from worker`, data);
		});
		worker.on('error', error => console.error(`Error from worker ${i}`, error));
		worker.on('exit', code => {
			if (code !== 0) {
				console.error(`Worker ${i} exited with code ${code}`);
			}
		});

		let onlineCount = 0;
		worker.on('online', () => {
			onlineCount++;

			if (onlineCount === numCpus) {
				console.timeEnd('Worker Spawning');
			}
		});
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
	// input: fs.createReadStream('2000_random_words_alpha.txt')
	input: fs.createReadStream('first_2000_words.txt')
	// input: fs.createReadStream('words_alpha.txt')
});

console.time('Time to add');
lineReader.on('line', (word) => {
	addWord(word);
});

lineReader.on('close', function () {
	console.timeEnd('Time to add');
	console.log(`${(process.memoryUsage().heapUsed / 1024 / 1024).toFixed(1)} mb. ${trieNodesCount.toLocaleString()} nodes in trie`);
	run();
});
