package transform

import (
	"github.com/odysseia/plato/models"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestRemovingAccents(t *testing.T) {
	sourceString := "ἀγαθός"
	expected := "αγαθος"

	accentsRemoved := RemoveAccents(sourceString)

	assert.Equal(t, expected, accentsRemoved)
}

func TestRemovingAccentsFromCombinedWord(t *testing.T) {
	sourceString := "ἄλγος –ους, τό"
	expected := "αλγος –ους, το"

	accentsRemoved := RemoveAccents(sourceString)

	assert.Equal(t, expected, accentsRemoved)
}

func TestTransformAGreekWord(t *testing.T) {
	expected := "αναλαμβανω"
	m := models.Meros{
		Greek:      "ἀναλαμβάνω",
		English:    "pick up",
		LinkedWord: "ἀνά",
	}

	strippedWord := RemoveAccents(m.Greek)
	word := models.Meros{
		Greek:      strippedWord,
		English:    m.English,
		LinkedWord: m.LinkedWord,
		Original:   m.Greek,
	}

	assert.Equal(t, expected, word.Greek)
	assert.Equal(t, m.Greek, word.Original)
}
