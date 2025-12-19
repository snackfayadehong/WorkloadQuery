<template>
    <div class="sidebar-wrapper" :class="{ 'is-collapsed': collapse }">
        <div class="logo-container">
            <transition name="logo-fade">
                <div v-if="!collapse" class="logo-full" key="full">
                    <el-icon class="logo-icon"><Monitor /></el-icon>
                    <span class="logo-text">SupperSystem</span>
                </div>
                <div v-else class="logo-mini" key="mini">
                    <span class="logo-text-short">SyS</span>
                </div>
            </transition>
        </div>

        <el-scrollbar class="menu-scrollbar">
            <el-menu
                router
                :default-active="$route.path"
                :collapse="collapse"
                :collapse-transition="false"
                background-color="#001529"
                text-color="#bfcbd9"
                active-text-color="#ffffff"
                class="sidebar-menu"
            >
                <template v-for="item in menu" :key="item.path">
                    <el-sub-menu v-if="item.children?.length" :index="item.path">
                        <template #title>
                            <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
                            <span>{{ item.label }}</span>
                        </template>
                        <el-menu-item v-for="child in item.children" :key="child.path" :index="child.path">
                            {{ child.label }}
                        </el-menu-item>
                    </el-sub-menu>

                    <el-menu-item v-else :index="item.path">
                        <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
                        <template #title>
                            <span>{{ item.label }}</span>
                        </template>
                    </el-menu-item>
                </template>
            </el-menu>
        </el-scrollbar>
    </div>
</template>

<script setup>
import { Monitor } from "@element-plus/icons-vue";
import menu from "@/router/menu"; //

defineProps({
    collapse: Boolean
});
</script>

<style scoped>
/* 基础容器优化 */
.sidebar-wrapper {
    height: 100%;
    background: #001529;
    display: flex;
    flex-direction: column;
    transition: width 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
    box-shadow: 2px 0 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
}

/* Logo 样式增强 */
.logo-container {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 16px;
    overflow: hidden;
    white-space: nowrap;
    background: rgba(255, 255, 255, 0.02);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.logo-full {
    display: flex;
    align-items: center;
    gap: 12px;
}

.logo-icon {
    font-size: 24px;
    color: #409eff;
}

.logo-text {
    font-size: 18px;
    font-weight: 800;
    color: #ffffff;
    letter-spacing: 1px;
}

.logo-mini {
    font-size: 20px;
    font-weight: 800;
    color: #409eff;
}

/* 菜单样式微调 */
.menu-scrollbar {
    flex: 1;
}

.sidebar-menu {
    border-right: none !important;
}

/* 菜单项激活态光效 */
:deep(.el-menu-item.is-active) {
    background: linear-gradient(90deg, rgba(64, 158, 255, 0.2) 0%, rgba(64, 158, 255, 0) 100%) !important;
    position: relative;
}

:deep(.el-menu-item.is-active)::before {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: #409eff;
    box-shadow: 2px 0 8px rgba(64, 158, 255, 0.5);
}

/* 悬停效果 */
:deep(.el-menu-item:hover) {
    color: #ffffff !important;
}

/* Logo 切换动画 */
.logo-fade-enter-active,
.logo-fade-leave-active {
    transition: all 0.2s ease;
}

.logo-fade-enter-from,
.logo-fade-leave-to {
    opacity: 0;
    transform: scale(0.8);
}

/* 处理折叠时文字消失的平滑度 */
:deep(.el-menu--collapse .el-sub-menu__title span),
:deep(.el-menu--collapse .el-menu-item span) {
    height: 0;
    width: 0;
    overflow: hidden;
    visibility: hidden;
    display: inline-block;
}
</style>
