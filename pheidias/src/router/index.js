import Vue from "vue";
import Router from "vue-router";

const routerOptions = [
    { path: "/", view: "HomePage" },
    { path: "/multi", view: "Sokrates" },
    { path: "/texts", view: "Herodotos" },
    { path: "/search", view: "Alexandros" },
    { path: "*", view: "NotFound" }
];

const routes = routerOptions.map(route => {
    return {
        ...route,
        component: () => import(`../views/${route.view}.vue`)
    };
});

Vue.use(Router);

export default new Router({
    mode: "history",
    routes
});