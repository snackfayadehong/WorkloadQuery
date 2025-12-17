import { createRouter, createWebHistory } from "vue-router";
import MainLayout from "../layout/MainLayout.vue";
import menuList from "./menu";

function loadView(view) {
    return () => import(`@/pages/${view}.vue`);
}

function buildRoutes(menu) {
    const routes = [];

    menu.forEach(item => {
        if (item.children && item.children.length > 0) {
            routes.push(...buildRoutes(item.children));
        } else if (item.path && item.component) {
            routes.push({
                path: item.path,
                component: loadView(item.component)
            });
        }
    });

    return routes;
}

const routes = [
    {
        path: "/",
        component: MainLayout,
        redirect: "/home",
        children: buildRoutes(menuList)
    }
];

export default createRouter({
    history: createWebHistory(),
    routes
});
