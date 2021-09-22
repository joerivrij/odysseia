<template>
  <div id="search">
    <v-app id="notfound">
      <v-content>
        <div class="text-center">
          <h1>Error - page not found</h1>
          <v-btn
              text
              v-for="item in menuItems"
              :key="item.title"
              :to="item.path"
          >
            <v-icon left dark>{{ item.icon }}</v-icon>
            {{ item.title }}
          </v-btn>
        </div>
      </v-content>
    </v-app>
  </div>
</template>
<script>
export default {
  name: "App",
  data(){
    return {
      appTitle: 'Odysseia',
      sidebar: false,
      menuItems: [
        { title: 'ΠΑΝΤΑΡΕΙ', path: '/', icon: 'home' },
      ]
    }
  },
};
</script>
