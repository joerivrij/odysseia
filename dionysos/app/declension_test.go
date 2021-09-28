package app

import (
	"github.com/odysseia/plato/elastic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFirstDeclensionNominativusFemSing(t *testing.T) {
	fixtureFile := "declensionsDionysos"
	mockCode := 200
	mockElasticClient, err := elastic.CreateMockClient(fixtureFile, mockCode)
	assert.Nil(t, err)

	handler := DeclensionHandler{
		ElasticClient: *mockElasticClient,
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"τιμή", "μάχη", "οἰκία", "δόξα"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "fem")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - fem - nom")
		}
	}
}

func TestGetFirstDeclensionGenitivusFemSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"τιμῆς", "μάχης", "οἰκίας", "δόξης"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "fem")
				assert.NotEqual(t, "", declension.Translation)
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - fem - gen")
		}
	}
}

func TestGetFirstDeclensionDativusFemSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μάχῃ", "οἰκίᾳ", "δόξῃ"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "dat")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - fem - dat")
		}
	}
}

func TestGetFirstDeclensionAccusativusFemSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"τιμήν", "μάχην", "οἰκίαν", "δόξαν"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - sing - fem - acc")
	}
}

func TestGetFirstDeclensionNominativusFemPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"τιμαι", "μάχαι", "οἰκίαι", "δόξαι"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - plural - fem - nom")
	}
}

func TestGetFirstDeclensionGenitivusFemPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"τιμῶν", "μάχῶν", "χωρῶν", "δόξῶν"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "gen")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - fem - gen")
		}
	}
}

func TestGetFirstDeclensionDativusFemPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μάχαις", "οἰκίαις", "δόξαις"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "dat")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - fem - dat")
		}
	}
}

func TestGetFirstDeclensionAccusativusFemPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"τιμᾱς", "μάχας", "οἰκίας", "χώρᾱς"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "fem")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - fem - acc")
		}
	}
}

func TestGetFirstDeclensionNominativusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίας", "πολίτης"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - nom")
	}
}

func TestGetFirstDeclensionGenitivusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίου", "πολίτου", "κριτοῦ"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - gen")
	}
}

func TestGetFirstDeclensionDativusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίᾳ", "πολίτῃ", "κριτῇ"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - dat")
	}
}

func TestGetFirstDeclensionAccusativusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίαν", "πολίτην"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - acc")
	}
}

func TestGetFirstDeclensionNominativusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίαι", "πολίται", "κριταί"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - nom")
	}
}

func TestGetFirstDeclensionGenitivusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίῶν", "πολίτῶν", "κριτῶν"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "gen")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - gen")
		}
	}
}

func TestGetFirstDeclensionDativusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίαις", "πολίταις", "κριταῖς"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "dat")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - dat")
		}
	}
}

func TestGetFirstDeclensionAccusativusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"νεανίας", "πολίτας", "κριτᾱ́ς"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - acc")
		}
	}
}

func TestGetSecondDeclensionNominativusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"δοῦλος", "πόλεμος", "θεός"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - nom")
		}
	}
}

func TestGetSecondDeclensionGenitivusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"δοῦλου", "πόλεμου", "θεoῦ"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - gen")
		}
	}
}

func TestGetSecondDeclensionDativusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"δοῦλῳ", "πόλεμῳ", "θεῷ"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - dat")
		}
	}
}

func TestGetSecondDeclensionAccusativusMascSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"πόλεμον", "θεόν"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - masc - acc")
		}
	}
}

func TestGetSecondDeclensionNominativusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"δοῦλοι", "πόλεμοι", "θεοί"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - nom")
		}
	}
}

func TestGetSecondDeclensionGenitivusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"πόλεμων", "θεῶν"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "gen")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - gen")
		}
	}
}

func TestGetSecondDeclensionDativusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"πόλεμοις", "θεοῖς"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - dat")
		}
	}
}

func TestGetSecondDeclensionAccusativusMascPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"πόλεμους", "θεούς"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "masc")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - masc - acc")
		}
	}
}

func TestGetSecondDeclensionNominativusNeutSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μῆλον", "δῶρον"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "neut")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - neut - nom")
		}
	}
}

func TestGetSecondDeclensionGenitivusNeutSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μῆλου", "δῶρου" }
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "gen")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - neut - gen")
		}
	}
}

func TestGetSecondDeclensionDativusNeutSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μῆλῳ", "δῶρῳ"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - sing - neut - dat")
	}
}

func TestGetSecondDeclensionAccusativusNeutSing(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μῆλον", "δῶρον"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "neut")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - sing - neut - acc")
		}
	}
}

func TestGetSecondDeclensionNominativusNeutPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μῆλα", "δῶρα"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "neut")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - neut - nom")
		}
	}
}

func TestGetSecondDeclensionGenitivusNeutPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μήλων", "δῶρων" }
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "gen")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - neut - gen")
		}
	}
}

func TestGetSecondDeclensionDativusNeutPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μήλοις", "δῶροις"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		assert.Contains(t, declensions.Results[0].Rule, "noun - plural - neut - dat")
	}
}

func TestGetSecondDeclensionAccusativusNeutPlural(t *testing.T) {
	handler := DeclensionHandler{
		BaseUrl: "http://minikube-odysseia.test",
		Version: "v1",
		ApiName: "alexandros",
		SearchWordEndPoint: "search",
	}

	words := []string{"μῆλα", "δῶρα"}
	for _, word := range words {
		declensions, err := handler.StartFindingRules(word)
		assert.Nil(t, err)
		if len(declensions.Results) > 1 {
			for _, declension := range declensions.Results{
				assert.Contains(t, declension.Rule, "neut")
			}
		} else {
			assert.Contains(t, declensions.Results[0].Rule, "noun - plural - neut - acc")
		}
	}
}