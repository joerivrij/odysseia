import Vue from 'vue'
import App from './App.vue'
import axios from 'axios';
import router from './router'
import store from './store'

Vue.config.productionTip = false

Vue.http = Vue.prototype.$apiClient = axios;
Vue.prototype.$sokratesUrl = store.state.SokratesUrl
Vue.prototype.$herodotosUrl = store.state.HerodotosUrl
if (process.env.NODE_ENV === 'development') {
  Vue.prototype.$sokratesUrl = "http://minikube-lexiko.test/sokrates/v1"
  Vue.prototype.$herodotosUrl = "http://minikube-lexiko.test/herodotos/v1"
}

console.log(`sokrates proxy found at: ${Vue.prototype.$sokratesUrl}`)
console.log(`herodotos proxy found at: ${Vue.prototype.$herodotosUrl}`)


new Vue({
  el: '#app',
  render: h => h(App),
  router,
  template: '<App/>'
}).$mount('#app');
