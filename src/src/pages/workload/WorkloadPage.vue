<template>
    <div class="workload-page">
        <el-card shadow="never" class="custom-page-header">
            <div class="header-flex">
                <div class="title-group">
                    <h1 class="main-title">工作量看板</h1>
                    <span class="sub-title">请选择时间段后点击查询以获取最新数据</span>
                </div>
                <div class="header-actions">
                    <el-button type="primary" :icon="Download" plain :loading="exportLoading" @click="handleExport"> 导出报表 </el-button>
                </div>
            </div>
        </el-card>

        <main class="dashboard-content">
            <WorkloadStats :data="rawData" />

            <WorkloadFilter v-model:search="searchQuery" v-model:type="filterType" v-model:dateRange="dateRange" :total="filteredData.length" @query="fetchList" />

            <el-card shadow="never" class="table-card">
                <WorkloadTable :data="paginatedData" :loading="loading" @view-detail="handleViewDetail" />

                <div class="pagination-container">
                    <el-pagination
                        v-model:current-page="currentPage"
                        v-model:page-size="pageSize"
                        :total="filteredData.length"
                        :page-sizes="[10, 20, 50]"
                        layout="total, sizes, prev, pager, next"
                        background
                    />
                </div>
            </el-card>
        </main>

        <WorkloadDetail v-model="detailVisible" :data="currentDetail" />
    </div>
</template>

<script setup>
import { ref } from "vue";
import { Download } from "@element-plus/icons-vue";
import { useWorkload } from "./composables/useWorkload";
import { ElMessage } from "element-plus";

// 引入子组件
import WorkloadStats from "./components/WorkloadStats.vue";
import WorkloadFilter from "./components/WorkloadFilter.vue";
import WorkloadTable from "./components/WorkloadTable.vue";
import WorkloadDetail from "./components/WorkloadDetail.vue";

// 从 Composable 导出逻辑
const { rawData, loading, searchQuery, filterType, dateRange, currentPage, pageSize, filteredData, paginatedData, fetchList } = useWorkload();

const detailVisible = ref(false);
const currentDetail = ref({});
const exportLoading = ref(false);

/**
 * 处理查看详情：将数据传递给弹窗并打开
 */
const handleViewDetail = row => {
    currentDetail.value = row;
    detailVisible.value = true;
};

/**
 * 处理报表导出逻辑
 */
const handleExport = async () => {
    if (filteredData.value.length === 0) {
        ElMessage.warning("当前无数据可导出");
        return;
    }
    exportLoading.value = true;
    try {
        // 这里未来对接后端的导出接口
        console.log("正在导出数据...", filteredData.value);
        await new Promise(resolve => setTimeout(resolve, 1500)); // 模拟请求
        ElMessage.success("报表生成成功，开始下载...");
    } finally {
        exportLoading.value = false;
    }
};
</script>

<style scoped>
.workload-page {
    /* 关键修复：使用变量适配暗黑模式 */
    background-color: var(--el-bg-color-page);
    padding: 24px;
    min-height: 100%;
    display: flex;
    flex-direction: column;
    gap: 20px;
    /* 解决滚动条冲突的核心 */
    box-sizing: border-box;
}

.custom-page-header {
    border-radius: 12px;
    border: none;
    background-color: var(--el-bg-color); /* 适配暗黑模式背景 */
    transition: background-color 0.3s;
}

.header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.main-title {
    margin: 0;
    font-size: 20px;
    color: var(--el-text-color-primary); /* 使用变量适配文字颜色 */
    font-weight: 700;
}

.sub-title {
    font-size: 13px;
    color: var(--el-text-color-secondary);
    margin-top: 4px;
    display: block;
}

.dashboard-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
    flex: 1;
}

.table-card {
    border-radius: 12px;
    border: none;
    background-color: var(--el-bg-color);
}

.pagination-container {
    margin-top: 24px;
    display: flex;
    justify-content: flex-end;
}

/* 针对暗黑模式下卡片边框的微调 */
:global(html.dark) .custom-page-header,
:global(html.dark) .table-card {
    border: 1px solid var(--el-border-color-lighter);
}
</style>
