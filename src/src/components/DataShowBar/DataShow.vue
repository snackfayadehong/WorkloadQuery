<template>
    <el-table :data="tableData" style="width: 100%" class="elTable">
        <!-- 入库信息 -->
        <el-table-column prop="name" label="姓名" width="150" align="center" />
        <el-table-column label="入库信息" align="center">
            <el-table-column prop="In-billAmount" label="单据数" align="center" />
            <el-table-column prop="In-specAmount" label="品规数" align="center" />
            <el-table-column prop="In-totalAmount" label="总价格" align="center" />
        </el-table-column>
        <!-- 出库信息 -->
        <el-table-column label="出库信息" align="center">
            <el-table-column prop="Out-billAmount" label="单据数" align="center" />
            <el-table-column prop="Out-specAmount" label="品规数" align="center" />
            <el-table-column prop="Out-toalAmount" label="总价格" align="center" />
        </el-table-column>
        <!-- 退还信息 -->
        <el-table-column label="退还信息" align="center">
            <el-table-column prop="Refund-billAmount" label="单据数" align="center" />
            <el-table-column prop="Refund-specAmount" label="品规数" align="center" />
            <el-table-column prop="Refund-totalAmount" label="总价格" align="center" />
        </el-table-column>
        <!-- 合计信息 -->
        <el-table-column label="合计信息" align="center">
            <el-table-column prop="Total-billAmount" label="单据数" align="center" />
            <el-table-column prop="Total-specAmount" label="品规数" align="center" />
            <el-table-column prop="Total-totalAmount" label="总价格" align="center" />
        </el-table-column>
    </el-table>
</template>

<script setup>
const tableData = [];
</script>

<style scoped>
.elTable {
    --el-table-border-color: #a4b6e1;
    --el-table-border: 1px solid #a4b6e1;
    --el-table-header-text-color: black;
    --el-text-color-secondary: black;
}
</style>
