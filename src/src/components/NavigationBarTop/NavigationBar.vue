<template>
    <el-tabs v-model="activeName" type="card" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="库房工作量统计" name="workload">
            <query-bar parent-name="workload" />
            <workload-data-show />
        </el-tab-pane>
        <el-tab-pane label="123" name="WorkLoadDataShow2">
            <query-bar parent-name="WorkLoadDataShow2" />
            <work-load-data-show2 />
        </el-tab-pane>
        <!-- <el-tab-pane label="调拨单未上账统计" name="departmentCollar">
            <query-bar parent-name="departmentCollar" />
            <depart-ment-collar-data-show />
        </el-tab-pane> -->
        <!-- <el-tab-pane label="计费未核对数据统计" name="UnCheckBills">
            <query-bar parent-name="UnCheckBills" />
            <un-check-bills-show />
        </el-tab-pane> -->
        <!-- <ElTabPane label="采购订单未到货统计" name="NoDeliveredPurchaseSummary">
            <query-bar parent-name="NoDeliveredPurchaseSummary" />
            <NoDeliveredPurchaseSummaryShow></NoDeliveredPurchaseSummaryShow>
        </ElTabPane> -->
        <!-- <el-tab-pane label="Task" name="fourth">Task</el-tab-pane> -->
    </el-tabs>
</template>
<script setup>
import { ref } from "vue";
import QueryBar from "../QueryBarTop/QueryBar.vue";
import WorkloadDataShow from "@/DataShowBar/WorkloadDataShow.vue";
// import DepartMentCollarDataShow from "../DataShowBar/DepartMentCollarDataShow.vue";
// import UnCheckBillsShow from "../DataShowBar/UnCheckBillsShow.vue";
// import NoDeliveredPurchaseSummaryShow from "../DataShowBar/NoDeliveredPurchaseSummaryShow.vue";
import WorkLoadDataShow2 from "@/workload/components/WorkloadDataShow2.vue";
import { ElTabPane } from "element-plus";

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
