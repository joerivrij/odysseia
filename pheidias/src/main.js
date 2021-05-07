import Vue from 'vue'
import App from './App.vue'
import axios from 'axios';
import router from './router'
import store from './store'

Vue.config.productionTip = false

Vue.http = Vue.prototype.$apiClient = axios;
Vue.prototype.$proxyUrl = store.state.proxyUrl + "/sokrates/v1"
if (process.env.NODE_ENV === 'development') {
  Vue.prototype.$proxyUrl = "http://minikube-lexiko.test/sokrates/v1"
}

console.log(`sokrates proxy found at: ${Vue.prototype.$proxyUrl}`)

new Vue({
  el: '#app',
  render: h => h(App),
  router,
  template: '<App/>'
}).$mount('#app');
