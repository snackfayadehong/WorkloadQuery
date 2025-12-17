<template>
    <el-table :data="dataList" style="width: 100%" class="elTable">
        <!-- 入库信息 -->
        <el-table-column prop="name" label="姓名" width="150" align="center" />
        <el-table-column label="入库信息" align="center">
            <el-table-column prop="prodAcBill" label="单据数" align="center" />
            <el-table-column prop="prodAcSpec" label="品规数" align="center" />
            <el-table-column prop="prodAcTotal" label="总价格" align="center" />
        </el-table-column>
        <!-- 出库信息 -->
        <el-table-column label="出库信息" align="center">
            <el-table-column prop="prodDpBill" label="单据数" align="center" />
            <el-table-column prop="prodDpSpec" label="品规数" align="center" />
            <el-table-column prop="prodDpTotal" label="总价格" align="center" />
        </el-table-column>
        <!-- 退还信息 -->
        <el-table-column label="退还信息" align="center">
            <el-table-column prop="refBill" label="单据数" align="center" />
            <el-table-column prop="refSpec" label="品规数" align="center" />
            <el-table-column prop="refTotal" label="总价格" align="center" />
        </el-table-column>
        <!-- 合计信息 -->
        <el-table-column label="合计信息" align="center">
            <el-table-column prop="total_bill_amount" label="单据数" align="center"> </el-table-column>
            <el-table-column prop="total_spec_amount" label="品规数" align="center" />
            <el-table-column prop="total_total_amount" label="总价格" align="center" />
        </el-table-column>
    </el-table>
</template>

<script setup>
import bus from "../../eventBus";
import { ref } from "vue";
import _ from "lodash";

var dataList = ref([]);

bus.on("getWorkloadData", res => {
    if (typeof res == "undefined") {
        alert("无数据");
        return;
    }
    dataList.value = res;
});
</script>

<style scoped></style>
