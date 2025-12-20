<template>
    <div class="workload-dashboard-page">
        <header class="hero-section">
            <div class="hero-content">
                <h1 class="welcome-text">工作量看板</h1>
                <p class="date-text">下面这个按钮可以导出全部QAQ</p>
                <div class="header-actions">
                    <el-button type="primary" :icon="Download" class="modern-action-btn" :loading="exportLoading"
                        @click="handleExport"> 导出报表数据 </el-button>
                </div>
            </div>
        </header>

        <main class="dashboard-grid">
            <WorkloadStats :data="rawData" />

            <WorkloadFilter v-model:search="searchQuery" v-model:type="filterType" v-model:dateRange="dateRange"
                :total="filteredData.length" @query="fetchList" />

            <el-card shadow="never" class="glass-panel table-card">
                <WorkloadTable :data="paginatedData" :loading="loading" @view-detail="handleViewDetail"
                    @export-row="handleExportRow" />

                <div class="pagination-container">

                    <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                        :total="filteredData.length" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next"
                        background />
                </div>
            </el-card>
        </main>

        <el-dialog v-model="detailVisible" :title="`业务处理明细汇总 - ${currentDetail?.operator || ''}`" width="1000px"
            destroy-on-close append-to-body class="modern-dialog">
            <WorkloadDetail :data="currentDetail" @close="detailVisible = false"
                @export-current="handleExportCurrent" />
        </el-dialog>
    </div>
</template>

<script setup>
import { ref } from "vue";
import { Download } from "@element-plus/icons-vue";
import { useWorkload } from "./composables/useWorkload";
import { ElMessage } from "element-plus";
import dayjs from "dayjs";
import ExportExcelUtity from "@/utity/exportExcel";

import WorkloadStats from "./components/WorkloadStats.vue";
import WorkloadFilter from "./components/WorkloadFilter.vue";
import WorkloadTable from "./components/WorkloadTable.vue";
import WorkloadDetail from "./components/WorkloadDetail.vue";

const { rawData, loading, searchQuery, filterType, dateRange, currentPage, pageSize, filteredData, paginatedData, fetchList } = useWorkload();

const detailVisible = ref(false);
const currentDetail = ref({});
const exportLoading = ref(false);

const handleViewDetail = row => {
    currentDetail.value = row;
    detailVisible.value = true;
};

const performExport = (data, fileName) => {
    try {
        ExportExcelUtity(data, fileName, "Workload");
        ElMessage.success(`${fileName} 导出成功`);
    } catch (error) {
        ElMessage.error("导出失败：" + error.message);
    }
};

const handleExport = async () => {
    if (filteredData.value.length === 0) {
        ElMessage.warning("当前筛选条件下无数据可供导出");
        return;
    }
    exportLoading.value = true;
    try {
        const fileName = `全院工作量汇总_${dayjs().format('YYYYMMDD_HHmmss')}.xlsx`;
        performExport(filteredData.value, fileName);
    } finally {
        exportLoading.value = false;
    }
};

const handleExportRow = (row) => {
    const fileName = `工作量明细_${row.operator}_${dayjs().format('YYYYMMDD')}.xlsx`;
    performExport(row, fileName);
};

const handleExportCurrent = (data) => {
    const fileName = `工作量明细_${data.operator}_${dayjs().format('YYYYMMDD')}.xlsx`;
    performExport(data, fileName);
};
</script>

<style scoped>
.workload-dashboard-page {
    /* 复用 HomePage 变量 */
    --bg-page: #f6f8fb;
    --bg-card: #ffffff;
    --text-primary: #262626;
    --card-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
    --glass-bg: rgba(255, 255, 255, 0.7);
    --glass-border: rgba(255, 255, 255, 0.4);
    --hero-gradient: linear-gradient(135deg, #4158d0 0%, #c850c0 46%, #ffcc70 100%);

    padding: 32px;
    background: var(--bg-page);
    min-height: 100%;
    display: flex;
    flex-direction: column;
    gap: 32px;
    box-sizing: border-box;
}

:global(html.dark) .workload-dashboard-page {
    --bg-page: #121212;
    --bg-card: #1e1e1e;
    --text-primary: #e5eaf3;
    --hero-gradient: linear-gradient(135deg, #1f2d5a 0%, #6d2d6d 46%, #a37c32 100%);
}

/* Hero Section */
.hero-section {
    padding: 48px;
    background: var(--hero-gradient);
    border-radius: 24px;
    color: white;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.welcome-text {
    font-size: 36px;
    font-weight: 800;
    margin: 0;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    /* 提升质感 */
}

.date-text {
    font-size: 16px;
    opacity: 0.8;
    margin: 12px 0 24px;
}

.modern-action-btn {
    background: rgba(255, 255, 255, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.3);
    color: white;
    backdrop-filter: blur(8px);
    transition: all 0.3s;
}

.modern-action-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    border-color: white;
}

.dashboard-grid {
    display: flex;
    flex-direction: column;
    gap: 32px;
}

/* 玻璃拟态面板 */
.glass-panel {
    background: var(--glass-bg) !important;
    backdrop-filter: blur(12px);
    border: 1px solid var(--glass-border) !important;
    border-radius: 20px;
}

.pagination-container {
    margin-top: 24px;
    display: flex;
    justify-content: flex-end;
}
</style>