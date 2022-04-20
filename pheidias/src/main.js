import Vue from 'vue'
import App from './App.vue'
import axios from 'axios';
import router from './router'
import store from './store'
import vuetify from '@/plugins/vuetify'
import "vuetify/dist/vuetify.min.css";

Vue.config.productionTip = false

Vue.http = Vue.prototype.$apiClient = axios;
Vue.prototype.$sokratesUrl = store.state.SokratesUrl
Vue.prototype.$herodotosUrl = store.state.HerodotosUrl
Vue.prototype.$alexandrosUrl = store.state.AlexandrosUrl
Vue.prototype.$dionysosUrl = store.state.DionysosUrl

if (process.env.NODE_ENV === 'development') {
  Vue.prototype.$sokratesUrl = "https://odysseia-greek.com/sokrates/v1"
  Vue.prototype.$herodotosUrl = "https://odysseia-greek.com/herodotos/v1"
  Vue.prototype.$alexandrosUrl = "https://odysseia-greek.com/alexandros/v1"
  Vue.prototype.$dionysosUrl = "https://odysseia-greek.com/dionysos/v1"
}

console.log(`sokrates proxy found at: ${Vue.prototype.$sokratesUrl}`)
console.log(`herodotos proxy found at: ${Vue.prototype.$herodotosUrl}`)
console.log(`alexandros proxy found at: ${Vue.prototype.$alexandrosUrl}`)
console.log(`dionysos proxy found at: ${Vue.prototype.$dionysosUrl}`)


new Vue({
  vuetify,
  el: '#app',
  render: h => h(App),
  router,
  template: '<App/>',
}).$mount('#app');
