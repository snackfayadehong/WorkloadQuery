<template>
    <div class="dict-compare-container">
        <header class="tool-hero-mini">
            <div class="hero-text">
                <h1 class="title">字典信息实时比对</h1>
                <p class="desc">支持通过材料代码或产品ID进行跨系统校对</p>
            </div>
        </header>

        <el-card shadow="never" class="glass-panel search-console">
            <div class="search-flex">
                <el-input v-model="keyword" placeholder="请输入材料代码或产品ID" class="modern-input" clearable
                    @keyup.enter="handleCompare">
                    <template #prefix>
                        <el-icon>
                            <Search />
                        </el-icon>
                    </template>
                    <template #suffix>
                        <span class="input-tip" v-if="keywordType">{{ keywordType }}</span>
                    </template>
                </el-input>
                <el-button type="primary" :loading="loading" class="search-btn" @click="handleCompare">
                    开始同步校验
                </el-button>
            </div>
        </el-card>

        <el-dialog v-model="choiceVisible" title="检测到多个本地条目匹配" width="750px" destroy-on-close append-to-body
            class="choice-dialog">
            <div class="choice-header-tip">该编码在本地库对应多个 产品ID，请选择具体项进行比对：</div>
            <el-table :data="multiOptions" border stripe @row-click="confirmSelection" style="cursor: pointer">
                <el-table-column prop="ProductInfoID" label="本地 ID" width="120" align="center" />
                <el-table-column prop="tymc" label="材料名称" min-width="180" />
                <el-table-column prop="ypgg" label="规格型号" min-width="180" />
                <el-table-column label="操作" width="100" align="center">
                    <template #default>
                        <el-button type="primary" link icon="Pointer">选择</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-dialog>

        <transition name="list-fade">
            <div v-if="compareResults.length" class="results-wrapper">
                <div class="result-header">
                    <div class="status-summary">
                        <el-badge :value="mismatchCount" :hidden="mismatchCount === 0" type="danger">
                            <span class="summary-text">比对完成：发现 {{ mismatchCount }} 项差异</span>
                        </el-badge>
                    </div>
                </div>

                <el-card shadow="never" class="glass-panel table-card">
                    <el-table :data="compareResults" style="width: 100%" :row-class-name="getRowClass">
                        <el-table-column label="字典属性项" prop="label" width="160">
                            <template #default="{ row }">
                                <span class="field-label">{{ row.label }}</span>
                            </template>
                        </el-table-column>

                        <el-table-column label="怡道系统 (Local DB)">
                            <template #default="{ row }">
                                <div class="val-box local">
                                    <span class="val">{{ row.localValue || '-' }}</span>
                                </div>
                            </template>
                        </el-table-column>

                        <el-table-column label="HIS系统 (Direct API)">
                            <template #default="{ row }">
                                <div class="val-box his">
                                    <span class="val" :class="{ 'text-danger': !row.isMatch }">
                                        {{ row.hisValue || '-' }}
                                    </span>
                                    <el-icon v-if="row.isMatch" class="icon-match">
                                        <CircleCheckFilled />
                                    </el-icon>
                                    <el-icon v-else class="icon-mismatch">
                                        <WarningFilled />
                                    </el-icon>
                                </div>
                            </template>
                        </el-table-column>

                        <el-table-column label="校验结果" width="120" align="center">
                            <template #default="{ row }">
                                <el-tag :type="row.isMatch ? 'success' : 'danger'" effect="dark">
                                    {{ row.isMatch ? '一致' : '冲突' }}
                                </el-tag>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
            </div>
        </transition>

        <el-empty v-if="!loading && !compareResults.length && hasSearched" :description="errorMsg || '未找到该材料的比对信息'">
            <template #extra>
                <p class="empty-sub-tip">请确认输入是否准确，或 HIS 系统是否已停用该材料代码</p>
            </template>
        </el-empty>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { Search, CircleCheckFilled, WarningFilled, Pointer } from "@element-plus/icons-vue";
import myAxios from '@/services/myAxios';
import { ElMessage } from 'element-plus';

const keyword = ref('');
const loading = ref(false);
const hasSearched = ref(false);
const compareResults = ref([]);
const errorMsg = ref('');

