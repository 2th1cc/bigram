package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func readData(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	return data, nil
}

func calculateProbabilities(data []string) map[string]float64 {
	bigramCounts := make(map[string]int)
	totalCount := 0

	for _, name := range data {
		name = "^" + name + "$" // add "^" and "$" markers to indicate the beginning and end of the name
		for i := 0; i < len(name)-1; i++ {
			bigram := name[i : i+2]
			bigramCounts[bigram]++
			totalCount++
		}
	}

	probabilities := make(map[string]float64)
	for bigram, count := range bigramCounts {
		probabilities[bigram] = float64(count) / float64(totalCount)
	}

	return probabilities
}

func selectFirstLetter(probabilities map[string]float64) string {
	rand.Seed(time.Now().UnixNano()) // initialize the random number generator
	sum := 0.0
	threshold := rand.Float64()

	for bigram, probability := range probabilities {
		if strings.HasPrefix(bigram, "^") {
			sum += probability
			if sum > threshold {
				return string(bigram[1])
			}
		}
	}

	return "" // should never happen
}

func generateName(probabilities map[string]float64) string {
	name := selectFirstLetter(probabilities)
	for {
		var lastBigram string
		if len(name) > 1 {
			lastBigram = name[len(name)-2:] + "$"
		} else {
			lastBigram = "^" + name + "$"
		}
		nextLetter := selectNextLetter(lastBigram, probabilities)
		if nextLetter == "" {
			break
		}
		name += nextLetter
	}
	return strings.Trim(name, "^$")
}

func selectNextLetter(lastBigram string, probabilities map[string]float64) string {
	rand.Seed(time.Now().UnixNano())
	sum := 0.0
	threshold := rand.Float64()

	for bigram, probability := range probabilities {
		if strings.HasPrefix(bigram, lastBigram) {
			sum += probability
			if sum > threshold {
				return string(bigram[1])
			}
		}
	}

	return "" // should never happen
}

func printProbabilities(probabilities map[string]float64) {
	fmt.Printf("%-3s %-7s\n", "Bigram", "Probability")
	for bigram, probability := range probabilities {
		fmt.Printf("%-3s %-7.5f\n", bigram, probability)
	}
}

func main() {
	data, err := readData("names.txt")
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	probabilities := calculateProbabilities(data)
	printProbabilities(probabilities)

	name := generateName(probabilities)
	fmt.Println("Generated name:", name)
}
