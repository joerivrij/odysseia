<template>
  <div id="app">
    <v-app id="homepage"
    >
      <v-parallax
          dark
          height="800"
          :src="schoolOfAthens">

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

              <p class="text--black text--primary">Sokrates will test your knowledge of Ancient Greek words (in Dutch only at this time)</p>
              <p class="text--black text--primary">Herodotos is where you can start translating Ancient Greek texts</p>
              <p class="text--black text--primary">Alexandros will let you search the database for word translations</p>

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
      schoolOfAthens: require('@/assets/school_of_athens.jpg')
    }
  }
}
</script>

<style scoped>

</style>