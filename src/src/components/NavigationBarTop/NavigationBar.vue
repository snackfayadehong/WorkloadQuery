<template>
    <el-tabs v-model="activeName" type="card" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="高值库房工作量统计" name="workload">
            <query-bar />
            <workload-data-show />
        </el-tab-pane>
        <el-tab-pane label="调拨单未上账统计" name="departmentCollar">
            <query-bar />
            <depart-ment-collar-data-show />
        </el-tab-pane>
        <el-tab-pane label="计费未核对数据统计" name="third">计费未核对数据统计</el-tab-pane>
        <!-- <el-tab-pane label="Task" name="fourth">Task</el-tab-pane> -->
    </el-tabs>
</template>
<script setup>
import { ref } from "vue";
import QueryBar from "../QueryBarTop/QueryBar.vue";
import WorkloadDataShow from "../DataShowBar/WorkloadDataShow.vue";
import DepartMentCollarDataShow from "../DataShowBar/DepartMentCollarDataShow.vue";

const activeName = ref("workload");
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
