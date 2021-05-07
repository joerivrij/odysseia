<template>
  <div id="quizarea" class="quizarea">

    <h2>Hoofdstuk {{this.$parent.selectedChapter}}</h2>
    <h4> Je hebt {{correctlyAnswered}} van de {{answered}} goed beantwoord ({{percentage}}%)</h4>
    <button class="categorieButton"
            v-on:click="resetProgress">
      <strong>Reset voortgang</strong></button>
    <h3>Je hebt gekozen voor {{categorie}}</h3>

    <button class="categorieButton"
            id="nominaButton"
            @click="categorie = 'nomina'"
            v-on:click="getNextQuestion">
      <strong>Nomina</strong></button>

    <button class="categorieButton"
            id="verbaButton"
            @click="categorie = 'verba'"
            v-on:click="getNextQuestion">
      <strong>Verba</strong></button>

    <button class="categorieButton"
            id="miscButton"
            @click="categorie = 'misc'"
            v-on:click="getNextQuestion">
      <strong>Misc</strong></button>


    <button class="categorieButton"
            id="logosButton"
            @click="categorie = 'logos'"
            v-on:click="getNextQuestion">
      <strong>Logos</strong></button>

    <h3>Vertaal: {{quizWord}}</h3><br>
    <div v-if="showAnswer">
      <h2 v-if="correct" class="rightText">Correct</h2>
      <h2 v-else class="wrongText">Verkeerd</h2>
      <h3 v-if="!correct"> Het juiste antwoord was: {{correctAnswer}} </h3>

      <button class="categorieButton"
              id="nextButton"
              v-on:click="getNextQuestion">
        <strong>Volgende</strong></button>
    </div>
    <div v-for="answer in answers" :key="answer">
      <button class="answerButton"
              @click="selectedAnswer = answer"
              v-on:click="postAnswer">{{answer}}</button>
    </div>

  </div>
</template>

<script>
export default {
  name: "QuizArea",
  data() {
    return {
      showAnswer: false,
      correctAnswer: "",
      quizWord: [],
      answers: [],
      selectedAnswer : "",
      correct: false,
      categorie: "nomina",
      answered: 0,
      correctlyAnswered: 0,
      percentage: 100,
    }
  },
  methods: {
    getQuestion: function () {
      let url = `${this.$proxyUrl}/createQuestion?category=${this.categorie}`
      if (this.$parent.selectedChapter !== "") {
        url = url + `&chapter=${this.$parent.selectedChapter}`
      }
      this.$apiClient.get(url)
          .then((response) => {
            let shuffeledArray = response.data.slice(1, 5)
            for (let i = shuffeledArray.length - 1; i > 0; i--) {
              const j = Math.floor(Math.random() * (i + 1));
              [shuffeledArray[i], shuffeledArray[j]] = [shuffeledArray[j], shuffeledArray[i]];
            }
            this.answers = shuffeledArray;
            this.quizWord = response.data[0];
            this.correctAnswer = response.data[1]
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    getNextQuestion: function () {
      this.showAnswer = false;
      this.setCategorieToParent(this.categorie)
      this.$root.$refs.navbar.getChapters()
      this.getQuestion()
    },
    postAnswer: function () {
      this.showAnswer = true;
      this.answered++
      this.$apiClient({
        method: 'post',
        url: `${this.$proxyUrl}/answer`,
        data: {
          "answerProvided" : this.selectedAnswer,
          "quizWord": this.quizWord,
          "category": this.categorie
        }
      })
          .then((response) => {
            this.correct = response.data['correct']
            if (this.correct) {
              this.correctlyAnswered++
            }

            this.percentage = Math.round(this.correctlyAnswered / this.answered * 100)
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    setCategorieToParent(categorie) {
      this.$emit("update-categorie", categorie);
    },
    resetProgress : function () {
      this.correctlyAnswered = 0
      this.answered = 0
      this.percentage = 100
    }
  },
  mounted() {
    this.getQuestion()
    this.correctlyAnswered = 0
    this.answered = 0
    this.percentage = 100
  },
  created() {
    this.$root.$refs.quizarea = this;
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
.quizarea {
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

.answerButton {
  background-color: cadetblue;
  color: #fff;
  border-width: 1em;
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
  width: 80%;
}

.answerButton:hover {
  background: darkseagreen;
}

.wrongText {
  color: red
}

.rightText {
  margin: 0.1em;
  color: green;
}
</style>