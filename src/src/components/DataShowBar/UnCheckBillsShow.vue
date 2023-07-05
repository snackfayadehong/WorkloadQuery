<template>
    <el-table :data="tableData" style="width: 100%" :border="true" class="elTable">
        <el-table-column fixed prop="Index" label="序号" width="50" align="center" />
        <el-table-column prop="BillNo" label="单号" align="center" />
        <el-table-column prop="PatName" label="制单人" align="center" />
        <el-table-column prop="Section" label="领用科室" align="center" />
        <el-table-column prop="Doctor" label="领用人" align="center" />
        <el-table-column prop="OperateName" label="出库科室" align="center" />
        <el-table-column prop="CheckDate" label="库房审核时间" align="center" />
        <el-table-column prop="CheckStatus" label="状态" align="center" />
    </el-table>
</template>

<script setup name="departmentCollar">
import { ref } from "vue";
import bus from "../../eventBus";
// 表格数据
const tableData = ref([]);
bus.on("getUnCheckBills", res => {
    tableData.value = res;
});
</script>

<style scoped></style>
