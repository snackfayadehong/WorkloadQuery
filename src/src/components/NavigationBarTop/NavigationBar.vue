<template>
    <el-tabs v-model="activeName" type="card" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="高值库房工作量统计" name="workload">
            <query-bar parent-name="workload" />
            <workload-data-show />
        </el-tab-pane>
        <el-tab-pane label="调拨单未上账统计" name="departmentCollar">
            <query-bar parent-name="departmentCollar" />
            <depart-ment-collar-data-show />
        </el-tab-pane>
        <el-tab-pane label="计费未核对数据统计" name="third">
            <query-bar parent-name="UnCheckBills" />
            <un-check-bills-show />
        </el-tab-pane>
        <!-- <el-tab-pane label="Task" name="fourth">Task</el-tab-pane> -->
    </el-tabs>
</template>
<script setup>
import { ref } from "vue";
import QueryBar from "../QueryBarTop/QueryBar.vue";
import WorkloadDataShow from "../DataShowBar/WorkloadDataShow.vue";
import DepartMentCollarDataShow from "../DataShowBar/DepartMentCollarDataShow.vue";
import UnCheckBillsShow from "../DataShowBar/UnCheckBillsShow.vue";

const activeName = ref("workload");
const handleClick = () => {};
</script>
<style>
.demo-tabs > .el-tabs__content {
    padding-left: 32px;
    font-size: 32px;
    font-weight: 600;
}
.elTable {
    font-size: 16px;
    --el-table-border-color: #a4b6e1;
    --el-table-border: 1px solid #a4b6e1;
    --el-table-header-text-color: black;
    --el-text-color-secondary: black;
}
</style>
