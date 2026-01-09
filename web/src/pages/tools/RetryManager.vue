<template>
    <div class="retry-manager-container">
        <section class="unified-hero-card">
            <div class="hero-content">
                <h1 class="hero-title">接口重试管理中心</h1>
                <p class="hero-desc">监控并补偿发送失败的出库单与退库单据</p>

                <div class="search-area glass-effect">
                    <el-select v-model="queryType" placeholder="单据类型" class="type-select">
                        <el-option label="领用/其他出库单" value="delivery" />
                        <el-option label="科室退库单" value="refund" />
                    </el-select>
                    <el-date-picker
                        v-model="dateRange"
                        type="datetimerange"
                        range-separator="至"
                        start-placeholder="开始时间"
                        end-placeholder="结束时间"
                        value-format="YYYY-MM-DD HH:mm:ss"
                        class="integrated-picker"
                    />
                    <el-button type="primary" :loading="loading" class="audit-btn" @click="handleQuery"> 查询待重试单据 </el-button>
                </div>
            </div>
        </section>

        <transition name="list-fade">
            <div v-if="tableData.length" class="results-container">
                <el-card shadow="never" class="table-card glass-panel">
                    <el-table :data="tableData" style="width: 100%" border stripe>
                        <el-table-column type="expand">
                            <template #default="props">
                                <div class="detail-wrapper">
                                    <el-descriptions title="单据详细信息" :column="2" border>
                                        <template v-if="queryType === 'delivery'">
                                            <el-descriptions-item label="明细序号">{{ props.row.detailSort }}</el-descriptions-item>
                                            <el-descriptions-item label="出库方式">{{ props.row.ckfs }}</el-descriptions-item>
                                            <el-descriptions-item label="供货库房">{{ props.row.storeHouseName }}</el-descriptions-item>
                                            <el-descriptions-item label="领用二级库房">{{ props.row.leaderDepartName }}</el-descriptions-item>
                                        </template>
                                        <template v-else>
                                            <el-descriptions-item label="单据类型">科室退库单</el-descriptions-item>
                                            <el-descriptions-item label="入库方式">{{ props.row.rkfs }}</el-descriptions-item>
                                            <el-descriptions-item label="操作提示">请确认 HIS 系统对应单据状态后进行重试</el-descriptions-item>
                                        </template>
                                    </el-descriptions>
                                </div>
                            </template>
                        </el-table-column>

                        <el-table-column label="单据编号" width="220">
                            <template #default="scope">
                                <span>{{ queryType === "delivery" ? scope.row.deliveryCode : scope.row.yddh }}</span>
                            </template>
                        </el-table-column>

                        <template v-if="queryType === 'delivery'">
                            <el-table-column prop="storeHouseName" label="供货库房" show-overflow-tooltip />
                            <el-table-column prop="leaderDepartName" label="领用二级库房" show-overflow-tooltip />
                            <el-table-column prop="ckfs" label="出库类型" width="120" />
                        </template>

                        <template v-else>
                            <el-table-column label="单据类型" width="150">
                                <template #default>科室退库单</template>
                            </el-table-column>
                            <el-table-column prop="rkfs" label="入库方式" show-overflow-tooltip />
                            <el-table-column label="" />
                        </template>

                        <el-table-column label="操作" width="150" align="center" fixed="right">
                            <template #default="scope">
                                <el-button type="primary" size="small" icon="Refresh" @click="executeRetry(scope.row)"> 触发重试 </el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
            </div>
        </transition>

        <el-empty v-if="!loading && hasSearched && tableData.length === 0" description="暂无符合条件的待重试单据" />
    </div>
</template>

<script setup>
import { ref, watch } from "vue"; // 引入 watch
import { Refresh } from "@element-plus/icons-vue";
import myAxios from "@/services/myAxios";
import { ElMessage, ElMessageBox } from "element-plus";

// --- 手动引入样式，修复提示框左上角显示问题 ---
import "element-plus/es/components/message-box/style/css";
import "element-plus/es/components/message/style/css";

const queryType = ref("delivery");
const dateRange = ref([]);
const loading = ref(false);
const hasSearched = ref(false);
const tableData = ref([]);

// --- 修复问题 1：切换类型时清空数据 ---
watch(queryType, () => {
    tableData.value = [];
    hasSearched.value = false;
});

// 查询逻辑
const handleQuery = async () => {
    if (!dateRange.value || dateRange.value.length === 0) {
        return ElMessage.warning("请选择查询时间范围");
    }

    loading.value = true;
    hasSearched.value = true;

    try {
        const res = await myAxios.post(`retry/list`, {
            queryType: queryType.value,
            startTime: dateRange.value[0],
            endTime: dateRange.value[1]
        });
        tableData.value = res.data || [];
    } catch (err) {
        ElMessage.error("查询失败");
    } finally {
        loading.value = false;
    }
};

// 触发重试逻辑
const executeRetry = row => {
    const displayId = queryType.value === "delivery" ? row.deliveryCode : row.yddh;

    ElMessageBox.confirm(`确认重新推送单据 ${displayId} 到 HIS 系统吗？`, "操作确认", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        center: true, // 居中显示
        draggable: true // 允许拖拽
    }).then(async () => {
        try {
            const res = await myAxios.post(`retry/execute`, {
                type: queryType.value,
                billNo: row.yddh,
                detailSort: row.detailSort || ""
            });

            if (res.code === 200) {
                ElMessage.success("重试指令已下发");
                handleQuery();
            }
        } catch (err) {
            ElMessage.error("触发重试异常");
        }
    });
};
</script>

<style scoped>
/* 样式部分保持不变 */
.retry-manager-container {
    padding: 24px;
    background-color: #f5f7fa;
    min-height: 100vh;
}

.unified-hero-card {
    padding: 40px;
    background: linear-gradient(135deg, #67c23a 0%, #4facfe 100%);
    border-radius: 24px;
    color: white;
    box-shadow: 0 12px 24px rgba(103, 194, 58, 0.2);
    margin-bottom: 24px;
}

.hero-content {
    max-width: 1000px;
    margin: 0 auto;
    text-align: center;
}

.hero-title {
    font-size: 28px;
    font-weight: 800;
}

.hero-desc {
    font-size: 15px;
    opacity: 0.9;
    margin: 10px 0 30px;
}

.glass-effect {
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    padding: 12px;
    border-radius: 16px;
    display: flex;
    gap: 12px;
    border: 1px solid rgba(255, 255, 255, 0.3);
}

.type-select {
    width: 170px;
}
.integrated-picker {
    flex: 1;
}
.table-card {
    border-radius: 16px;
    border: none;
}
.detail-wrapper {
    padding: 20px 40px;
    background-color: #fcfdfe;
}
:deep(.el-table__expanded-cell) {
    padding: 0 !important;
}
</style>
