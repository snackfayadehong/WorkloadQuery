<template>
  <div class="workload-table-wrapper">
    <el-skeleton :loading="loading" animated :rows="10">
      <template #default>
        <el-table :data="data" border stripe style="width: 100%" class="modern-dashboard-table"
          header-cell-class-name="modern-table-header" row-key="operator">
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="expand-content">
                <el-row :gutter="40">
                  <el-col :span="14">
                    <h4 class="preview-title">核心业务预览 (金额最高项)</h4>
                    <div class="top-items-grid">
                      <div class="top-item success-bg">
                        <span class="type-label">入库主项：</span>
                        <span class="item-val">{{ getTopCategory(row.inbound) }}</span>
                      </div>
                      <div class="top-item danger-bg">
                        <span class="type-label">出库主项：</span>
                        <span class="item-val">{{ getTopCategory(row.outbound) }}</span>
                      </div>
                    </div>
                  </el-col>
                  <el-col :span="10">
                    <h4 class="preview-title">单据分类统计</h4>
                    <el-descriptions :column="1" border size="small">
                      <el-descriptions-item label="入库总计">{{ row.inbound?.length || 0 }} 类物料</el-descriptions-item>
                      <el-descriptions-item label="出库总计">{{ row.outbound?.length || 0 }} 类物料</el-descriptions-item>
                      <el-descriptions-item label="退还总计">{{ row.return?.length || 0 }} 类物料</el-descriptions-item>
                    </el-descriptions>
                  </el-col>
                </el-row>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="operator" label="操作人员" align="center" min-width="120">
            <template #default="{ row }">
              <span class="operator-name-modern">{{ row.operator }}</span>
            </template>
          </el-table-column>

          <el-table-column label="入库汇总" align="center" min-width="160">
            <template #default="{ row }">
              <span class="amount-val success">{{ formatCurrency(calculateTotal(row.inbound)) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="出库汇总" align="center" min-width="160">
            <template #default="{ row }">
              <span class="amount-val danger">{{ formatCurrency(calculateTotal(row.outbound)) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="退还统计" align="center" min-width="160">
            <template #default="{ row }">
              <span class="amount-val warning">{{ formatCurrency(calculateTotal(row.return)) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="管理" align="center" width="160" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link icon="View" @click="$emit('view-detail', row)">明细</el-button>
              <el-button type="success" link icon="Download" @click="$emit('export-row', row)">导出</el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup>
const props = defineProps({
  data: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false }
});

defineEmits(['view-detail', 'export-row']);

/**
 * 核心逻辑：计算金额最大的分类用于预览
 */
const getTopCategory = (items) => {
  // 增加对非数组或空数组的拦截
  if (!items || !Array.isArray(items) || items.length === 0) return "暂无业务记录";
  // 使用解构赋值确保不污染原数据，并增加可选链
  const sorted = [...items].sort((a, b) => (b.totalAmount || 0) - (a.totalAmount || 0));
  const top = sorted[0];
  return `${top.category} (￥${(top.totalAmount || 0).toLocaleString()})`;
  // if (!items || items.length === 0) return "暂无业务记录";
  // const top = [...items].sort((a, b) => (b.totalAmount || 0) - (a.totalAmount || 0))[0];
  // return `${top.category} (￥${top.totalAmount.toLocaleString()})`;
};

const formatCurrency = (val) => {
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY' }).format(val);
};

const calculateTotal = (items) => {
  if (!items || !Array.isArray(items)) return 0;
  return items.reduce((sum, item) => sum + (item.totalAmount || 0), 0);
};
</script>

<style scoped>
.workload-table-wrapper {
  background: transparent;
  width: 100%;
}

/* 表格整体风格：轻量、圆角 */
.modern-dashboard-table {
  --el-table-border-color: #f0f0f0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
}

.operator-name-modern {
  font-weight: 700;
  color: #262626;
  font-size: 16px;
}

.amount-val {
  font-family: 'Monaco', monospace;
  font-weight: 800;
  font-size: 17px;
}

.success {
  color: #52c41a;
}

.danger {
  color: #ff4d4f;
}

.warning {
  color: #faad14;
}

/* 展开行预览样式 */
.expand-content {
  padding: 24px 40px;
  background: rgba(64, 158, 255, 0.02);
}

.preview-title {
  margin: 0 0 16px;
  font-size: 14px;
  color: #8c8c8c;
  font-weight: 500;
}

.top-items-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.top-item {
  padding: 10px 16px;
  border-radius: 8px;
  font-size: 14px;
  border: 1px solid transparent;
}

.success-bg {
  background: #f6ffed;
  border-color: #b7eb8f;
  color: #389e0d;
}

.danger-bg {
  background: #fff1f0;
  border-color: #ffa39e;
  color: #cf1322;
}

.item-val {
  font-weight: bold;
  margin-left: 4px;
}

/* 表头美化 */
:deep(.modern-table-header) {
  background-color: #fafafa !important;
  color: #262626;
  font-weight: 700;
  height: 50px;
}

:global(html.dark) .modern-dashboard-table {
  --el-table-border-color: #333;
}

:global(html.dark) .operator-name-modern {
  color: #e5eaf3;
}

:global(html.dark) .expand-content {
  background: #1a1a1a;
}
</style>