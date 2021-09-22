<template>
  <div id="app">
    <v-app id="homepage"
    >
      <v-parallax
          dark
          id="opacity-parallax"
          height="800"
          src="https://upload.wikimedia.org/wikipedia/commons/4/49/%22The_School_of_Athens%22_by_Raffaello_Sanzio_da_Urbino.jpg">

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
  name: "HomePage"
}
</script>

<style scoped>
#opacity-parallax {
  opacity: 90%;
}
</style>