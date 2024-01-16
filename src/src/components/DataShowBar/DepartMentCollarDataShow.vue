<template>
    <el-table :data="tableData" style="width: 100%" :border="true" class="elTable">
        <el-table-column fixed prop="Index" label="序号" width="50" align="center" />
        <el-table-column prop="DepartmentCollarCode" label="领用单号" align="center" />
        <el-table-column prop="LeadingDepartmentName" label="领用科室" align="center" />
        <el-table-column prop="LeaderName" label="领用人" align="center" />
        <el-table-column prop="BLMakerName" label="制单人" align="center" />
        <el-table-column prop="CreateTime" label="总库出库时间" align="center" />
        <el-table-column prop="StatisticalTime" label="数据统计时间" align="center" />
        <el-table-column prop="Code" label="材料代码" align="center" />
        <el-table-column prop="ProdName" label="材料名称" align="center" />
        <el-table-column prop="SpecModelName" label="型号|规格" align="center" />
        <el-table-column prop="UnitName" label="单位" align="center" />
        <el-table-column prop="ChargePrice" label="单价" align="center" />
        <el-table-column prop="Amount" label="数量" align="center" />
    </el-table>
</template>

<script setup name="departmentCollar">
import { ref } from "vue";
import bus from "../../eventBus";
// 表格数据
const tableData = ref([]);
bus.on("getDepartmenCollarData", res => {
    tableData.value = res;
});
</script>

<style scoped></style>
