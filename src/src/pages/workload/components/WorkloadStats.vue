<template>
    <div class="stats-container">
        <div class="stats-cards-grid">
            <div v-for="item in stats" :key="item.title" class="modern-stat-card">
                <div class="card-icon" :style="{ background: item.color + '15', color: item.color }">
                    <el-icon :size="24">
                        <component :is="item.icon" />
                    </el-icon>
                </div>
                <div class="card-body">
                    <span class="label">{{ item.title }}</span>
                    <h2 class="value">{{ item.value }}</h2>
                </div>
            </div>
        </div>

        <div class="analysis-glass-card">
            <div class="chart-box">
                <svg viewBox="0 0 36 36" class="circular-chart">
                    <path class="circle-bg"
                        d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" />

                    <path class="circle inbound" :style="{ strokeDasharray: `${percents.inbound}, 100` }"
                        d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" />

                    <path class="circle outbound"
                        :style="{ strokeDasharray: `${percents.outbound}, 100`, strokeDashoffset: `-${percents.inbound}` }"
                        d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" />

                    <path class="circle return-path"
                        :style="{ strokeDasharray: `${percents.return}, 100`, strokeDashoffset: `-${percents.inbound + percents.outbound}` }"
                        d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" />
                </svg>
                <div class="chart-center">
                    <span class="center-label">工作占比</span>
                </div>
            </div>

            <div class="chart-legend">
                <div class="legend-item"><span class="dot inbound"></span> 入库 {{ percents.inbound }}%</div>
                <div class="legend-item"><span class="dot outbound"></span> 出库 {{ percents.outbound }}%</div>
                <div class="legend-item"><span class="dot return"></span> 退还 {{ percents.return }}%</div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from "vue";
import { Download, Upload, RefreshLeft, User } from "@element-plus/icons-vue";

const props = defineProps(['data']);

// 计算各项业务占总金额的百分比
const percents = computed(() => {
    const i = props.data.reduce((s, r) => s + (r.inbound?.reduce((ss, x) => ss + (x.totalAmount || 0), 0) || 0), 0);
    const o = props.data.reduce((s, r) => s + (r.outbound?.reduce((ss, x) => ss + (x.totalAmount || 0), 0) || 0), 0);
    const r = props.data.reduce((s, r) => s + (r.return?.reduce((ss, x) => ss + (x.totalAmount || 0), 0) || 0), 0);

    const total = i + o + r || 1; // 防止除以0

    return {
        inbound: Math.round((i / total) * 100),
        outbound: Math.round((o / total) * 100),
        return: Math.round((r / total) * 100)
    };
});

const stats = computed(() => {
    const sum = (type) => props.data.reduce((s, row) => s + (row[type]?.reduce((ss, i) => ss + (i.totalAmount || 0), 0) || 0), 0);
    const format = (val) => new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY' }).format(val);

    return [
        { title: "总入库金额", value: format(sum('inbound')), icon: Download, color: "#67C23A" },
        { title: "总出库金额", value: format(sum('outbound')), icon: Upload, color: "#F56C6C" },
        { title: "退还总额", value: format(sum('return')), icon: RefreshLeft, color: "#E6A23C" },
        { title: "操作员总数", value: props.data.length, icon: User, color: "#409EFF" }
    ];
});
</script>

<style scoped>
.stats-container {
    display: grid;
    grid-template-columns: 1fr 300px;
    gap: 24px;
    margin-bottom: 24px;
}

.stats-cards-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
}

.modern-stat-card {
    background: var(--bg-card, #fff);
    padding: 24px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
    transition: transform 0.3s;
}

.modern-stat-card:hover {
    transform: translateY(-5px);
}

.analysis-glass-card {
    background: var(--bg-card, #fff);
    padding: 24px;
    border-radius: 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

/* 环形图样式 */
.circular-chart {
    width: 140px;
    height: 140px;
    transform: rotate(-90deg);
}

/* 旋转90度让起始点在顶部 */
.circle-bg {
    fill: none;
    stroke: #f0f2f5;
    stroke-width: 3.5;
}

.circle {
    fill: none;
    stroke-width: 3.8;
    stroke-linecap: round;
    transition: all 0.5s ease;
}

.inbound {
    stroke: #67C23A;
}

.outbound {
    stroke: #F56C6C;
}

.return-path {
    stroke: #E6A23C;
}

/* 退还路径颜色 */

.chart-box {
    position: relative;
}

.chart-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    text-align: center;
}

.center-label {
    font-size: 13px;
    color: #909399;
    font-weight: 500;
}

.chart-legend {
    margin-top: 20px;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 14px;
    color: #606266;
}

.dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
}

.dot.inbound {
    background: #67C23A;
}

.dot.outbound {
    background: #F56C6C;
}

.dot.return {
    background: #E6A23C;
}

:global(html.dark) .modern-stat-card,
:global(html.dark) .analysis-glass-card {
    background: #1e1e1e;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

@media (max-width: 1200px) {
    .stats-container {
        grid-template-columns: 1fr;
    }

    .stats-cards-grid {
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    }
}
</style>