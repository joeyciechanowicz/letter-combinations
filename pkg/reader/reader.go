package reader

import (
	"bufio"
	"io"
	"log"
	"os"
)

func ReadFile(filename string, cb func(string)) {
	fileHandle, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer closeIO(fileHandle)

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		cb(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}


func closeIO(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}