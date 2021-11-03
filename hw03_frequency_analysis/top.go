package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type FreqDist struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	words := strings.Fields(text)

	wordsFrequency := fillWordFrequency(words)

	freqStructures := turnToFreqDist(wordsFrequency)

	sorted := sortWords(freqStructures)

	return takeFirst10(sorted)
}

func fillWordFrequency(words []string) map[string]int {
	res := make(map[string]int)

	for _, val := range words {
		_, isset := res[val]

		if isset {
			res[val]++
			continue
		}

		res[val] = 1
	}

	return res
}

func turnToFreqDist(words map[string]int) []FreqDist {
	res := make([]FreqDist, 0, len(words))

	for key, value := range words {
		res = append(res, FreqDist{key, value})
	}

	return res
}

func sortWords(frequencies []FreqDist) []FreqDist {
	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].Count != frequencies[j].Count {
			return frequencies[i].Count > frequencies[j].Count
		}

		return frequencies[i].Word < frequencies[j].Word
	})

	return frequencies
}

func takeFirst10(words []FreqDist) []string {
	keys := make([]string, 0, len(words))

	for _, val := range words {
		keys = append(keys, val.Word)
		if len(keys) >= 10 {
			break
		}
	}

	return keys
}
