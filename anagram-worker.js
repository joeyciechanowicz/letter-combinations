const { workerData, parentPort } = require('worker_threads')


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
			return false;
		}
	}

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

	searchSet(currentWord.sortedSet, 0, workerData.trie, currentWord, wordList);

	return wordList;
}

let maxWord;
let maxWordsList = [];

let ticks = 0;
const tickSize = 100;

for (let i = workerData.start; i < workerData.end; i++) {
	const currentWord = workerData.words[i];
	if (currentWord.canSkip) {
		continue;
	}

	if (++ticks % tickSize === 0) {
		parentPort.postMessage({tick: tickSize});
	}

	const anagrams = findAnagrams(currentWord);

	if (anagrams.length > maxWordsList.length) {
		maxWord = workerData.words[i].word;
		maxWordsList = anagrams;
	}
}

const ticksLeft = (workerData.end - workerData.start - ticks) % 100;
parentPort.postMessage({tick: ticksLeft});

parentPort.postMessage({maxWord, anagrams: maxWordsList.map(x => x.word)});
