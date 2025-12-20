// src/router/menu.js
import {
    House,
    DataAnalysis,
    Tools,
    Menu,
    Collection
} from "@element-plus/icons-vue";

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
    },
    {
        path: '/tools',
        label: '辅助工具',    // 统一使用 label
        icon: Tools,        // 统一使用组件对象
        children: [
            {
                path: '/tools/hub', // 建议使用绝对路径，方便菜单高亮
                label: '工具概览',
                icon: Menu,
                component: "tools/ToolHub"
            },
            {
                path: '/tools/dict-compare',
                label: '字典对比',
                icon: Collection,
                component: "tools/DictCompareTool"
            }
        ]
    }
];