<template>
  <div id="dictionary">
    <v-app
        id="dictionaryarea"
        :style="{background: $vuetify.theme.themes[theme].background}"
    >
      <v-card color="primary" dark>
        <v-card-text>
          Explore hundreds of words available on the database for more
          information see
          <a
              class="grey--text text--lighten-3"
              href="https://github.com/joerivrij/odysseia/tree/master/demokritos/lexiko"
              target="_blank"
          >the GitHub repository</a
          >
        </v-card-text>
        <v-card-text>
          <v-autocomplete
              :loading="loading"
              :search-input.sync="search"
              hide-no-data
              color="white"
              item-text="Description"
              item-value="API"
              label="What Greek word are you looking for?"
              placeholder="Start typing to Search"
              prepend-icon="mdi-database-search"
              auto-select-first
              return-object
          ></v-autocomplete>
        </v-card-text>
        <v-divider></v-divider>
        <v-expand-transition>
          <v-card light color="background">
            <v-card-text>
              <h2>Results</h2>
              <br />
              <v-data-table
                  dense
                  :headers="headers"
                  :items="searchResults"
                  :items-per-page="10"
                  item-key="name"
                  class="elevation-1"
              ></v-data-table>
            </v-card-text>
          </v-card>
        </v-expand-transition>
      </v-card>
    </v-app>
  </div>
</template>

<script>
export default {
  name: "DictionaryArea",
  computed: {
    theme(){
      return (this.$vuetify.theme.dark) ? 'dark' : 'light'
    }
  },
  data() {
    return {
      headers: [
        {
          text: 'Greek',
          align: 'start',
          sortable: true,
          value: 'greek',
        },
        { text: 'English', value: 'english' },
      ],
      searchResults: [],
      errors: [],
      loading: false,
      search: null,
      select: null,
    }
  },
  methods: {
    submitSearch: function (value) {
      this.loading = true
      this.searchResults = []
      let url = `${this.$alexandrosUrl}/search?word=${value}`
      this.$apiClient.get(url)
          .then((response) => {
            this.searchResults = response.data
            setTimeout(() => {
              this.loading = false
            }, 1500)
          })
          .catch(e => {
            console.log(e)
            setTimeout(() => {
              this.loading = false
            }, 1500)
          })
    },
  },
  watch: {
    search (val) {
      val && val !== this.select && this.submitSearch(val)
    },
  },
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
