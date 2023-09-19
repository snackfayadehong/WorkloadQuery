<template>
    <div class="demo-date-picker">
        <div class="block">
            <p class="dateTime-select">时间区间:</p>
            <el-config-provider>
                <el-date-picker v-model="queryDate" type="daterange" value-format="YYYY-MM-DD" start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-config-provider>
        </div>
        <el-button type="primary" class="qry-btn" @click="query">查询</el-button>
        <el-button v-if="props.parentName == 'NoDeliveredPurchaseSummary'" type="primary" class="qry-btn" @click="exportExcel">导出</el-button>
    </div>
</template>

<script setup>
import myAxios from "../../api/myAxios.js";
import { ref } from "vue";
import bus from "../../eventBus";
import dayjs from "dayjs";
import ExportExcelUtity from "../../utity/exportExcel";
// import { getCurrentInstance } from "vue";
// 接口返回信息
let res = "";

const props = defineProps({
    // 父组件名
    parentName: {
        type: String,
        default() {
            return "";
        }
    }
});

// 格式化时间
let end = dayjs().format("YYYY-MM-DD");
let start = dayjs().subtract(6, "day").format("YYYY-MM-DD");
// 时间
const queryDate = ref("");
queryDate.value = [start, end];

// 通过 getCurrentInstance 方法获取当前组件实例 并通过.parent.proxy.name获取父组件名称
// let commName = "";
// let currentCpn = getCurrentInstance();

// 查询工作量数据
const query = async () => {
    // 获取父组件名称,根据夫组件名称查询不同数据
    // commName = currentCpn.parent.proxy.name;
    // console.log(commName);

    let Timevalue = JSON.parse(JSON.stringify(queryDate.value));
    if (props.parentName == "") {
        alert("组件名为空");
        return;
    }
    if (props.parentName == "workload") {
        // 工作量查询
        res = await myAxios.post("/getWorkload", { startTime: Timevalue[0], endTime: Timevalue[1] });
        if (res == "") {
            alert("无数据");
            return;
        }
        bus.emit("getWorkloadData", res.Data);
    } else if (props.parentName == "departmentCollar") {
        // 未上账单据查询
        res = await myAxios.post("/getNoAccountEntry", { startTime: Timevalue[0], endTime: Timevalue[1] });
        if (res == "") {
            alert("无数据");
            return;
        }
        bus.emit("getDepartmenCollarData", res.Data);
    } else if (props.parentName == "UnCheckBills") {
        // 未核对数据查询
        res = await myAxios.post("/getUnCheckBills", { startTime: Timevalue[0], endTime: Timevalue[1] });
        if (res == "") {
            alert("无数据");
            return;
        }
        bus.emit("getUnCheckBills", res.Data);
    } else if (props.parentName == "NoDeliveredPurchaseSummary") {
        // 采购订单未到货数据统计
        res = await myAxios.post("/getNoDeliveredPurchaseSummary", { startTime: Timevalue[0], endTime: Timevalue[1] });
        if (res == "") {
            alert("无数据");
            return;
        }
        bus.emit("getNodeliveredPurchaseSummary", res.Data);
    }
};
const exportExcel = () => {
    if (typeof res.Data == "undefined" || res.Data.length == 0) {
        alert("请先查询信息");
        reutrn;
    }
    let now = new Date();
    let ExcelName = "采购订单未到货信息_" + dayjs(now).format("YYYY-MM-DD_HH:mm:ss") + ".xlsx";
    ExportExcelUtity(res.Data, ExcelName, "NodeliveredPurchase");
};
</script>
<style scoped>
.dateTime-select {
    margin: 0;
}
.demo-date-picker {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 0;
    flex-wrap: wrap;
    font-size: 14px;
}
.demo-date-picker .block {
    padding-bottom: 15px;
    text-align: left;
    /* border-right: solid 1px var(--el-border-color); */
    flex: 1;
}
.demo-date-picker .qry-btn {
    margin-right: 80px;
}
.demo-date-picker .block:last-child {
    border-right: none;
}
</style>
