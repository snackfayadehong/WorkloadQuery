<template>
    <div class="dict-compare-container">
        <header class="tool-hero">
            <div class="hero-left">
                <h1 class="title">字典数据一致性审查</h1>
                <p class="desc">跨系统字典比对工具 • 实时同步校验</p>
            </div>
            <div class="hero-right">
                <div class="search-box glass-panel">
                    <el-input v-model="keyword" placeholder="输入 材料代码 或 产品ID" class="search-input"
                        @keyup.enter="handleCompare">
                        <template #prefix><el-icon>
                                <Search />
                            </el-icon></template>
                        <template #suffix><span class="mode-tag">{{ keywordType }}</span></template>
                    </el-input>
                    <el-button type="primary" :loading="loading" @click="handleCompare">立即比对</el-button>
                </div>
            </div>
        </header>

        <el-dialog v-model="choiceVisible" title="选择匹配条目" width="800px" append-to-body destroy-on-close
            class="custom-dialog">
            <el-table :data="multiOptions" @row-click="confirmSelection" highlight-current-row style="cursor: pointer">
                <el-table-column prop="ProductInfoID" label="产品 ID" width="100" />
                <el-table-column prop="ypdm" label="编码" width="120" />
                <el-table-column prop="tymc" label="材料名称" min-width="150" show-overflow-tooltip />
                <el-table-column prop="ypgg" label="规格型号" min-width="150" />
            </el-table>
        </el-dialog>

        <main v-if="compareResults.length" class="compare-workspace">
            <div class="audit-summary glass-panel" :class="{ 'has-diff': mismatchCount > 0 }">
                <div class="summary-item">
                    <span class="s-label">比对项总数</span>
                    <span class="s-value">{{ compareResults.length }}</span>
                </div>
                <div class="summary-divider"></div>
                <div class="summary-item">
                    <span class="s-label">异常差异项</span>
                    <span class="s-value danger">{{ mismatchCount }}</span>
                </div>
                <div class="summary-status">
                    <el-tag :type="mismatchCount === 0 ? 'success' : 'danger'" effect="dark" round>
                        {{ mismatchCount === 0 ? '数据完全一致' : '检测到字段冲突' }}
                    </el-tag>
                </div>
            </div>

            <div class="diff-list">
                <div class="diff-header-row">
                    <div class="col-label">属性维度</div>
                    <div class="col-local">怡道系统记录</div>
                    <div class="col-indicator">状态</div>
                    <div class="col-his">HIS系统 (Result)</div>
                </div>

                <div v-for="(item, index) in compareResults" :key="item.field" class="diff-card-row"
                    :class="{ 'is-mismatch': !item.isMatch }" :style="{ animationDelay: index * 0.05 + 's' }">
                    <div class="col-label">
                        <span class="field-name">{{ item.label }}</span>
                        <span class="field-key">{{ item.field }}</span>
                    </div>

                    <div class="col-local">
                        <div class="val-display">{{ item.localValue || '未维护' }}</div>
                    </div>

                    <div class="col-indicator">
                        <el-icon v-if="item.isMatch" class="match-icon">
                            <SuccessFilled />
                        </el-icon>
                        <el-icon v-else class="mismatch-icon">
                            <WarningFilled />
                        </el-icon>
                    </div>

                    <div class="col-his">
                        <div class="val-display" :title="!item.isMatch ? '建议修正为 HIS 值' : ''">
                            {{ item.hisValue || '未返回' }}
                        </div>
                    </div>
                </div>
            </div>
        </main>

        <div v-else-if="!loading && hasSearched" class="empty-state">
            <el-empty :description="errorMsg || '未检索到比对数据'">
                <template #extra>
                    <p class="empty-hint">请核实编码是否存在或 HIS 是否已将其停用</p>
                </template>
            </el-empty>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { Search, SuccessFilled, WarningFilled } from "@element-plus/icons-vue";
import myAxios from '@/services/myAxios';
import { ElMessage } from 'element-plus';

const keyword = ref('');
const loading = ref(false);
const hasSearched = ref(false);
const compareResults = ref([]);
const errorMsg = ref('');
const choiceVisible = ref(false);
const multiOptions = ref([]);

const keywordType = computed(() => {
    if (!keyword.value) return '待输入';
    return /^\d{1,6}$/.test(keyword.value) ? 'ID' : 'Code';
});

const mismatchCount = computed(() =>
    compareResults.value.filter(item => !item.isMatch).length
);

