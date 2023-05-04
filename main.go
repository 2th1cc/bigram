package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	// read data from file
	data, err := ioutil.ReadFile("names.txt")
	if err != nil {
		panic(err)
	}

	// convert data to lowercase string
	text := strings.ToLower(string(data))

	// calculate bigram probabilities
	bigrams := make(map[string]int)
	for i := 0; i < len(text)-1; i++ {
		bigram := text[i : i+2]
		bigrams[bigram]++
	}

	// calculate total number of bigrams
	total := 0
	for _, count := range bigrams {
		total += count
	}

	// calculate probabilities
	probabilities := make(map[string]float64)
	for bigram, count := range bigrams {
		probabilities[bigram] = float64(count) / float64(total)
	}

	// print table of bigram probabilities
	fmt.Println("Bigram probabilities:")
	for bigram, probability := range probabilities {
		fmt.Printf("%s: %.4f\n", bigram, probability)
	}

	// generate a name
	rand.Seed(time.Now().UnixNano())
	name := ""
	lastLetter := "^"
	for {
		// get possible next bigrams
		var candidates []string
		for bigram := range probabilities {
			if strings.HasPrefix(bigram, lastLetter) {
				candidates = append(candidates, bigram)
			}
		}

		// break out of loop if there are no candidates
		if len(candidates) == 0 {
			break
		}

		// choose next bigram randomly based on probabilities
		sum := 0.0
		for _, bigram := range candidates {
			sum += probabilities[bigram]
		}
		randVal := rand.Float64() * sum
		nextBigram := ""
		for _, bigram := range candidates {
			randVal -= probabilities[bigram]
			if randVal <= 0 {
				nextBigram = bigram
				break
			}
		}

		// add next letter to name
		nextLetter := string(nextBigram[1])
		if nextLetter == "$" {
			break
		}
		name += nextLetter
		lastLetter = nextLetter
	}

	// capitalize first letter of name
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	if name != "" {
		fmt.Println("Generated name:", name)
	} else {
		fmt.Println("Unable to generate name")
	}
}
