<template>
  <v-app id="odysseia">
      <v-card flat>
        <v-container fluid>
          <v-row class="child-flex">
            <div>
              <v-toolbar
                  color=secondary>
                <v-toolbar-title>
                  <router-link to="/" tag="span" style="cursor: pointer">
                    {{ appTitle }}
                  </router-link>
                </v-toolbar-title>
                <v-spacer></v-spacer>
                <v-toolbar-items class="hidden-xs-only">
                <v-btn
                    text
                    v-for="item in menuItems"
                    :key="item.title"
                    :to="item.path">
                  <v-icon left dark>{{ item.icon }}</v-icon>
                  {{ item.title }}
                </v-btn>
                  <v-icon>

                  </v-icon>
                </v-toolbar-items>
              </v-toolbar>
            </div>
          </v-row>
        </v-container>
      </v-card>
    <v-content>
      <h2></h2>
      <router-view></router-view>
    </v-content>
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
        { title: 'Home', path: '/', icon: 'home' },
        { title: 'Sokrates', path: '/multi', icon: 'list' },
        { title: 'Herodotos', path: '/texts', icon: 'subject' },
        { title: 'Alexandros', path: '/search', icon: 'search' }
      ]
    }
  },
};
</script>


<style>
</style>