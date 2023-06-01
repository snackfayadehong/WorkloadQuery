<template>
    <div class="demo-date-picker">
        <div class="block">
            <p class="dateTime-select">时间区间:</p>
            <el-config-provider>
                <el-date-picker v-model="queryDate" type="daterange" value-format="YYYY-MM-DD" start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-config-provider>
        </div>
        <el-button type="primary" class="qry-btn" @click="query">查询</el-button>
    </div>
</template>

<script setup>
import myAxios from "../../api/myAxios.js";
import { ref } from "vue";

const queryDate = ref("");

const query = () => {};
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
