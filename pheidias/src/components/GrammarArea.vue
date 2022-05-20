<template>
  <div id="grammar">
    <v-app
        id="grammararea"
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
                  Copy word
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step step="2">
                  Search for results
                </v-stepper-step>
              </v-stepper-header>

              <v-stepper-items>
                <v-stepper-content step="1">
                  <v-card
                      class="mb-12"
                      color="white"
                      height="14em"
                  >Copy or type a word into the search bar<br>
                    For example:<br>
                    ἔβαλλε <br>
                    φέροντος <br>
                    ἀληθῆ <br>
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
                  >Press enter or the magnifying glass. Results will be shown in the table</v-card>

                  <v-btn
                      text
                      v-on:click="displayInfo=!displayInfo">
                    Close
                  </v-btn>
                </v-stepper-content>
              </v-stepper-items>
            </v-stepper>
          </div>
      <v-card color="background" dark>
      <v-app-bar
      height="150em"
      prominent
      color="primary"
      >
        <v-app-bar-title>Grammar</v-app-bar-title>
        <template v-slot:extension>
        <v-text-field
            label="What Greek word are you trying to decline?"
            v-model="queryWord"
            @keydown.enter="queryGrammarWord"
        ></v-text-field>
          <v-btn icon
                 @click="queryGrammarWord">
            <v-icon>mdi-magnify</v-icon>
            </v-btn>
        </template>
        </v-app-bar>
        <v-expand-transition>
          <v-card light color="background">
            <v-card-text>
              <h2>Results</h2>
              <br/>
              <v-data-table
                  dense
                  :headers="headers"
                  :items="grammarResults"
                  :items-per-page="10"
                  item-key="name"
                  class="elevation-1"
              ></v-data-table>
            </v-card-text>
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
  name: "GrammarArea",
  computed: {
    theme(){
      return (this.$vuetify.theme.dark) ? 'dark' : 'light'
    }
  },
  data() {
    return {
      grammarResults: [],
      headers: [
        {
          text: 'Queried',
          align: 'start',
          sortable: true,
          value: 'word',
        },
        { text: 'Rule', value: 'rule' },
        { text: 'Root', value: 'rootWord' },
        { text: 'Translation', value: 'translation' },
      ],
      queryWord: "",
      errors: [],
      select: null,
      stepper: 1,
      displayInfo: false,
    }
  },
  methods: {
    queryGrammarWord: function () {
      let url = `${this.$dionysiosUrl}/checkGrammar?word=${this.queryWord}`
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
                          "word": this.queryWord,
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
              "word"  :  this.queryWord,
              "translation"   :  "No translation found",
              "rootWord"      :  this.queryWord,
              "rule" : "No rule found"
            }]

            this.errors.push(e)
          })
    },
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

* {
  box-sizing: border-box;
}

/* Style the search field */
form.livesearch input[type="text"] {
  padding: 10px;
  font-size: 17px;
  border: 1px solid grey;
  text-align: center;
  width: 40%;
  background: #f1f1f1;
}
</style>