// 新增控制弹窗的变量
const choiceVisible = ref(false);
const multiOptions = ref([]);

// 智能识别输入类型
const keywordType = computed(() => {
    if (!keyword.value) return '';
    return /^\d{1,6}$/.test(keyword.value) ? 'ID 模式' : 'Code 模式';
});

const mismatchCount = computed(() =>
    compareResults.value.filter(item => !item.isMatch).length
);

const handleCompare = async () => {
    if (!keyword.value) {
        ElMessage.warning('请输入查询关键字');
        return;
    }

    loading.value = true;
    compareResults.value = [];
    errorMsg.value = '';
    hasSearched.value = true;

    try {
        const res = await myAxios.post('/dict/compare', { keyword: keyword.value });

        // 处理 201 多条记录冲突
        if (res.code === 201) {
            multiOptions.value = res.data;
            choiceVisible.value = true;
            loading.value = false;
            return;
        }

        if (res.code === 0) {
            if (res.data && res.data.length > 0) {
                compareResults.value = res.data;
            } else {
                errorMsg.value = res.message || 'HIS 系统未返回数据，该材料可能已停用';
            }
        } else {
            errorMsg.value = res.message || '查询失败';
        }
    } catch (err) {
        errorMsg.value = '连接 HIS 接口超时或本地数据库查询异常';
    } finally {
        loading.value = false;
    }
};

/**
 * 选中某条具体的本地记录 [逻辑修复]
 */
const confirmSelection = (row) => {
    keyword.value = row.ProductInfoID; // 使用具体的 ID 覆盖当前输入
    choiceVisible.value = false;
    handleCompare(); // 重新发起比对请求
};

const getRowClass = ({ row }) => {
    return row.isMatch ? '' : 'mismatch-row';
};
</script>

<style scoped>
.dict-compare-container {
    padding: 32px;
    background: var(--el-bg-color-page);
    min-height: 100%;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.tool-hero-mini {
    padding: 30px 40px;
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    border-radius: 20px;
    color: white;
    box-shadow: 0 10px 20px rgba(79, 172, 254, 0.2);
}

.title {
    font-size: 24px;
    font-weight: 800;
    margin: 0;
}

.desc {
    font-size: 14px;
    opacity: 0.8;
    margin-top: 8px;
}

.search-console {
    padding: 10px;
    border-radius: 16px;
}

.search-flex {
    display: flex;
    gap: 16px;
}

.input-tip {
    font-size: 12px;
    color: #409eff;
    padding: 0 8px;
    border-left: 1px solid #eee;
}

/* 弹窗样式 */
.choice-header-tip {
    margin-bottom: 16px;
    color: var(--el-text-color-secondary);
    font-size: 14px;
}

.results-wrapper {
    animation: slideUp 0.5s ease-out;
}

.result-header {
    margin-bottom: 16px;
    display: flex;
    justify-content: flex-end;
}

.summary-text {
    font-weight: 600;
    color: #606266;
}

.table-card {
    border-radius: 20px;
    overflow: hidden;
}

.field-label {
    font-weight: 700;
    color: #909399;
}

.val-box {
    padding: 8px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.val {
    font-family: 'Monaco', 'Courier New', monospace;
    font-weight: 700;
    font-size: 16px;
}

.icon-match {
    color: #67c23a;
    font-size: 20px;
}

.icon-mismatch {
    color: #f56c6c;
    font-size: 20px;
    animation: pulse 2s infinite;
}

:deep(.mismatch-row) {
    background-color: #fff1f0 !important;
}

:deep(.mismatch-row:hover) {
    background-color: #ffccc7 !important;
}

.text-danger {
    color: #f56c6c;
    text-decoration: underline wavy #ffccc7;
}

@keyframes pulse {
    0% {
        transform: scale(1);
        opacity: 1;
    }

    50% {
        transform: scale(1.2);
        opacity: 0.7;
    }

    100% {
        transform: scale(1);
        opacity: 1;
    }
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.empty-sub-tip {
    margin-top: 10px;
    color: #999;
    font-size: 13px;
}
</style>