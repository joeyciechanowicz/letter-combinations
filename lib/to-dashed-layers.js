const alphabet = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'];


/**
 * eehl ->
 * [
 * 		'----e',
 * 		'----e--h---l'
 * 	]
 *
 *
 * @param sortedLetters (char[])
 */
function toDashedLayers(sortedLetters) {
	let alphabetIndex = 0;
	let layers = [];

	let currentString = '';

	for (let i = 0; i < sortedLetters.length; i++) {
		const currentChar = sortedLetters[i];

		// add dashes to our current string until we hit the right letter
		while (alphabet[alphabetIndex] !== currentChar) {
			currentString += '-';
			alphabetIndex++;
		}

		alphabetIndex++;

		currentString += currentChar;

		if (i < sortedLetters.length - 1) {
			if (sortedLetters[i + 1] === currentChar) {
				// we have two successive characters, need a new layer
				layers.push(currentString);
				currentString = '';
				alphabetIndex = 0;
			}
		}
	}

	layers.push(currentString);

	return layers;
}

module.exports = toDashedLayers;

