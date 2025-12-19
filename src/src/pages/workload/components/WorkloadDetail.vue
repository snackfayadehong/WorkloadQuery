<template>
  <div class="workload-detail-container">
    <el-card shadow="never" class="detail-card">
      <template #header>
        <div class="detail-header">
          <div class="user-info">
            <el-icon class="user-icon"><UserFilled /></el-icon>
            <span class="operator-name">{{ data?.operator || '未知人员' }}</span>
            <el-tag size="small" effect="plain">业务汇总明细</el-tag>
          </div>
          <div class="header-actions">
            <el-button type="success" size="small" icon="Download" @click="$emit('export-current', data)">导出此人明细</el-button>
            <el-button type="info" link @click="$emit('close')">收起详情</el-button>
          </div>
        </div>
      </template>

      <div class="sections-stack">
        <div v-if="data?.inbound?.length" class="detail-section inbound">
          <h4 class="section-title">
            <el-icon><Download /></el-icon> 入库汇总 (Inbound)
            <el-tag type="success" size="small" effect="light" class="item-count">{{ data.inbound.length }} 分类</el-tag>
          </h4>
          <el-table :data="data.inbound" border stripe size="default" class="custom-detail-table">
            <el-table-column prop="category" label="物料分类名称" min-width="180" />
            <el-table-column prop="specCount" label="品规总数" align="center" width="100" />
            <el-table-column prop="billCount" label="单据数量" align="center" width="100" />
            <el-table-column prop="totalAmount" label="合计金额" align="right" width="150">
              <template #default="{ row }">
                <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, {minimumFractionDigits: 2}) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="data?.outbound?.length" class="detail-section outbound">
          <h4 class="section-title">
            <el-icon><Upload /></el-icon> 出库汇总 (Outbound)
            <el-tag type="danger" size="small" effect="light" class="item-count">{{ data.outbound.length }} 分类</el-tag>
          </h4>
          <el-table :data="data.outbound" border stripe size="default" class="custom-detail-table">
            <el-table-column prop="category" label="物料分类名称" min-width="180" />
            <el-table-column prop="specCount" label="品规总数" align="center" width="100" />
            <el-table-column prop="billCount" label="单据数量" align="center" width="100" />
            <el-table-column prop="totalAmount" label="合计金额" align="right" width="150">
              <template #default="{ row }">
                <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, {minimumFractionDigits: 2}) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="data?.return?.length" class="detail-section return">
          <h4 class="section-title">
            <el-icon><RefreshLeft /></el-icon> 退还汇总 (Return)
            <el-tag type="warning" size="small" effect="light" class="item-count">{{ data.return.length }} 分类</el-tag>
          </h4>
          <el-table :data="data.return" border stripe size="default" class="custom-detail-table">
            <el-table-column prop="category" label="物料分类名称" min-width="180" />
            <el-table-column prop="specCount" label="品规总数" align="center" width="100" />
            <el-table-column prop="billCount" label="单据数量" align="center" width="100" />
            <el-table-column prop="totalAmount" label="合计金额" align="right" width="150">
              <template #default="{ row }">
                <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, {minimumFractionDigits: 2}) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <el-empty v-if="!hasAnyData" :image-size="120" description="选定统计时间内该操作员无相关业务流水记录" />
    </el-card>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { UserFilled, Download, Upload, RefreshLeft } from "@element-plus/icons-vue";

const props = defineProps({
  data: { type: Object, default: () => ({}) }
});

// 增加 export-current 事件定义
defineEmits(['close', 'export-current']);

const hasAnyData = computed(() => {
  return (props.data?.inbound?.length || 0) > 0 || 
         (props.data?.outbound?.length || 0) > 0 || 
         (props.data?.return?.length || 0) > 0;
});
</script>

<style scoped>
.workload-detail-container { margin-top: 0; animation: slideInUp 0.3s ease-out; }
@keyframes slideInUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
.detail-card { border-top: 2px solid #409eff; border-radius: 8px; }
.detail-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; align-items: center; gap: 12px; }
.user-info { display: flex; align-items: center; gap: 10px; }
.user-icon { font-size: 20px; color: #409eff; }
.operator-name { font-weight: bold; font-size: 16px; }
.sections-stack { display: flex; flex-direction: column; gap: 32px; }
.detail-section { padding: 10px; border-radius: 4px; }
.section-title { margin: 0 0 16px 0; display: flex; align-items: center; gap: 10px; font-size: 17px; font-weight: 700; color: var(--el-text-color-primary); padding-bottom: 8px; border-bottom: 1px dashed #ebeef5; }
.item-count { margin-left: 4px; font-weight: normal; }
.inbound .section-title { color: var(--el-color-success); }
.outbound .section-title { color: var(--el-color-danger); }
.return .section-title { color: var(--el-color-warning); }
.custom-detail-table { font-size: 14px; }
.amount-font { font-family: 'Monaco', 'Courier New', monospace; font-weight: 600; color: var(--el-text-color-primary); }
:deep(.el-table__header th) { background-color: var(--el-fill-color-light) !important; font-weight: 700; color: var(--el-text-color-primary); }
</style>