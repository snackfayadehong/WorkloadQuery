/**
 * 菜单配置
 * 注意：
 * 1. label：菜单显示文字
 * 2. path：路由路径
 * 3. component：页面组件路径
 * 4. children：子菜单（可选）
 */

export default [
    {
        label: "首页",
        path: "/home",
        component: "home/HomePage"
    },
    {
        label: "工作量",
        path: "/workload",
        component: "workload/WorkloadPage"
    },
    {
        label: "采购管理",
        children: [
            {
                label: "采购汇总",
                path: "/purchase",
                component: "purchase/PurchaseSummaryPage"
            }
        ]
    }
];
