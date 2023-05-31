<template>
    <el-tabs v-model="activeName" type="card" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="高值库房工作量统计" name="first">
            <query-bar />
            <data-show />
        </el-tab-pane>
        <!-- <el-tab-pane label="Config" name="second">Config</el-tab-pane>
        <el-tab-pane label="Role" name="third">Role</el-tab-pane>
        <el-tab-pane label="Task" name="fourth">Task</el-tab-pane> -->
    </el-tabs>
</template>
<script setup>
import { ref } from "vue";
import DataShow from "../DataShowBar/DataShow.vue";
import QueryBar from "../QueryBarTop/QueryBar.vue";

const activeName = ref("first");
const handleClick = () => {};
</script>
<style>
.demo-tabs > .el-tabs__content {
    padding-left: 32px;
    color: #6b778c;
    font-size: 32px;
    font-weight: 600;
}
</style>
