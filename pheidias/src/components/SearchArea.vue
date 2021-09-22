<template>
  <div id="search">
    <v-app id="searcharea">
      <v-card color="primary" dark>
        <v-card-text>
          Explore hundreds of words available on the database for more
          information see
          <a
              class="grey--text text--lighten-3"
              href="https://github.com/joerivrij/odysseia/tree/master/demokritos/lexiko"
              target="_blank"
          >the GitHub repository</a
          >.
        </v-card-text>
        <v-card-text>
          <v-autocomplete
              v-model="model"
              :items="items"
              :loading="loading"
              :search-input.sync="search"
              color="white"
              hide-no-data
              hide-selected
              item-text="Description"
              item-value="API"
              label="What Greek word are you looking for?"
              placeholder="Start typing to Search"
              prepend-icon="mdi-database-search"
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
  name: "SearchArea",
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
      items: [],
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
            for (let i = 0; i < response.data.length; i++) {
              this.items.push(response.data[i]['greek']);
            }

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
