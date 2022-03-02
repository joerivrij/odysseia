<template>
  <div id="quizarea">
    <v-app
        id="navbar"
        :style="{background: $vuetify.theme.themes[theme].background}"
    >
      <v-main>
        <div class="text-center">
          <h2>Method: {{this.selectedMethod}} - Category: {{ this.category}} - Chapter {{this.selectedChapter}}</h2>
          <h4>Available Methods</h4>
          <v-row justify="center" align="center">
            <v-menu
                v-for="(method, index) in this.methods"
                :key="method"
                transition="slide-y-transition"
            >
              <template v-slot:activator="{ attrs, on }">
                <v-btn
                    v-bind="attrs"
                    v-on="on"
                    rounded
                    color="primary"
                    class="ma-2"
                    dark
                    @click="selectedMethod = methods[index];"
                    v-on:click="getCategories(methods[index])"
                >
                  {{ method }}
                </v-btn>
              </template>
              <v-list>
                <v-list-item
                    v-for="item in categories"
                    :key="item"
                    @click="category = item;setChapter(1)"
                    v-on:click="getChapters"
                >
                  <v-list-item-title v-text="item"></v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-row>
          <br />
          <br />
          <br />
            <v-container>
              <v-row justify="center" align="center">
                <v-col
                    cols="12"
                    md="4"
                >
                  <v-text-field
                      v-model="inputChapter"
                      min="1"
                      :max=chapters
                      :rules="numberRules"
                      :label=labelText
                      type="number"
                      required
                  ></v-text-field>
                  <v-btn
                      :disabled="!valid"
                      color="primary"
                      dark
                      rounded
                      @click="validate()"
                  >
                    Set Chapter
                  </v-btn>
                </v-col>
              </v-row>
            </v-container>
          <br>
          <br>

          <h3>Translate:</h3>
          <h3>{{quizWord}}</h3>
          <br />
          <div v-if="showAnswer">
            <div style="min-height: 4px">
              <v-progress-linear
                  v-model="value"
                  width="100%"
              ></v-progress-linear>
            </div>
            <v-alert
                v-if="correct"
                dense
                prominent
                type="success"
                icon="done"
                color="#1de9b6"
                transition="slide-y-transition"
            >
              <v-row align="center">
                <v-col class="grow">
                  The correct answer was indeed {{correctAnswer}}
                </v-col>
                <v-col class="shrink">
                  <v-btn
                      id="nextButton"
                      v-on:click="getNextQuestion"
                      rounded
                      color="primary"
                      dark
                  >
                    Next
                  </v-btn>
                </v-col>
              </v-row>
            </v-alert>
            <v-alert
                v-if="!correct"
                dense
                prominent
                type="error"
                icon="close"
                color="#ff5252"
                transition="slide-y-transition"
            >
              <v-row align="center">
                <v-col class="grow">
                  The correct answer was: {{correctAnswer}}
                </v-col>
                <v-col class="shrink">
                  <v-btn
                      id="nextButton"
                      v-on:click="getNextQuestion"
                      rounded
                      color="primary"
                      dark
                  >
                    Next
                  </v-btn>
                </v-col>
              </v-row>
            </v-alert>
          </div>
          <div v-if="!showAnswer">
            <v-btn
                v-for="(item, index) in answers"
                :key="index"
                @click="selectedAnswer = item"
                v-on:click="postAnswer();showLoader();"
                class="ma-2"
                color="triadic"
            >
              {{ item }}
            </v-btn>
          </div>
          <br />
          <br />
          <h4>
            You answered {{correctlyAnswered}} of {{answered}} correctly ({{percentage
            }}%)
          </h4>
          <v-btn v-on:click="resetProgress" rounded color="primary" dark>
            Reset Progress
          </v-btn>

          <v-sparkline
              :value="this.graphNumbers"
              :gradient="['#1feaea', '#ffd200', '#f72047']"
              :smooth="radius || false"
              :padding="padding"
              :line-width="width"
              :stroke-linecap="lineCap"
              :gradient-direction="gradientDirection"
              :fill="fill"
              :type="type"
              :auto-line-width="autoLineWidth"
              auto-draw
          ></v-sparkline>
        </div>
      </v-main>
    </v-app>
  </div>
</template>

<script>


