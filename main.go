package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

type Bigram struct {
	First  byte
	Second byte
}

func main() {
	names, err := readFile("names.txt")
	if err != nil {
		log.Fatal(err)
	}

	bigrams := make(map[Bigram]int)
	for _, name := range names {
		for i := 0; i < len(name)-1; i++ {
			bigram := Bigram{name[i], name[i+1]}
			bigrams[bigram]++
		}
	}

	printBigramTable(bigrams)

	newName := generateName(bigrams)
	fmt.Println("New name:", newName)
}

func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names = append(names, name)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func printBigramTable(bigrams map[Bigram]int) {
	fmt.Println("Bigram\tFrequency")
	for bigram, freq := range bigrams {
		fmt.Printf("%c%c\t%d\n", bigram.First, bigram.Second, freq)
	}
}

func generateName(bigrams map[Bigram]int) string {

	firstLetter := randomLetter(bigrams, ' ')
	name := string(firstLetter)

	for lastLetter := byte(' '); lastLetter != '$'; {

		bigram := Bigram{lastLetter, firstLetter}
		nextLetter := randomLetter(bigrams, bigram.Second)
		name += string(nextLetter)

		lastLetter = firstLetter
		firstLetter = nextLetter
	}

	return strings.TrimSpace(name)
}

func randomLetter(bigrams map[Bigram]int, prevLetter byte) byte {
	var letters []byte
	var frequencies []int
	totalFreq := 0

	for bigram, freq := range bigrams {
		if bigram.First == prevLetter {
			letters = append(letters, bigram.Second)
			frequencies = append(frequencies, freq)
			totalFreq += freq
		}
	}

	if totalFreq <= 0 {
		return ' '
	}

	r := rand.Intn(totalFreq)
	for i, freq := range frequencies {
		r -= freq
		if r < 0 {
			return letters[i]
		}
	}

	return ' '
}
