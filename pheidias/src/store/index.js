import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        SokratesUrl: document.location.origin + "/sokrates/v1",
        HerodotosUrl: document.location.origin + "/herodotos/v1",
        AlexandrosUrl: document.location.origin + "/alexandros/v1",
        DionysiosUrl: document.location.origin + "/dionysios/v1"
    },
})