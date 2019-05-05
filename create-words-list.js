const fs = require('fs');
const readline = require('readline');

const lineReader = readline.createInterface({
    input: fs.createReadStream('./words_no-names-or-places.txt')
});

const words = [];

lineReader.on('line', function (line) {
    if (line.length >= 3 && line.length <= 9) {
        words.push(line);
    }
});

lineReader.on('close', () => {
    fs.writeFileSync('3-to-9-letter-words.txt', words.join('\n'));
});

