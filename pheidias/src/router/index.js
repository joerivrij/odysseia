import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
        {
            path: '/',
            name: 'homepage',
            component: () => import("../views/HomePage")
        },
        {
            path: '/herodotos/',
            name: 'sentence-area',
            component: () => import("../views/Herodotos")
        },
    ]
})