let msgInstance = null; // 消息实体
let warnCount = 0; // 警告计数
let resetTimer = null; // 计时器
const handleCompare = async () => {
    if (!keyword.value) {
        warnCount = Math.min(warnCount + 1, 99);

        const tipContent = h(
            'div',
            {
                style: 'display: flex; align-items: center; gap: 8px;'
            },
            [
                h(
                    'span',
                    null,
                    warnCount > 1 ? '请勿重复点击！' : '请输入要比对的编码或 ID'
                ),
                warnCount > 1
                    ? h(
                        'span',
                        {
                            style: `
                              background: #f56c6c;
                              color: #fff;
                              padding: 0 6px;
                              border-radius: 10px;
                              font-size: 10px;
                              height: 16px;
                              line-height: 16px;
                          `
                        },
                        warnCount
                    )
                    : null
            ]
        );

        // 如果已有 message，先关闭
        if (msgInstance) {
            msgInstance.close();
        }

        msgInstance = ElMessage({
            message: tipContent,
            type: 'warning',
            duration: 2000
        });
        // 3秒后如果没有再次点击，自动清零计数，下次点击重新从“请输入”开始
        clearTimeout(resetTimer);
        resetTimer = setTimeout(() => {
            warnCount = 0;
        }, 3000);

        return;
    }
    loading.value = true;
    compareResults.value = [];
    hasSearched.value = true;

    try {
        const res = await myAxios.post('/dict/compare', { keyword: keyword.value });
        if (res.code === 201) {
            multiOptions.value = res.data;
            choiceVisible.value = true;
        } else if (res.code === 0) {
            compareResults.value = res.data;
        } else {
            errorMsg.value = res.message;
        }
    } catch (err) {
        errorMsg.value = '通信异常';
    } finally {
        loading.value = false;
    }
};

const confirmSelection = (row) => {
    keyword.value = row.ProductInfoID;
    choiceVisible.value = false;
    handleCompare();
};
</script>

<style scoped>
.dict-compare-container {
    padding: 30px;
    background-color: #f8fafc;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

/* Hero 布局优化 */
.tool-hero {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 40px;
    background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
    border-radius: 24px;
    color: white;
}

.title {
    font-size: 28px;
    font-weight: 800;
    margin: 0;
    letter-spacing: 1px;
}

.desc {
    margin: 8px 0 0;
    opacity: 0.6;
    font-size: 14px;
}

.search-box {
    display: flex;
    gap: 12px;
    padding: 8px;
    background: rgba(255, 255, 255, 0.1) !important;
    border-radius: 16px;
}

.mode-tag {
    font-size: 12px;
    color: #38bdf8;
    font-weight: bold;
    padding: 0 8px;
}

/* 汇总状态栏 */
.audit-summary {
    display: flex;
    align-items: center;
    padding: 20px 30px;
    margin-bottom: 20px;
    border-left: 5px solid #64748b;
}

.audit-summary.has-diff {
    border-left-color: #ef4444;
    background: #fef2f2 !important;
}

.summary-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.s-label {
    font-size: 12px;
    color: #94a3b8;
}

.s-value {
    font-size: 20px;
    font-weight: 800;
}

.s-value.danger {
    color: #ef4444;
}

.summary-divider {
    width: 1px;
    height: 30px;
    background: #e2e8f0;
    margin: 0 30px;
}

.summary-status {
    margin-left: auto;
}

/* 镜像对比列表 */
.diff-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.diff-header-row {
    display: grid;
    grid-template-columns: 200px 1fr 80px 1fr;
    padding: 10px 20px;
    color: #94a3b8;
    font-size: 13px;
    font-weight: bold;
}

.diff-card-row {
    display: grid;
    grid-template-columns: 200px 1fr 80px 1fr;
    align-items: center;
    background: white;
    padding: 20px;
    border-radius: 16px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    animation: slideIn 0.5s ease backwards;
}

.diff-card-row:hover {
    transform: scale(1.01);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

.diff-card-row.is-mismatch {
    border: 1px solid #fecaca;
    background: #fffafb;
}

.field-name {
    display: block;
    font-weight: 700;
    color: #334155;
}

.field-key {
    font-size: 11px;
    color: #94a3b8;
    font-family: monospace;
}

.val-display {
    font-family: 'JetBrains Mono', 'Monaco', monospace;
    font-size: 15px;
    font-weight: 600;
    color: #475569;
    word-break: break-all;
}

.is-mismatch .col-his .val-display {
    color: #ef4444;
    text-decoration: underline;
}

.col-indicator {
    display: flex;
    justify-content: center;
}

.match-icon {
    color: #22c55e;
    font-size: 24px;
}

.mismatch-icon {
    color: #ef4444;
    font-size: 24px;
    animation: pulse 2s infinite;
}

@keyframes slideIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes pulse {
    0% {
        transform: scale(1);
        opacity: 1;
    }

    50% {
        transform: scale(1.1);
        opacity: 0.7;
    }

    100% {
        transform: scale(1);
        opacity: 1;
    }
}

:global(html.dark) .diff-card-row {
    background: #1e293b;
}

:global(html.dark) .val-display {
    color: #cbd5e1;
}
</style>