<template>
  <div id="quizarea">
    <v-app
        id="navbar"
        :style="{background: $vuetify.theme.themes[theme].background}"
    >
      <v-main>
        <div class="text-center">
          <div style="margin-bottom:2em;">
            <v-btn
                class="ma-2"
                color="primary"
                v-on:click="displayInfo=!displayInfo;stepper=1"
            >
              Howto
              <v-icon
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
                  Select Method and Category
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step
                    :complete="stepper > 2"
                    step="2"
                >
                  Select Chapter
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step step="3">
                  Answer Questions
                </v-stepper-step>
              </v-stepper-header>

              <v-stepper-items>
                <v-stepper-content step="1">
                  <v-card
                      class="mb-12"
                      color="white"
                      height="14em"
                  >Press one of the buttons below "Available Methods". And a popup will appear with each category for that method.
                  <br>
                    Some of these are based on books (such as Aristophanes -> frogs) and others on more general terms.
                    <br>
                    MOUSEION is Dutch only.
                  </v-card>

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
                      height="14em"
                  >Chapters can be chosen by typing and pressing SET CHAPTER.
                    <br>
                    Between brackets (for example: Chapters (1 -93)) are the available chapters
                  <br>
                  You can only set a chapter after you have chosen a method and catergory. Each chapter is 10 words at this time.</v-card>


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
                      height="14em"
                  >After having set a Method with a Category and a chapter (default is 1) you can start testing your knowledge.
                    <br>
                    For example:
                    <br>
                    Translate:
                    ἀλλά
                    <br>
                    Choose one of the options below and click the right answer. You will be shown whether you are right or not!
                    A graph is displayed just for fun.
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
            <v-container v-if="category.length">
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

          <div v-if="quizWord.length">
          <h2>Method: {{this.selectedMethod}} - Category: {{ this.category}} - Chapter {{this.selectedChapter}}</h2>
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
                color="#e9501d"
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
          <div v-if="!showAnswer && showButtons">
            <v-btn
                v-for="(item, index) in answers"
                :key="index"
                @click="selectedAnswer = item"
                v-on:click="postAnswer();showLoader();"
                class="ma-4"
                color="triadic"
                :width="widthStyle"
            >
              <span>{{ item }}</span>
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

          <div>
            <div style="width:75%; margin: 5em auto;">
              <v-data-table
                  :disable-sort="true"
                  :headers="headers"
                  :items="historyTable"
                  :items-per-page="5"
                  class="elevation-1"
              >
                <template v-slot:item.input="{ item }">
                  <v-chip
                      :color="item.color"
                  >
                    {{ item.input }}
                  </v-chip>
                </template>
              </v-data-table>
           </div>
          </div>
          <v-sparkline
              :value="this.graphNumbers"
              :gradient="['#1affcb', '#e9b61d', '#e9501d']"
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
      headers: [
        { text: 'Greek', value: 'greek', align: 'center' },
        { text: 'Input', value: 'input', align: 'center' },
        {text: 'Answer', value: 'answer', align: 'center'},
        {text: 'Method', value: 'method', align: 'center'},
        {text: 'Category', value: 'category', align: 'center'}
      ],
      historyTable: [],
      widthStyle : "50%",
      valid: true,
      showButtons: false,
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
      stepper: 1,
      displayInfo: false,
    }
  },
  methods: {
    showLoader () {
      this.value = 0

      this.interval = setInterval(() => {
        if (this.value >= 100) {
          clearInterval(this.interval)
        }
        this.value += 6
      }, 100)
    },
    hideAlert: function () {
      setTimeout(() => {
        if(this.showAnswer === true) {
          this.getQuestion()
          this.showAnswer = false
        }
      }, 2000);
    },
    async getQuestion () {
      this.showButtons = false
      let url = `${this.$sokratesUrl}/createQuestion?method=${this.selectedMethod}&category=${this.category}&chapter=${this.selectedChapter}`
      this.$apiClient.get(url)
          .then(async (response) => {
            this.quizWord = response.data[0];
            this.correctAnswer = response.data[1]
            let slicedArray = response.data.slice(1, 5)
            this.answers = await this.createNewArray(slicedArray);
            this.showButtons = true
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    async createNewArray(shuffeledArray) {
      for (let i = shuffeledArray.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffeledArray[i], shuffeledArray[j]] = [shuffeledArray[j], shuffeledArray[i]];
      }

      return shuffeledArray
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
            this.showAnswer = true;
            if (this.correct) {
              this.correctlyAnswered++
            }

            let color = "#1de9b6"
            if(!this.correct) {
              color = "#e9501d"
            }

            let lastAnswer = {
              greek: this.quizWord,
              color: color,
              answer: this.correctAnswer,
              input: this.selectedAnswer,
              category: this.category,
              method: this.selectedMethod,
            }

            this.historyTable.unshift(lastAnswer)

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
    isMobile() {
      return screen.width <= 800;
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
    if (this.isMobile()) {
      this.widthStyle = "90%"
    }
    else {
      this.flex = 6;
    }
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
