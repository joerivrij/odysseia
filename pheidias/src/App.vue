<template>
  <v-app id="odysseia">
      <v-card flat>
              <v-toolbar
                  dense
                  color=triadic>
                <v-toolbar-title>
                  <router-link to="/" tag="span" style="cursor: pointer">
                    {{ appTitle }}
                  </router-link>
                </v-toolbar-title>
                <v-spacer></v-spacer>
                <v-btn
                    text
                    v-for="item in menuItems"
                    :key="item.title"
                    :to="item.path">
                  <v-icon left dark>{{ item.icon }}</v-icon>
                  {{ item.title }}
                </v-btn>
              </v-toolbar>
      </v-card>
    <v-main>
      <h2></h2>
      <router-view></router-view>
    </v-main>
  </v-app>
</template>


<script>
export default {
  name: "App",
  data(){
    return {
      appTitle: 'Odysseia',
      sidebar: false,
      menuItems: [
        { title: 'Quiz', path: '/quiz', icon: 'list' },
        { title: 'Text', path: '/texts', icon: 'subject' },
        { title: 'Search', path: '/search', icon: 'search' }
      ]
    }
  },
};
</script>


<style>
</style>