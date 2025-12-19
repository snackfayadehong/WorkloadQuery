<template>
    <div class="workload-page">
        <el-card shadow="never" class="custom-page-header">
            <div class="header-flex">
                <div class="title-group">
                    <h1 class="main-title">工作量看板</h1>
                    <span class="sub-title">请选择时间段后点击查询以获取最新数据</span>
                </div>
                <div class="header-actions">
                    <el-button type="primary" :icon="Download" plain :loading="exportLoading" @click="handleExport"> 导出报表数据 </el-button>
                </div>
            </div>
        </el-card>

        <main class="dashboard-content">
            <WorkloadStats :data="rawData" />

            <WorkloadFilter 
                v-model:search="searchQuery" 
                v-model:type="filterType" 
                v-model:dateRange="dateRange" 
                :total="filteredData.length" 
                @query="fetchList" 
            />

            <el-card shadow="never" class="table-card">
                <WorkloadTable 
                    :data="paginatedData" 
                    :loading="loading" 
                    @view-detail="handleViewDetail" 
                    @export-row="handleExportRow"
                />

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

        <el-dialog
            v-model="detailVisible"
            :title="`业务处理明细汇总 - ${currentDetail?.operator || ''}`"
            width="1000px"
            destroy-on-close
            append-to-body
            class="custom-workload-dialog"
        >
            <WorkloadDetail 
                :data="currentDetail" 
                @close="detailVisible = false" 
                @export-current="handleExportCurrent"
            />
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="detailVisible = false">关闭窗口</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref } from "vue";
import { Download } from "@element-plus/icons-vue";
import { useWorkload } from "./composables/useWorkload";
import { ElMessage } from "element-plus";
import dayjs from "dayjs"; // 用于文件名日期

// 引入导出工具类
import ExportExcelUtity from "@/utity/exportExcel";

// 引入子组件
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

/**
 * 通用执行导出逻辑
 */
const performExport = (data, fileName) => {
    try {
        ExportExcelUtity(data, fileName, "Workload");
        ElMessage.success(`${fileName} 导出成功`);
    } catch (error) {
        ElMessage.error("导出失败：" + error.message);
    }
};

/**
 * 1. 顶部全量导出
 */
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

/**
 * 2. 表格行单人导出
 */
const handleExportRow = (row) => {
    const fileName = `工作量明细_${row.operator}_${dayjs().format('YYYYMMDD')}.xlsx`;
    performExport(row, fileName);
};

/**
 * 3. 详情页内单人导出
 */
const handleExportCurrent = (data) => {
    const fileName = `工作量明细_${data.operator}_${dayjs().format('YYYYMMDD')}.xlsx`;
    performExport(data, fileName);
};
</script>

<style scoped>
.workload-page { background-color: var(--el-bg-color-page); padding: 24px; min-height: 100vh; display: flex; flex-direction: column; gap: 20px; box-sizing: border-box; }
.custom-page-header { border-radius: 12px; border: none; background-color: var(--el-bg-color); }
.header-flex { display: flex; justify-content: space-between; align-items: center; }
.main-title { margin: 0; font-size: 24px; color: var(--el-text-color-primary); font-weight: 800; letter-spacing: -0.5px; }
.sub-title { font-size: 14px; color: var(--el-text-color-secondary); margin-top: 6px; display: block; }
.dashboard-content { display: flex; flex-direction: column; gap: 20px; flex: 1; }
.table-card { border-radius: 12px; border: none; background-color: var(--el-bg-color); }
.pagination-container { margin-top: 24px; display: flex; justify-content: flex-end; }
:deep(.custom-workload-dialog) { border-radius: 16px; overflow: hidden; }
:deep(.custom-workload-dialog .el-dialog__header) { margin-right: 0; padding: 20px 24px; border-bottom: 1px solid var(--el-border-color-lighter); background-color: var(--el-fill-color-blank); }
:deep(.custom-workload-dialog .el-dialog__title) { font-weight: 700; font-size: 18px; }
:deep(.custom-workload-dialog .el-dialog__body) { padding: 24px; }
</style>