<template>
    <el-dialog v-model="visible" :title="`操作详情 - ${data.operator || '未知'}`" width="80%" top="5vh" destroy-on-close class="custom-detail-dialog" @close="$emit('update:modelValue', false)">
        <div class="detail-summary">
            <el-row :gutter="20">
                <el-col :span="8">
                    <div class="summary-item inbound">
                        <span class="label">入库总计</span>
                        <span class="value">{{ formatCurrency(totals.inbound) }}</span>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="summary-item outbound">
                        <span class="label">出库总计</span>
                        <span class="value">{{ formatCurrency(totals.outbound) }}</span>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="summary-item return">
                        <span class="label">退还总计</span>
                        <span class="value">{{ formatCurrency(totals.return) }}</span>
                    </div>
                </el-col>
            </el-row>
        </div>

        <el-tabs v-model="activeTab" class="detail-tabs">
            <el-tab-pane label="入库明细" name="inbound">
                <el-table :data="data.inbound" border height="400px" stripe size="small">
                    <el-table-column prop="billNo" label="单据编号" width="180" />
                    <el-table-column prop="itemName" label="物料名称" min-width="150" />
                    <el-table-column prop="quantity" label="数量" width="100" align="right" />
                    <el-table-column prop="totalAmount" label="金额" width="120" align="right">
                        <template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template>
                    </el-table-column>
                    <el-table-column prop="createTime" label="操作时间" width="160" />
                </el-table>
            </el-tab-pane>

            <el-tab-pane label="出库明细" name="outbound">
                <el-table :data="data.outbound" border height="400px" stripe size="small">
                    <el-table-column prop="billNo" label="单据编号" width="180" />
                    <el-table-column prop="itemName" label="物料名称" min-width="150" />
                    <el-table-column prop="totalAmount" label="金额" width="120" align="right">
                        <template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template>
                    </el-table-column>
                    <el-table-column prop="createTime" label="操作时间" width="160" />
                </el-table>
            </el-tab-pane>
        </el-tabs>

        <template #footer>
            <div class="dialog-footer">
                <el-button @click="visible = false">关闭</el-button>
                <el-button type="primary" icon="Download">导出该员报表</el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from "vue";

const props = defineProps({
    modelValue: Boolean,
    data: {
        type: Object,
        default: () => ({ operator: "", inbound: [], outbound: [], return: [] })
    }
});

const emit = defineEmits(["update:modelValue"]);

const visible = ref(props.modelValue);
const activeTab = ref("inbound");

// 同步弹窗显示状态
watch(
    () => props.modelValue,
    val => (visible.value = val)
);
watch(visible, val => emit("update:modelValue", val));

// 格式化金额
const formatCurrency = val => {
    return new Intl.NumberFormat("zh-CN", { style: "currency", currency: "CNY" }).format(val || 0);
};

// 计算各类型总额
const totals = computed(() => {
    const calc = arr => (arr || []).reduce((s, i) => s + (i.totalAmount || 0), 0);
    return {
        inbound: calc(props.data.inbound),
        outbound: calc(props.data.outbound),
        return: calc(props.data.return)
    };
});
</script>

<style scoped>
.detail-summary {
    margin-bottom: 24px;
}
.summary-item {
    padding: 16px;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.summary-item .label {
    font-size: 13px;
    opacity: 0.8;
}
.summary-item .value {
    font-size: 20px;
    font-weight: bold;
}

/* 针对明暗模式的背景适配 */
.inbound {
    background: var(--el-color-success-light-9);
    color: var(--el-color-success);
}
.outbound {
    background: var(--el-color-danger-light-9);
    color: var(--el-color-danger);
}
.return {
    background: var(--el-color-warning-light-9);
    color: var(--el-color-warning);
}

:global(html.dark) .summary-item {
    background: rgba(255, 255, 255, 0.05);
}

.detail-tabs {
    margin-top: 10px;
}
</style>
