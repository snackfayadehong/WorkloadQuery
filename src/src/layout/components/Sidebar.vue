<template>
    <div class="sidebar-wrapper">
        <!-- Logo -->
        <div class="logo">
            <span class="logo-text" :class="{ hide: collapse }"> WorkloadQuery </span>
            <span class="logo-short" :class="{ show: collapse }"> WQ </span>
        </div>

        <el-menu router :default-active="$route.path" :collapse="collapse" background-color="#001529" text-color="#bfcbd9" active-text-color="#409eff" class="sidebar-menu">
            <template v-for="item in menu" :key="item.path">
                <!-- 有子菜单 -->
                <el-sub-menu v-if="item.children && item.children.length" :index="item.path">
                    <template #title>
                        <el-icon v-if="item.icon">
                            <component :is="item.icon" />
                        </el-icon>
                        <span>{{ item.label }}</span>
                    </template>

                    <el-menu-item v-for="child in item.children" :key="child.path" :index="child.path">
                        {{ child.label }}
                    </el-menu-item>
                </el-sub-menu>

                <!-- 无子菜单 -->
                <el-menu-item v-else :index="item.path">
                    <el-icon v-if="item.icon">
                        <component :is="item.icon" />
                    </el-icon>
                    <span>{{ item.label }}</span>
                </el-menu-item>
            </template>
        </el-menu>
    </div>
</template>

<script>
import menu from "@/router/menu";

export default {
    name: "Sidebar",
    props: {
        collapse: Boolean
    },
    data() {
        return {
            menu
        };
    }
};
</script>

<style>
/* ======================
   基础结构
====================== */

.sidebar-wrapper {
    height: 100%;
    background: #001529;
}

/* ======================
   Logo 动画
====================== */

.logo {
    height: 56px;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    color: #fff;
    font-weight: bold;
}

.logo-text,
.logo-short {
    position: absolute;
    transition: all 0.2s ease;
}

.logo-text {
    opacity: 1;
    transform: translateX(0);
}

.logo-text.hide {
    opacity: 0;
    transform: translateX(-10px);
}

.logo-short {
    opacity: 0;
    transform: scale(0.8);
}

.logo-short.show {
    opacity: 1;
    transform: scale(1);
}

/* ======================
   菜单样式
====================== */

.sidebar-menu {
    border-right: none;
}

/* ======================
   菜单文字动画（关键）
   ❗ 不使用 display:none
====================== */

.el-menu-item span,
.el-sub-menu__title span {
    display: inline-block;
    white-space: nowrap;
    overflow: hidden;
    width: 120px;
    opacity: 1;
    transition: opacity 0.2s ease, width 0.2s ease;
}

/* 折叠状态 */
.el-menu--collapse .el-menu-item span,
.el-menu--collapse .el-sub-menu__title span {
    width: 0;
    opacity: 0;
}
</style>
