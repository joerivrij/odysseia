<template>
  <v-app id="odysseia">
      <v-card flat>
              <v-toolbar
                  color=triadic>
                <v-menu
                    bottom
                    left
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                        primary
                        icon
                        v-bind="attrs"
                        v-on="on"
                    >
                      <v-icon dark>mdi-menu</v-icon>
                    </v-btn>
                  </template>

                  <v-list
                      nav
                      width="25em"
                  >
                    <v-list-item v-for="(item, i) in menuItems" :key="i" :to="item.path" link>
                      <v-list-item-icon>
                        <v-icon>{{item.icon}}</v-icon>
                      </v-list-item-icon>
                      <v-list-item-title>{{item.title}}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </v-toolbar>
      </v-card>
    <v-main>
      <h2></h2>
      <router-view></router-view>
      <v-btn
          v-scroll="onScroll"
          v-show="fab"
          fab
          dark
          fixed
          bottom
          right
          color="triadic"
          @click="toTop"
      >
        <v-icon>keyboard_arrow_up</v-icon>
      </v-btn>
    </v-main>
  </v-app>
</template>


<script>
export default {
  name: "App",
  data(){
    return {
      fab: false,
      appTitle: 'Odysseia',
      closeOnClick: true,
      menuItems: [
        { title: 'Home', path: '/', icon: 'mdi-home-variant' },
        { title: 'Quiz', path: '/quiz', icon: 'mdi-alphabet-greek' },
        { title: 'Texts', path: '/texts', icon: 'mdi-bookshelf' },
        { title: 'Grammar', path: '/grammar', icon: 'mdi-feather' },
        { title: 'Dictionary', path: '/dictionary', icon: 'search' }
      ]
    }
  },
  methods: {
    onScroll (e) {
      if (typeof window === 'undefined') return
      const top = window.pageYOffset ||   e.target.scrollTop || 0
      this.fab = top > 20
    },
    toTop () {
      this.$vuetify.goTo(0)
    }
  }
};
</script>


<style>
</style>