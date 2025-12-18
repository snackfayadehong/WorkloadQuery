/**
 * 菜单配置
 * 注意：
 * 1. label：菜单显示文字
 * 2. path：路由路径
 * 3. component：页面组件路径
 * 4. children：子菜单（可选）
 */

// src/router/menu.js
import { House, DataAnalysis, Document, Setting } from "@element-plus/icons-vue";

export default [
    {
        path: "/home",
        label: "首页",
        icon: House,
        component: "home/HomePage"
    },
    {
        path: "/workload",
        label: "工作量",
        icon: DataAnalysis,
        component: "workload/WorkloadPage"
    }
];
