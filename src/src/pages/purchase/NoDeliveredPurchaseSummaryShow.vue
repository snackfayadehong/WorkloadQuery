<!-- <template>
    <div>
        <el-table :data="purchaseSummary" lazy="true" :load="load" style="width: 100%; margin-bottom: 20px" row-key="PurchaseSummaryID" border default-expand-all>
            <el-table-column prop="PurchaseSummaryID" label="订单ID" sortable />
            <el-table-column prop="PurchaseSummaryCode" label="订单号" sortable />
            <el-table-column prop="AllMoney" label="订单金额" sortable />
            <el-table-column prop="DepartmentName" label="入库科室" sortable />
            <el-table-column prop="StatusName" label="订单状态" sortable />
            <el-table-column prop="MakeName" label="开单人员" sortable />
            <el-table-column prop="AuditorDate" label="订单日期" sortable />
            <el-table-column prop="GoodStatusName" label="订单物流状态" sortable />
            <el-table-column prop="SupplierName" label="供货公司" sortable />
            <el-table-column prop="Remark" label="订单备注" sortable />
        </el-table>
    </div>
</template> -->
<template>
    <el-table :data="purchaseSummary" row-key="PurchaseSummaryID" style="width: 100%">
        <el-table-column type="expand">
            <template #default="Props">
                <el-table :data="Props.row.children" lazy :load="load">
                    <el-table-column label="院内代码" prop="Code" />
                    <el-table-column label="产品名称" prop="ProdName" />
                    <el-table-column label="院内规格" prop="HospitalSpec" />
                    <el-table-column label="供货公司" prop="EnterpriseName" />
                    <el-table-column label="单价" prop="UnitPrice" />
                    <el-table-column label="规格" prop="SpecName" />
                    <el-table-column label="数量" prop="Amount" />
                    <el-table-column label="到货数量" prop="FactInAmount" />
                    <el-table-column label="关闭数量" prop="RefundAmount" />
                    <el-table-column label="未到货数量" prop="NotDelivered" />
                    <el-table-column label="明细备注" prop="Remark" />
                </el-table>
            </template>
        </el-table-column>
        <el-table-column prop="PurchaseSummaryID" label="订单ID" sortable />
        <el-table-column prop="PurchaseSummaryCode" label="订单号" sortable />
        <el-table-column prop="AllMoney" label="订单金额" sortable />
        <el-table-column prop="DepartmentName" label="入库科室" sortable />
        <el-table-column prop="StatusName" label="订单状态" sortable />
        <el-table-column prop="MakeName" label="开单人员" sortable />
        <el-table-column prop="AuditorDate" label="订单日期" sortable />
        <el-table-column prop="GoodStatusName" label="订单物流状态" sortable />
        <el-table-column prop="SupplierName" label="供货公司" sortable />
        <el-table-column prop="Remark" label="订单备注" sortable />
    </el-table>
</template>
<script setup>
import bus from "../../eventBus";
import { ref } from "vue";

var purchaseSummary = ref([]);
bus.on("getNodeliveredPurchaseSummary", res => {
    if (typeof res == "undefined") {
        alert("无数据");
        return;
    }
    purchaseSummary.value = res;
});
const load = () => {
    console.log("123");
};
</script>
