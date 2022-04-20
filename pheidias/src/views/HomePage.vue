<template>
  <div id="app">
    <v-app id="homepage"
    >
      <v-parallax
          class=".v-parallax__image-container"
          dark
          height="900"
          :src="schoolOfAthens"
      >
        <v-row
            align="center"
            justify="center"
        >
          <v-col
              class="text-center"
              cols="12"
          >

            <h2 class="text-h2 text--primary">
              Welcome to Odysseia!
            </h2>
            <h5 class="text-md-h5 text--primary">
              This is an attempt at creating an app to help you learn Ancient Greek.<br>
              To start your own Odysseia select any of the components in the top navbar.
            </h5>
          </v-col>
        </v-row>
      </v-parallax>
    </v-app>
  </div>
</template>

<script>
export default {
  name: "HomePage",
  data () {
    return {
      schoolOfAthens: require('@/assets/school_of_athens.jpg'),
      drawer: true,
      group: null,
    }
  }
}
</script>

<style scoped>
.v-parallax__image-container {
  opacity: 0.1;
}
</style>