<template>
  <div id="searcharea">
    <div id="navbar">
      <main id="page-wrap" class="navbar">
      </main>
    </div>
    <div id="sentencearea" class="sentencearea">
        <h1 class="cover-heading">Alexandros Search</h1>
          <form v-on:submit.prevent="submitSearch" class="livesearch">
            <input type="text" v-model="searchWord" placeholder="Type here" @keyup="submitSearch">
          </form>
      <h2>Querying:</h2>
      <h3>{{searchWord}}</h3>
      <div v-for="result in searchResult" :key="result" class="buttonContainer">
        <a>Greek: <strong>{{result.greek}}</strong> English: <strong>{{result.english}}</strong></a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SearchArea",
  data() {
    return {
      searchWord: "",
      searchResult: "",
      errors: [],
    }
  },
  methods: {
    submitSearch: function () {
      let url = `${this.$alexandrosUrl}/search?word=${this.searchWord}`
      this.$apiClient.get(url)
          .then((response) => {
            console.log(response.data)
            this.searchResult = response.data
          })
          .catch(e => {
            this.errors.push(e)
            this.searchResult = ""
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


.sentencearea {
  background: lightblue;
  text-align: center;
  width: 100%;
  height: 35em;
  padding: 2em 0 0 0;
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

* {
  box-sizing: border-box;
}

/* Style the search field */
form.livesearch input[type=text] {
  padding: 10px;
  font-size: 17px;
  border: 1px solid grey;
  text-align: center;
  width: 40%;
  background: #f1f1f1;
}


</style>