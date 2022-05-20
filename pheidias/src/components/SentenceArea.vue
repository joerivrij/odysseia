<template>
  <div id="herodotos">
    <v-app
        id="sentencearea"
        :style="{background: $vuetify.theme.themes[theme].background}"
    >
      <v-main>
        <div class="text-center">
          <div style="margin-bottom:2em;">
            <v-btn
                class="ma-2"
                color="primary"
                dark
                v-on:click="displayInfo=!displayInfo;stepper=1"
            >
              Howto
              <v-icon
                  dark
                  right
              >
                info
              </v-icon>
            </v-btn>
            <v-stepper v-model="stepper" v-if="displayInfo">
              <v-stepper-header>
                <v-stepper-step
                    :complete="stepper > 1"
                    step="1"
                >
                  Select Author
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step
                    :complete="stepper > 2"
                    step="2"
                >
                  Select Book
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step step="3">
                  Type Translation
                </v-stepper-step>
              </v-stepper-header>

              <v-stepper-items>
                <v-stepper-content step="1">
                  <v-card
                      class="mb-12"
                      color="white"
                      height="5em"
                  >Press the Authors button and choose one of the available Authors (a pop up will appear)</v-card>

                  <v-btn
                      color="primary"
                      @click="stepper = 2"
                  >
                    Next Step
                  </v-btn>

                  <v-btn
                      text
                      v-on:click="displayInfo=!displayInfo">
                    Close
                  </v-btn>
                </v-stepper-content>

                <v-stepper-content step="2">
                  <v-card
                      class="mb-12"
                      color="white"
                      height="5em"
                  >Press the Books button and choose one of the available Books (a pop up will appear)</v-card>

                  <v-btn
                      color="primary"
                      @click="stepper = 3"
                  >
                    Next Step
                  </v-btn>

                  <v-btn
                      text
                      v-on:click="displayInfo=!displayInfo">
                    Close
                  </v-btn>
                </v-stepper-content>

                <v-stepper-content step="3">
                  <v-card
                      class="mb-12"
                      color="white"
                      height="12em"
                  >Type your translation and when done hit the CHECK button. You can toggle Translation assistance where
                    Odysseia will attempt to help you decline and translate the word (limited at this time).<br>
                    Any typos will be displayed after you have given your translation along side a "official" translation
                    </v-card>
                  <v-btn
                      text
                      v-on:click="displayInfo=!displayInfo">
                    Close
                  </v-btn>
                </v-stepper-content>
              </v-stepper-items>
            </v-stepper>
            </div>
          <v-card class="mx-auto" max-width="344">
            <v-card-text>
              <p class="text-h4 text--primary">
                {{ this.selectedAuthor }} Book: {{ this.selectedBook }}
              </p>
              <v-menu top :close-on-content-click="closeOnContentClick">
                <template v-slot:activator="{ on, attrs }">
                  <v-btn color="primary" dark v-bind="attrs" v-on="on" rounded>
                    Authors
                  </v-btn>
                </template>
                <v-list>
                  <v-list-item
                      v-for="(author, index) in authors"
                      :key="index"
                      v-on:click="setAuthorTo(author)"
                  >
                    <v-list-item-title>{{ author }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
              <v-menu top :close-on-content-click="closeOnContentClick">
                <template v-slot:activator="{ on, attrs }">
                  <v-btn color="primary" dark v-bind="attrs" v-on="on" rounded>
                    Books
                  </v-btn>
                </template>
                <v-list>
                  <v-list-item
                      v-for="(book, index) in books"
                      :key="index"
                      v-on:click="setBookTo(book)"
                  >
                    <v-list-item-title>{{ book }}</v-list-item-title>
                  </v-list-item>
              </v-list>
              </v-menu>
              <div class="text--primary">
                <v-textarea
                    readonly
                    v-model="sentence"
                ></v-textarea>
              </div>
              <h4>Possible Typos</h4>
              <v-list
                  v-for="(typo, index) in possibleTypos"
                  :key="index">
                <v-list-item>
                  <v-list-item-title>{{typo.verified}} => {{ typo.provided}}</v-list-item-title>
                </v-list-item>
              </v-list>
              <v-container fluid>
                <v-textarea
                    clearable
                    v-model="translationText"
                    clear-icon="close"
                    label="Type your translation here."
                    value="Type your translation here."
                ></v-textarea>
              </v-container>
            </v-card-text>
            <v-card-actions>
              <v-btn
                  text
                  color="teal accent-4"
                  @click="reveal = true"
                  v-on:click="checkAnswer();"
              >
                Check
              </v-btn>
              <v-btn
                  text
                  color="teal accent-4"
                  @click="reveal = false"
                  v-on:click="getNewSentence();"
              >
                Next
              </v-btn>
            </v-card-actions>
            <v-card
                v-if="reveal"
                class="transition-fast-in-fast-out v-card--reveal"
                style="height: 100%"
            >
              <v-card-text class="pb-0">
                <p class="text-h4 text--primary">Translation</p>
                <p>Correctness: {{this.translationPercentage}}%</p>
                {{ this.databaseAnswer }}
              </v-card-text>
              <v-card-actions class="pt-0">
                <v-btn text color="teal accent-4" @click="reveal = false">
                  Close
                </v-btn>
              </v-card-actions>
            </v-card>
            <div style="margin-bottom:1em; margin-left:1em;">
            <v-switch
                v-model="mobileView"
                label="Translation Assistance"
                color="primary"
                value="primary"
                hide-details
            ></v-switch>
            </div>
            <v-card-text v-if="mobileView">
              <v-chip-group
                  active-class="deep-purple accent-4 white--text"
                  column
              >
                <v-chip
                    v-for="(word, index) in seperatedWords"
                    :key="index"
                    v-on:click="queryWord(word)"
                    class="ma-1"
                    color="triadic"
                >
                  {{ word }}
                </v-chip>
              </v-chip-group>
            </v-card-text>
            <v-expand-transition>
              <v-card
                  elevation="24"
                  max-width="444"
                  class="mx-auto"
              >
                <v-system-bar lights-out></v-system-bar>
                <div class="text-overline mb-4">
                  Translation
                </div>
                <v-carousel
                    :continuous="false"
                    :cycle="cycle"
                    :show-arrows="true"
                    hide-delimiter-background
                    height="300"
                >
                  <v-carousel-item
                      v-for="(result, i) in grammarResults"
                      :key="i"
                  >
                    <v-sheet
                        height="100%"
                        tile
                    >
                      <v-row
                          class="fill-height"
                          align="center"
                          justify="center"
                      >
                        <v-list-item-title class="text-h5 mb-1">
                          {{ result.word }}
                        </v-list-item-title>
                        <v-list-item-title class="text-h5 mb-1">
                          {{ result.translation }}
                        </v-list-item-title>
                        <v-list-item-subtitle><strong>root:</strong> {{result.rootWord}}</v-list-item-subtitle>
                        <v-list-item-subtitle><strong>rule:</strong> {{result.rule}}</v-list-item-subtitle>
                      </v-row>
                    </v-sheet>
                  </v-carousel-item>
                </v-carousel>
                <v-list two-line>
                  <v-list-item>
                    <v-list-item-action>
                      <v-switch
                          v-model="cycle"
                          label="Cycle Translations"
                          inset
                      ></v-switch>
                    </v-list-item-action>
                  </v-list-item>
                </v-list>
              </v-card>
            </v-expand-transition>
          </v-card>
        </div>
      </v-main>
    </v-app>
  </div>
</template>

<script>
export default {
  name: "SentenceArea",
  computed: {
    theme(){
      return (this.$vuetify.theme.dark) ? 'dark' : 'light'
    }
  },
  data() {
    return {
      errors: [],
      authors: [],
      books: [],
      grammarResults: [],
      possibleTypos: [],
      seperatedWords: [],
      selectedAuthor: "",
      selectedBook: "",
      sentence: "",
      currentSentenceId: "",
      translationText: "",
      translationPercentage: "",
      databaseAnswer: "",
      closeOnContentClick: true,
      reveal: false,
      cycle: false,
      mobileView: false,
      displayInfo: false,
      stepper: 1,
    }
  },
  methods: {
    setAuthorTo(author) {
      this.selectedAuthor = author
      this.getBooks(this.selectedAuthor)
    },
    setBookTo(book) {
      this.selectedBook = book
      this.getNewSentence()
    },
    queryWord: function (word) {
      let url = `${this.$dionysiosUrl}/checkGrammar?word=${word}`
      this.$apiClient.get(url)
          .then((response) => {
            let extraTranslation = []
            for (let i = 0; i < response.data.results.length; i++) {
              if (response.data.results[i].translation === "") {
                response.data.results[i].translation = "No translation found"
              }

              if (response.data.results[i].translation.length > 25) {
                let words = response.data.results[i].translation.split(";")
                if (words.length > 1) {
                  for (let j = 0; j < words.length; j++) {
                    let rule = response.data.results[i].rule
                    let rootWord = response.data.results[i].rootWord
                    extraTranslation.push(
                        {
                          "word": word,
                          "rule": rule,
                          "rootWord": rootWord,
                          "translation": words[j]
                        }
                    )
                  }

                  response.data.results.splice(i, 1)
                }
              }
            }

            this.grammarResults = response.data.results
            if (extraTranslation !== []) {
              for (let i= 0; i < extraTranslation.length; i++) {
                this.grammarResults.push(extraTranslation[i])
              }
            }
          })
          .catch(e => {
            this.grammarResults =  [{
              "word"  :  word,
              "translation"   :  "No translation found",
              "rootWord"      :  word,
              "rule" : "No rule found"
            }]

            this.errors.push(e)
          })
    },
    getAuthors: function () {
      let url = `${this.$herodotosUrl}/authors`
      this.$apiClient.get(url)
          .then((response) => {
            let i;
            for (i in response.data.authors) {
              const author = response.data.authors[i].author
              const nameCapitalized = author.charAt(0).toUpperCase() + author.slice(1)
              this.authors.push(nameCapitalized)
            }

            this.selectedAuthor = this.authors[0]
            this.getBooks(this.selectedAuthor)
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    getBooks: function (author) {
      this.books = []
      if (author === "") {
        author = this.selectedAuthor
      }

      let url = `${this.$herodotosUrl}/authors/${author}/books`
      this.$apiClient.get(url)
          .then((response) => {
            let i;
            for (i in response.data.books) {
              const book = response.data.books[i].book
              this.books.push(book)
            }
            this.setBookTo(this.books[0])
          })
          .catch(e => {
            console.log(e)
            this.errors.push(e)
          })
    },
    getNewSentence: function () {
      const author = this.selectedAuthor.toLowerCase()
      let url = `${this.$herodotosUrl}/createQuestion?author=${author}&book=${this.selectedBook}`
      this.$apiClient.get(url)
          .then((response) => {
            this.sentence = response.data.sentence
            let cleanedText = response.data.sentence
                .replaceAll(",", "")
                .replaceAll(".", "")
                .replaceAll("—", " ")
            this.seperatedWords = cleanedText.split(" ")
            this.currentSentenceId = response.data.sentenceId
          })
          .catch(e => {
            console.log(e)
            this.errors.push(e)
          })
    },
    checkAnswer: function () {
      const author = this.selectedAuthor.toLowerCase()
      this.$apiClient({
        method: 'post',
        url: `${this.$herodotosUrl}/checkSentence`,
        data: {
          "answerSentence": this.translationText,
          "sentenceId": this.currentSentenceId,
          "author": author
        }
      }).then((response) => {
        this.translationPercentage = response.data.levenshteinPercentage
        this.databaseAnswer = response.data.quizSentence
        this.searchPossibleTypos(response.data.nonMatchingWords)
      })
    },
    searchPossibleTypos: function (words) {
      this.possibleTypos = []
      let i;
      for (i in words) {
        let j;
        for (j in words[i].matches) {
          const percentage = parseFloat(words[i].matches[j].percentage);
          if (words[i].matches[j].levenshtein <= 2 && percentage > 50.00) {
            let jsonMap = {
              "provided": words[i].word,
              "verified": words[i].matches[j].match
            }
            this.possibleTypos.push(jsonMap)
          }
        }
      }
    },
  },
  mounted() {
    this.getAuthors();
  },
  created() {
  }
}
</script>

<style scoped>
h2 {
  margin: 1em;
}

h3 {
  margin: 0.5em;
}

h4 {
  margin: 0.5em;
}
</style>
