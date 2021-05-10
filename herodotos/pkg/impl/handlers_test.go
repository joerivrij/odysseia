package impl

import "testing"
import "github.com/stretchr/testify/assert"

func TestLevenshteinAsInt(t *testing.T) {
	sourceString := "This is the display of the inquiry of Herodotos of Halikarnassos"
	targetString := "This is the historical inquiry of Herodotos of Halikarnassos"
	expected := 11

	levenshteinDist := levenshteinDistance(sourceString, targetString)

	assert.Equal(t, expected, levenshteinDist)
}

func TestLevenshteinToPercentageZero(t *testing.T) {
	source := 2
	longestWord := 2
	percentage := levenshteinDistanceInPercentage(source, longestWord)

	expected := float32(0)

	assert.Equal(t, expected, percentage)
}

func TestLevenshteinToPercentageHunderd(t *testing.T) {
	source := 0
	longestWord := 20
	percentage := levenshteinDistanceInPercentage(source, longestWord)

	expected := float32(100)

	assert.Equal(t, expected, percentage)
}

func TestLevenshteinToPercentageMixed(t *testing.T) {
	source := 11
	longestWord := 64
	percentage := levenshteinDistanceInPercentage(source, longestWord)

	expected := float32(82.8125)

	assert.Equal(t, expected, percentage)
}

func TestLongestSentence(t *testing.T) {
	sourceString := "this is a lot longer"
	targetString := "short"

	longestSentence := longestStringOfTwo(sourceString, targetString)
	expected := 20

	assert.Equal(t, expected, longestSentence)
}

func TestFindMatchingWords(t *testing.T) {
	sourceString := "This inquiry the Herodotos of Halikarnassos"
	targetString := "This inquiry tHe Herodotus of Halikarnassus"

	sameWords, bothWords := findMatchingWordsWithSpellingAllowance(sourceString, targetString)

	expectedSame := 4
	expectedSameWithSpelling := 2
	assert.Equal(t, expectedSame, len(sameWords))
	assert.Equal(t, expectedSameWithSpelling, len(bothWords))
}

func TestStreamlineString(t *testing.T) {
	sourceString := "This inquiry the Herodotos of Halikarnassos"
	matchingWords := []string{"the", "of"}

	newSentence := streamlineSentenceBeforeCompare(matchingWords, sourceString)

	expected := "This inquiry Herodotos Halikarnassos"
	assert.Equal(t, expected, newSentence)

}