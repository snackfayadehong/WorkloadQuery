<template>
    <el-container style="height: 100vh">
        <!-- 顶部导航栏 -->
        <el-header height="56px" class="header">
            <TopBar @toggle-menu="isCollapse = !isCollapse" />
        </el-header>

        <el-container>
            <!-- 左侧菜单 -->
            <el-aside class="sidebar" :class="{ collapsed: isCollapse }">
                <Sidebar :collapse="isCollapse" />
            </el-aside>

            <!-- 主内容区 -->
            <el-main class="content">
                <div class="page-container">
                    <router-view />
                </div>
            </el-main>
        </el-container>
    </el-container>
</template>

<script>
import TopBar from "./components/TopBar.vue";
import Sidebar from "./components/Sidebar.vue";

export default {
    name: "MainLayout",
    components: { TopBar, Sidebar },
    data() {
        return {
            isCollapse: false
        };
    }
};
</script>
