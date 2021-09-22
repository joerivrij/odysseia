<template>
  <div id="herodotos">
    <v-app id="sentencearea">
      <v-content>
        <div class="text-center">
          <v-card class="mx-auto" max-width="344">
            <v-card-text>
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
              <p class="text-h4 text--primary">
                {{ this.selectedAuthor }}
              </p>
              <p>RESERVED FOR CHAPTER</p>
              <div class="text--primary">
                {{ this.sentence }}
              </div>

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

            <v-expand-transition>
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
            </v-expand-transition>
          </v-card>
        </div>
      </v-content>
    </v-app>
  </div>
</template>

<script>
export default {
  name: "SentenceArea",
  data() {
    return {
      authors: [],
      selectedAuthor: "",
      sentence: "",
      currentSentenceId: "",
      translationText: "",
      translationPercentage: "",
      databaseAnswer: "",
      closeOnContentClick: true,
      reveal: false,
    }
  },
  methods: {
    setAuthorTo(author) {
      this.selectedAuthor = author
      this.getNewSentence()
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
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    getNewSentence: function () {
      const author = this.selectedAuthor.toLowerCase()
      let url = `${this.$herodotosUrl}/createQuestion?author=${author}`
      this.$apiClient.get(url)
          .then((response) => {
            this.sentence = response.data.sentence
            this.currentSentenceId = response.data.sentenceId
          })
          .catch(e => {
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
        console.log(response.data)
        this.translationPercentage = response.data.levenshteinPercentage
        this.databaseAnswer = response.data.quizSentence
      })
    },

    setDefaultAuthor: function () {
      if (this.authors.length === 0) {setTimeout(() => {
        this.selectedAuthor = this.authors[0]
        this.getNewSentence()
      }, 500)
      }
      this.selectedAuthor = this.authors[0]
      this.getNewSentence()

    }
  },
  mounted() {
    this.getAuthors();
    this.setDefaultAuthor();
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
