package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	file, err := os.Open("names.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		data = append(data, buf[:n]...)
	}
	names := strings.Split(string(data), "\n")

	bigramFreq := make(map[string]int)

	for _, name := range names {
		bigrams := getBigrams(name)
		for _, bigram := range bigrams {
			bigramFreq[bigram]++
		}
	}

	totalBigrams := len(getBigrams(strings.Join(names, "")))
	bigramProb := make(map[string]float64)
	for bigram, freq := range bigramFreq {
		bigramProb[bigram] = float64(freq) / float64(totalBigrams)
	}

	name := generateName(bigramProb)
	fmt.Println(name)

	displayTable(bigramProb)

}

func getBigrams(name string) []string {
	var bigrams []string
	for i := 0; i < len(name)-1; i++ {
		bigrams = append(bigrams, string(name[i])+string(name[i+1]))
	}
	return bigrams
}

func getFirstLetter(bigramProb map[string]float64) string {
	rand.Seed(time.Now().UnixNano())
	var letters []string
	var probs []float64
	for bigram, prob := range bigramProb {
		if strings.HasPrefix(bigram, "^") {
			letters = append(letters, string(bigram[1]))
			probs = append(probs, prob)
		}
	}
	return letters[randomChoice(probs)]
}

func randomChoice(probs []float64) int {
	rand.Seed(time.Now().UnixNano())
	sumProb := 0.0
	for _, prob := range probs {
		sumProb += prob
	}
	r := rand.Float64() * sumProb
	for i, prob := range probs {
		r -= prob
		if r <= 0 {
			return i
		}
	}
	return len(probs) - 1
}

func generateName(bigramProb map[string]float64) string {
	name := getFirstLetter(bigramProb)
	for {
		var nextLetter string
		var probs []float64
		lastBigram := string(name[len(name)-1]) + "$"
		for bigram, prob := range bigramProb {
			if strings.HasPrefix(bigram, lastBigram) {
				nextLetter = string(bigram[1])
				break
			}
			if strings.HasPrefix(bigram, string(name[len(name)]-1)) {
				probs = append(probs, prob)
			}
		}
		if nextLetter == "" {
			nextLetter = string(randomChoice(probs))
		}
		if nextLetter == "$" {
			break
		}
		name += nextLetter
	}
	return name
}

func displayTable(bigramProb map[string]float64) {
	fmt.Printf("%-5s%-5s%s\n", "Prev", "Curr", "Probability")
	for bigram, prob := range bigramProb {
		fmt.Printf("%-5s%-5s%.4f\n", string(bigram[0]), string(bigram[1]), prob)
	}
}
