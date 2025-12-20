// src/router/index.js
import { createRouter, createWebHistory } from "vue-router";
import MainLayout from "@/layout/MainLayout.vue";
import menuList from "./menu";
import NProgress from "nprogress"; // 引入进度条
import "nprogress/nprogress.css"; // 引入进度条样式

// 配置 NProgress (去掉右上角的螺旋加载圈)
NProgress.configure({ showSpinner: false });

const pages = import.meta.glob("../pages/**/*.vue");

function loadView(view) {
    if (!view) return null; // 过滤调没有组件的父级菜单
    const key = `../pages/${view}.vue`;
    const loader = pages[key];
    if (!loader) {
        console.error(`❌ 页面不存在: ${key}`);
    }
    return loader;
}

/**
 * 递归生成路由配置
 * 将嵌套的 menu 结构平铺或嵌套转换成 Vue Router 识别的格式
 */
const formatRoutes = (menu) => {
    let result = [];
    menu.forEach(item => {
        // 如果有组件，则注册路由
        if (item.component && typeof item.component === 'string') {
            result.push({
                path: item.path,
                name: item.name || item.path,
                component: loadView(item.component),
                meta: { title: item.label }
            })
        }
        // 如果有子菜单，则递归处理
        if (item.children && item.children.length > 0) {
            result.push(...formatRoutes(item.children));
        }
    })
    return result

}

const routes = [
    {
        path: "/",
        component: MainLayout,
        redirect: "/home",
        children: [
            ...formatRoutes(menuList), // 使用递归后的平铺路由
            {
                path: "/:pathMatch(.*)*",
                name: "NotFound",
                component: () => import("@/pages/error/NotFound.vue"),
                meta: { title: "404 - 页面不存在" }
            }
        ]
    }
];
const router = createRouter({
    history: createWebHistory(),
    routes,
    // 优化 1: 切换路由时自动滚动到页面顶部
    scrollBehavior() {
        return { top: 0 };
    }
});

// 优化 2: 增加前置守卫开启进度条
router.beforeEach((to, from, next) => {
    NProgress.start();
    next();
});

router.afterEach((to) => {
    // 优化 3: 进度条结束
    NProgress.done();
    document.title = to.meta.title ? `${to.meta.title} | WorkloadQuery` : 'WorkloadQuery';
});

export default router;