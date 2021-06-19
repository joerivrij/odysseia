package main

import (
	"github.com/odysseia/plato/models"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestRemovingAccents(t *testing.T) {
	sourceString := "ἀγαθός"
	expected := "αγαθος"

	accentsRemoved := removeAccents(sourceString)

	assert.Equal(t, expected, accentsRemoved)
}

func TestRemovingAccentsFromCombinedWord(t *testing.T) {
	sourceString := "ἄλγος –ους, τό"
	expected := "αλγος –ους, το"

	accentsRemoved := removeAccents(sourceString)

	assert.Equal(t, expected, accentsRemoved)
}

func TestTransformAGreekWord(t *testing.T) {
	expected := "αναλαμβανω"
	m := models.Meros{
		Greek:      "ἀναλαμβάνω",
		English:    "pick up",
		LinkedWord: "ἀνά",
	}

	strippedWord := removeAccents(m.Greek)
	word := models.Meros{
		Greek:      strippedWord,
		English:    m.English,
		LinkedWord: m.LinkedWord,
		Original:   m.Greek,
	}

	assert.Equal(t, expected, word.Greek)
	assert.Equal(t, m.Greek, word.Original)
}
