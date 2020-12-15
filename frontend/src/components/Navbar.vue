<template>
  <div id="navbar">
    <main id="page-wrap" class="navbar">
        <h2 v-on:click="collapseBar()">Beschikbare hoofdstukken</h2>
          <div v-for="chapter in chapters" :key="chapter" class="buttonContainer">
            <button class="chapterButton" v-if="showBar">
              <a v-on:click="setChapterToParent(chapter)">
          {{chapter}}</a>
            </button>
          </div>
    </main>
  </div>
</template>

<script>
export default {
  name: "Navbar",
  data() {
    return {
      chapters : [],
      showBar: false,
    }
  },
  created() {
    this.getChapters()
  },
  components: {
  },
  methods: {
    getChapters : function () {
      let url = `${this.$proxyUrl}/chapters`
      this.$apiClient.get(url)
          .then((response) => {
            let endOfArray = response.data['chapters']
            for (let i = 1; i <= endOfArray; i++) {
              this.chapters.push(i)
            }
          })
          .catch(e => {
            this.errors.push(e)
          })
    },
    setChapterToParent(chapter) {
      this.$emit("update-chapter", chapter);
      this.showBar = !this.showBar
    },
    collapseBar : function () {
      this.showBar = !this.showBar
    }
  },
}
</script>

<style>
#navbar {
  height: 100%;
  background-color: cadetblue;
  color: white;
  cursor: pointer;
  padding: 18px;
  width: 100%;
  border: none;
  outline: none;
  font-size: 15px;
  text-align: center;
}

.buttonContainer {
  display: inline-block;
}

a {
  cursor: pointer;
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