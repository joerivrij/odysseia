import Vue from 'vue'
import App from './App.vue'
import axios from 'axios';
import router from './router'
import store from './store'

Vue.config.productionTip = false

Vue.http = Vue.prototype.$apiClient = axios;
Vue.prototype.$proxyUrl = store.state.proxyUrl + "/api/v1"

new Vue({
  el: '#app',
  render: h => h(App),
  router,
  template: '<App/>'
}).$mount('#app');
