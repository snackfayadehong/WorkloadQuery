<template>
  <div class="workload-detail-wrapper">
    <div class="detail-summary-row">
      <div class="summary-item success">
        <div class="s-icon"><el-icon>
            <Download />
          </el-icon></div>
        <div class="s-info">
          <span class="s-label">入库总计</span>
          <span class="s-value">￥{{ calculateSectionTotal(data?.inbound).toLocaleString() }}</span>
        </div>
      </div>
      <div class="summary-item danger">
        <div class="s-icon"><el-icon>
            <Upload />
          </el-icon></div>
        <div class="s-info">
          <span class="s-label">出库总计</span>
          <span class="s-value">￥{{ calculateSectionTotal(data?.outbound).toLocaleString() }}</span>
        </div>
      </div>
      <div class="summary-item warning">
        <div class="s-icon"><el-icon>
            <RefreshLeft />
          </el-icon></div>
        <div class="s-info">
          <span class="s-label">退还总计</span>
          <span class="s-value">￥{{ calculateSectionTotal(data?.return).toLocaleString() }}</span>
        </div>
      </div>
    </div>

    <el-card shadow="never" class="modern-detail-card">
      <template #header>
        <div class="detail-header">
          <div class="user-info">
            <div class="user-avatar-mini">
              <el-icon>
                <UserFilled />
              </el-icon>
            </div>
            <span class="operator-name">{{ data?.operator || '未知人员' }}</span>
            <el-tag size="small" round effect="plain" class="status-tag">数据汇总就绪</el-tag>
          </div>
          <div class="header-actions">
            <el-button type="success" size="small" icon="Download" plain @click="$emit('export-current', data)">导出
              Excel</el-button>
            <el-button type="info" link @click="$emit('close')">收起详情</el-button>
          </div>
        </div>
      </template>

      <div class="sections-stack">
        <div v-if="data?.inbound?.length" class="detail-section">
          <h4 class="modern-section-title in">
            <span class="indicator"></span> 入库汇总 (Inbound)
            <span class="count-badge">{{ data.inbound.length }} 类</span>
          </h4>
          <el-table :data="data.inbound" border stripe class="glass-table">
            <el-table-column prop="category" label="物料分类名称" min-width="180" />
            <el-table-column prop="specCount" label="品规数" align="center" width="100" />
            <el-table-column prop="billCount" label="单据数" align="center" width="100" />
            <el-table-column prop="totalAmount" label="合计金额" align="right" width="150">
              <template #default="{ row }">
                <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, { minimumFractionDigits: 2 })
                  }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="data?.outbound?.length" class="detail-section">
          <h4 class="modern-section-title out">
            <span class="indicator"></span> 出库汇总 (Outbound)
            <span class="count-badge">{{ data.outbound.length }} 类</span>
          </h4>
          <el-table :data="data.outbound" border stripe class="glass-table">
            <el-table-column prop="category" label="物料分类名称" min-width="180" />
            <el-table-column prop="specCount" label="品规数" align="center" width="100" />
            <el-table-column prop="billCount" label="单据数" align="center" width="100" />
            <el-table-column prop="totalAmount" label="合计金额" align="right" width="150">
              <template #default="{ row }">
                <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, { minimumFractionDigits: 2 })
                  }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="data?.return?.length" class="detail-section">
          <h4 class="modern-section-title ret">
            <span class="indicator"></span> 退还汇总 (Return)
            <span class="count-badge">{{ data.return.length }} 类</span>
          </h4>
          <el-table :data="data.return" border stripe class="glass-table">
            <el-table-column prop="category" label="物料分类名称" min-width="180" />
            <el-table-column prop="specCount" label="品规数" align="center" width="100" />
            <el-table-column prop="billCount" label="单据数" align="center" width="100" />
            <el-table-column prop="totalAmount" label="合计金额" align="right" width="150">
              <template #default="{ row }">
                <span class="amount-font">￥{{ row.totalAmount?.toLocaleString(undefined, { minimumFractionDigits: 2 })
                  }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <el-empty v-if="!hasAnyData" :image-size="120" description="暂无相关业务流水记录" />
    </el-card>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { UserFilled, Download, Upload, RefreshLeft } from "@element-plus/icons-vue";

const props = defineProps({
  data: { type: Object, default: () => ({}) }
});

defineEmits(['close', 'export-current']);

const hasAnyData = computed(() => {
  return (props.data?.inbound?.length || 0) > 0 ||
    (props.data?.outbound?.length || 0) > 0 ||
    (props.data?.return?.length || 0) > 0;
});

const calculateSectionTotal = (list) => {
  if (!list) return 0;
  return list.reduce((s, i) => s + (i.totalAmount || 0), 0);
};
</script>

<style scoped>
.workload-detail-wrapper {
  padding-bottom: 20px;
}

/* 顶部迷你汇总卡片 */
.detail-summary-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.summary-item {
  padding: 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s;
}

.summary-item.success {
  background: #f6ffed;
  border: 1px solid #b7eb8f;
  color: #52c41a;
}

.summary-item.danger {
  background: #fff1f0;
  border: 1px solid #ffa39e;
  color: #f5222d;
}

.summary-item.warning {
  background: #fffbe6;
  border: 1px solid #ffe58f;
  color: #faad14;
}

.s-icon {
  font-size: 24px;
  opacity: 0.8;
}

.s-label {
  font-size: 12px;
  display: block;
  opacity: 0.7;
  margin-bottom: 2px;
}

.s-value {
  font-size: 18px;
  font-weight: 800;
  font-family: 'Monaco', monospace;
}

/* 卡片样式优化 */
.modern-detail-card {
  border-radius: 16px !important;
  border: 1px solid #ebeef5 !important;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar-mini {
  width: 32px;
  height: 32px;
  background: #f0f7ff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
}

.operator-name {
  font-weight: 700;
  font-size: 18px;
  color: #303133;
}

/* 业务板块标题 */
.modern-section-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 16px;
}

.indicator {
  width: 4px;
  height: 16px;
  border-radius: 2px;
  margin-right: 10px;
}

.in .indicator {
  background: #67C23A;
}

.out .indicator {
  background: #F56C6C;
}

.ret .indicator {
  background: #E6A23C;
}

.count-badge {
  margin-left: 10px;
  font-size: 12px;
  font-weight: normal;
  color: #909399;
  background: #f4f4f5;
  padding: 2px 8px;
  border-radius: 10px;
}

.amount-font {
  font-family: 'Monaco', 'Courier New', monospace;
  font-weight: 700;
  color: #303133;
}

:deep(.el-table) {
  --el-table-header-bg-color: #fafafa;
  border-radius: 8px;
  overflow: hidden;
}
</style>