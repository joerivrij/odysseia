import Vue from 'vue'
import App from './App.vue'
import axios from 'axios';
import router from './router'

Vue.config.productionTip = false

Vue.http = Vue.prototype.$apiClient = axios;
Vue.prototype.$proxyUrl = "http://localhost:5000/api/v1";

new Vue({
  el: '#app',
  render: h => h(App),
  router,
  template: '<App/>'
}).$mount('#app');
