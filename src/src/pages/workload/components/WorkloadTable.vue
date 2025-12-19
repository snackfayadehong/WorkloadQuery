<template>
  <div class="workload-table-wrapper">
    <el-skeleton :loading="loading" animated :rows="10">
      <template #template>
        <div class="skeleton-container">
          <el-skeleton-item variant="h3" style="width: 30%; margin-bottom: 20px" />
          <div v-for="i in 6" :key="i" class="skeleton-row">
            <el-skeleton-item variant="text" style="flex: 1" />
            <el-skeleton-item variant="text" style="flex: 2" />
            <el-skeleton-item variant="text" style="flex: 2" />
            <el-skeleton-item variant="text" style="flex: 1" />
          </div>
        </div>
      </template>

      <template #default>
        <el-table 
          :data="data" 
          border 
          stripe 
          style="width: 100%"
          header-cell-class-name="custom-table-header"
          row-key="operator"
        >
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="expand-content">
                <el-descriptions title="业务量明细汇总" :column="3" border size="small">
                  <el-descriptions-item label="入库项数">
                    <el-tag size="small" type="success">{{ row.inbound?.length || 0 }} 项</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="出库项数">
                    <el-tag size="small" type="danger">{{ row.outbound?.length || 0 }} 项</el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="退还项数">
                    <el-tag size="small" type="warning">{{ row.return?.length || 0 }} 项</el-tag>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="operator" label="操作人员" align="center" min-width="120">
            <template #default="{ row }">
              <span class="operator-name">{{ row.operator }}</span>
            </template>
          </el-table-column>

          <el-table-column label="入库总统计" align="center" min-width="150">
            <template #default="{ row }">
              <span class="amount-text success">{{ formatCurrency(calculateTotal(row.inbound)) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="出库总统计" align="center" min-width="150">
            <template #default="{ row }">
              <span class="amount-text danger">{{ formatCurrency(calculateTotal(row.outbound)) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="退还统计" align="center" min-width="150">
            <template #default="{ row }">
              <span class="amount-text warning">{{ formatCurrency(calculateTotal(row.return)) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" align="center" width="160" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link icon="View" @click="$emit('view-detail', row)">详情</el-button>
              <el-button type="success" link icon="Download" @click="$emit('export-row', row)">导出</el-button>
            </template>
          </el-table-column>
          
          <template #empty>
            <el-empty description="暂无工作量数据" :image-size="100" />
          </template>
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

// 增加 export-row 事件定义
defineEmits(['view-detail', 'export-row']);

const formatCurrency = (val) => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
  }).format(val);
};

const calculateTotal = (items) => {
  if (!items || !Array.isArray(items)) return 0;
  return items.reduce((sum, item) => sum + (item.totalAmount || 0), 0);
};
</script>

<style scoped>
.workload-table-wrapper { background: transparent; width: 100%; }
.skeleton-container { padding: 20px; }
.skeleton-row { display: flex; gap: 20px; margin-bottom: 15px; }
.expand-content { padding: 20px 50px; background: var(--el-fill-color-lighter); }
.operator-name { font-weight: 600; font-size: 18px; color: var(--el-text-color-primary); }
.amount-text { font-family: 'Monaco', monospace; font-weight: bold; font-size: 18px; }
.success { color: var(--el-color-success); }
.danger { color: var(--el-color-danger); }
.warning { color: var(--el-color-warning); }
:deep(.custom-table-header) { background-color: var(--el-fill-color-light) !important; color: var(--el-text-color-primary); font-weight: 700; }
:global(html.dark) .expand-content { background: #1d1d1d; }
:global(html.dark) .el-table { --el-table-border-color: #333; --el-table-header-bg-color: #1a1a1a; }
</style>