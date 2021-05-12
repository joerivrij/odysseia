<template>
  <div id="herodotos">
    <div id="navbar">
      <main id="page-wrap" class="navbar">
        <h2>Available Authors</h2>
        <div v-for="author in authors" :key="author" class="buttonContainer">
          <button class="chapterButton" v-on:click="setAuthorTo(author)">
            <a>{{author}}</a>
          </button>
        </div>
      </main>
    </div>
  <div id="sentencearea" class="sentencearea">
    <h2>Sentence to translate</h2>
    <h3>{{sentence}}</h3>
    <p>Official Translation:</p>
    <h4>{{databaseAnswer}}</h4>
    <p>Correctness:</p>
    <a> {{translationPercentage}}%</a>
    <p>Your Translation:</p>
    <h4 style="white-space: pre-line;">{{ translationText }}</h4>
    <p>Type your translation here:</p>
    <textarea v-model="translationText" class="textbox" placeholder="Type your translation here"></textarea>
    <br>
    <button class="categorieButton"
            id="translateButton"
            v-on:click="checkAnswer">
      <strong>Check Translation</strong></button>
  </div>
  </div>
</template>

<script>
export default {
  name: "SentenceArea",
  data() {
    return {
      authors: [],
      chosenAuthor: "",
      sentence: "",
      currentSentenceId: "",
      translationText: "",
      translationPercentage: "",
      databaseAnswer: "",
    }
  },
  methods: {
    setAuthorTo(author) {
      this.chosenAuthor = author
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
      const author = this.chosenAuthor.toLowerCase()
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
      const author = this.chosenAuthor.toLowerCase()
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
    }
  },
    mounted() {
      this.getAuthors();
    },
    created() {
    }
  }

</script>

<style scoped>

h4 {
  margin-top: 2em;
}

h3 {
  margin-top: 0.5em;
}

a {
  cursor: pointer;
}

.textbox {
  width: 50%;
  height: 15em;
}
.sentencearea {
  background: lightblue;
  text-align: center;
  width: 100%;
  padding: 2em 0 0 0;
}

.categorieButton {
  background-color: #f96332;
  color: #fff;
  border-width: 2px;
  font-weight: 400;
  font-size: .8571em;
  line-height: 1.35em;
  margin: 10px 1px;
  border: none;
  border-radius: .1875rem;
  padding: 1em 2em;
  cursor: pointer;
  text-align: center;
  vertical-align: middle;
}

.categorieButton:hover {
  background: rgba(249, 99, 49, 0.71);
}

.wrongText {
  color: red
}

.rightText {
  margin: 0.1em;
  color: green;
}

#navbar {
  height: 100%;
  background-color: cadetblue;
  color: white;
  cursor: pointer;
  padding: 2em 0 2em 0;
  width: 100%;
  border: none;
  outline: none;
  font-size: 15px;
  text-align: center;
}

.buttonContainer {
  display: inline-block;
}

.chapterButton {
  background-color: #f96332;
  color: #fff;
  border-width: 2px;
  font-weight: 400;
  font-size: .8571em;
  line-height: 1.35em;
  margin: 10px 1px;
  border: none;
  border-radius: .1875rem;
  padding: 11px 22px;
  cursor: pointer;
  text-align: center;
  vertical-align: middle;
}

.chapterButton:hover {
  background: rgba(249, 99, 49, 0.71);
}
</style>