<template>
    <el-container style="height: 100vh">
        <el-aside width="220px">
            <el-menu router>
                <template v-for="item in menu" :key="item.label">
                    <el-sub-menu v-if="item.children" :index="item.label">
                        <template #title>{{ item.label }}</template>

                        <el-menu-item v-for="child in item.children" :key="child.path" :index="child.path">
                            {{ child.label }}
                        </el-menu-item>
                    </el-sub-menu>

                    <el-menu-item v-else :index="item.path">
                        {{ item.label }}
                    </el-menu-item>
                </template>
            </el-menu>
        </el-aside>

        <el-container>
            <el-header>
                <NavigationBar />
            </el-header>

            <el-main>
                <router-view />
            </el-main>
        </el-container>
    </el-container>
</template>

<script>
import menu from "@/router/menu";
import NavigationBar from "@/components/NavigationBarTop/NavigationBar.vue";

export default {
    name: "MainLayout",
    components: { NavigationBar },
    data() {
        return {
            menu
        };
    }
};
</script>