export default {
  name: "QuizArea",
  computed: {
    theme(){
      return (this.$vuetify.theme.dark) ? 'dark' : 'light'
    }
  },
  data() {
    return {
      valid: true,
      alignments: [
        'start',
        'center',
        'end',
      ],
      inputChapter: 1,
      showAnswer: false,
      correctAnswer: "",
      quizWord: [],
      answers: [],
      selectedAnswer : "",
      correct: false,
      category: "",
      selectedMethod: "",
      answered: 0,
      correctlyAnswered: 0,
      percentage: 100,
      graphNumbers: [0],
      chapters : '',
      labelText: `Chapters: (1 - 1)`,
      numberRules: [
        v => !!v || 'Chapter is required',
        v => (v && v <= this.chapters) || 'Chapter cannot exceed chapters',
        v => (v && v !== 0) || 'non zero'
      ],
      closeOnContentClick: true,
      selectedChapter : 1,
      categories: [],
      methods: [],
      value: 0,
      interval: 0,
      width: 2,
      radius: 10,
      padding: 8,
      lineCap: 'round',
      gradientDirection: 'top',
      fill: false,
      type: 'trend',
      autoLineWidth: false,
    }
  },
  methods: {
    showLoader () {
      this.value = 0

      this.interval = setInterval(() => {
        if (this.value === 100) {
          clearInterval(this.interval)
        }
        this.value += 4
      }, 100)
    },
    hideAlert: function () {
      setTimeout(() => {
        if(this.showAnswer === true) {
          this.getQuestion()
          this.showAnswer = false
        }
      }, 3000);
    },
    getQuestion: function () {
      console.log('getting question')
      let url = `${this.$sokratesUrl}/createQuestion?method=${this.selectedMethod}&category=${this.category}&chapter=${this.selectedChapter}`
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
    getCategories: function (method) {
      let url = `${this.$sokratesUrl}/methods/${method}/categories`
      this.$apiClient.get(url)
          .then((response) => {
            this.categories = []
            let i;
            for (i in response.data.categories) {
              const category = response.data.categories[i].category
              this.categories.push(category)
            }
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    getNextQuestion: function () {
      clearInterval(this.interval)
      this.value = 0
      this.showAnswer = false;
      this.getQuestion()
    },
    postAnswer: function () {
      this.showAnswer = true;
      this.answered++
      this.$apiClient({
        method: 'post',
        url: `${this.$sokratesUrl}/answer`,
        data: {
          "answerProvided" : this.selectedAnswer,
          "quizWord": this.quizWord,
          "category": this.category
        }
      })
          .then((response) => {
            this.correct = response.data['correct']
            if (this.correct) {
              this.correctlyAnswered++
            }

            this.percentage = Math.round(this.correctlyAnswered / this.answered * 100)
            let inNumbers = Math.round(this.correctlyAnswered / this.answered * 10)
            this.graphNumbers.push(inNumbers)
          })
          .catch(e => {
            this.errors.push(e)
          })
      this.hideAlert()
    },
    setCategory(category) {
      this.category = category
    },
    resetProgress : function () {
      this.correctlyAnswered = 0
      this.answered = 0
      this.percentage = 100
    },
    getChapters : function () {
      let url = `${this.$sokratesUrl}/methods/${this.selectedMethod}/categories/${this.category}/chapters`
      this.$apiClient.get(url)
          .then((response) => {
            this.chapters = response.data['lastChapter']
            this.labelText = `Chapters: (1 - ${response.data['lastChapter']})`
          })

          .catch(e => {
            this.errors.push(e)
          })
    },
    getMethods : function () {
      let url = `${this.$sokratesUrl}/methods`
      this.$apiClient.get(url)
          .then((response) => {
            this.methods = []
            let i;
            for (i in response.data.methods) {
              const method = response.data.methods[i].method
              this.methods.push(method)
            }
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    setChapter(chapter) {
      this.selectedChapter = chapter
      this.getQuestion()
    },
    validate() {
      const isNotZero = this.inputChapter < 1
      const higher = this.inputChapter > this.chapters
      if (!isNotZero && !higher) {
        this.setChapter(this.inputChapter)
      }
    },
  },
  mounted() {
    this.getMethods()
    this.correctlyAnswered = 0
    this.answered = 0
    this.percentage = 100
  },
  created() {
    this.getMethods()
  },
  beforeDestroy () {
    clearInterval(this.interval)
  },

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
