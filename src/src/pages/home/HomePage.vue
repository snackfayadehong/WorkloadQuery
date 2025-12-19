<template>
    <div class="home-dashboard">
        <header class="hero-section">
            <div class="hero-content">
                <h1 class="welcome-text">{{ welcomeMessage }}</h1>
                <p class="date-text">{{ currentTime.format("YYYY年MM月DD日 dddd HH:mm:ss") }}</p>
                <div class="system-status-pill">
                    <span class="pulse-dot"></span>
                    后端 Go 服务：已连接 (172.21.1.75)
                </div>
            </div>
        </header>

        <div class="dashboard-grid">
            <section class="stats-row">
                <div v-for="item in stats" :key="item.title" class="modern-stat-card">
                    <div class="card-icon" :style="{ background: item.color + '15', color: item.color }">
                        <el-icon :size="24">
                            <component :is="item.icon" />
                        </el-icon>
                    </div>
                    <div class="card-body">
                        <span class="label">{{ item.title }}</span>
                        <h2 class="value">{{ item.value }}</h2>
                        <div class="trend-tag" :style="{ background: item.color + '10', color: item.color }">
                            {{ item.trend }}
                        </div>
                    </div>
                </div>
            </section>

            <div class="bottom-layout">
                <el-card shadow="never" class="glass-panel quick-panel">
                    <template #header>
                        <div class="panel-header">
                            <el-icon><Menu /></el-icon>
                            <span>快捷入口</span>
                        </div>
                    </template>
                    <div class="entry-grid">
                        <div class="entry-item" @click="$router.push('/workload')">
                            <div class="icon-circle blue">
                                <el-icon><Monitor /></el-icon>
                            </div>
                            <span class="entry-label">工作量查询</span>
                        </div>
                        <div class="entry-item">
                            <div class="icon-circle green">
                                <el-icon><Document /></el-icon>
                            </div>
                            <span class="entry-label">账单对账</span>
                        </div>
                        <div class="entry-item">
                            <div class="icon-circle orange">
                                <el-icon><Management /></el-icon>
                            </div>
                            <span class="entry-label">采购汇总</span>
                        </div>
                    </div>
                </el-card>

                <el-card shadow="never" class="glass-panel activity-panel">
                    <template #header>
                        <div class="panel-header">
                            <el-icon><Timer /></el-icon>
                            <span>系统日志</span>
                        </div>
                    </template>
                    <el-scrollbar height="220px">
                        <el-timeline>
                            <el-timeline-item timestamp="2025-12-18" type="primary" :hollow="true"> 完成 WorkloadPage 模块组件化重构 </el-timeline-item>
                            <el-timeline-item timestamp="2025-12-17"> 优化后端 WorkloadQuery 接口查询性能 </el-timeline-item>
                            <el-timeline-item timestamp="2025-12-16"> 财务对账系统数据导入完成 </el-timeline-item>
                        </el-timeline>
                    </el-scrollbar>
                </el-card>
            </div>
        </div>
    </div>
</template>

<script setup>
import { Monitor, Document, Menu, Timer, Management, Search, DataAnalysis, WarnTriangleFilled, User } from "@element-plus/icons-vue";
import { useHomeData } from "./composables/useHomeData";

const { stats, welcomeMessage, currentTime } = useHomeData();
</script>

<style scoped>
/* 1. 布局变量定义：实现明/暗一致性 */
.home-dashboard {
    --bg-page: #f6f8fb;
    --bg-card: #ffffff;
    --text-primary: #262626;
    --text-regular: #8c8c8c;
    --card-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
    --glass-bg: rgba(255, 255, 255, 0.7);
    --glass-border: rgba(255, 255, 255, 0.4);
    --hero-gradient: linear-gradient(135deg, #4158d0 0%, #c850c0 46%, #ffcc70 100%);

    padding: 32px;
    background: var(--bg-page);
    min-height: 100%;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    gap: 32px;
    transition: all 0.3s ease;
}

/* 暗黑模式变量覆盖 */
:global(html.dark) .home-dashboard {
    --bg-page: #121212;
    --bg-card: #1e1e1e;
    --text-primary: #e5eaf3;
    --text-regular: #a3a3a3;
    --card-shadow: none;
    --glass-bg: rgba(30, 30, 30, 0.6);
    --glass-border: rgba(255, 255, 255, 0.1);
    /* 暗黑模式下渐变调深，降低饱和度防止刺眼 */
    --hero-gradient: linear-gradient(135deg, #1f2d5a 0%, #6d2d6d 46%, #a37c32 100%);
}

/* 2. 欢迎看板样式 */
.hero-section {
    padding: 48px;
    background: var(--hero-gradient);
    border-radius: 24px;
    color: white;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    transition: background 0.5s ease;
}
.welcome-text {
    font-size: 36px;
    font-weight: 800;
    margin: 0;
}
.date-text {
    font-size: 16px;
    opacity: 0.8;
    margin: 12px 0 24px;
    font-family: monospace;
}

.system-status-pill {
    display: inline-flex;
    align-items: center;
    background: rgba(255, 255, 255, 0.15);
    padding: 8px 16px;
    border-radius: 50px;
    font-size: 14px;
    backdrop-filter: blur(8px);
    border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 3. 指标卡片样式 */
.dashboard-grid {
    display: flex;
    flex-direction: column;
    gap: 40px;
}
.stats-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 24px;
}

.modern-stat-card {
    background: var(--bg-card);
    padding: 24px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    gap: 20px;
    box-shadow: var(--card-shadow);
    border: 1px solid transparent;
    transition: all 0.3s ease;
}

:global(html.dark) .modern-stat-card {
    border-color: var(--glass-border);
}

.modern-stat-card:hover {
    transform: translateY(-8px);
    border-color: #409eff;
}

.card-icon {
    width: 56px;
    height: 56px;
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
}
.label {
    font-size: 14px;
    color: var(--text-regular);
}
.value {
    font-size: 28px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 4px 0;
}

/* 4. 底部面板与玻璃拟态 */
.bottom-layout {
    display: grid;
    grid-template-columns: 1fr 1.5fr;
    gap: 32px;
}

.glass-panel {
    background: var(--glass-bg) !important;
    backdrop-filter: blur(12px);
    border: 1px solid var(--glass-border) !important;
    border-radius: 20px;
}

.panel-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 700;
    color: var(--text-primary);
}

.entry-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 15px;
}
.entry-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    cursor: pointer;
    padding: 20px;
    border-radius: 16px;
    transition: all 0.2s;
}
.entry-item:hover {
    background: rgba(64, 158, 255, 0.1);
}
.entry-label {
    color: var(--text-regular);
    font-size: 14px;
}

.icon-circle {
    width: 52px;
    height: 52px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 22px;
}

/* 渐变图标适配暗黑模式 */
:global(html.dark) .icon-circle {
    filter: brightness(0.85);
}
.icon-circle.blue {
    background: linear-gradient(135deg, #667eea, #764ba2);
}
.icon-circle.green {
    background: linear-gradient(135deg, #42e695, #3bb2b8);
}
.icon-circle.orange {
    background: linear-gradient(135deg, #f6d365, #fda085);
}

/* 动画部分 */
.pulse-dot {
    width: 8px;
    height: 8px;
    background: #52c41a;
    border-radius: 50%;
    margin-right: 10px;
    animation: pulse 2s infinite;
}

@keyframes pulse {
    0% {
        box-shadow: 0 0 0 0 rgba(82, 196, 26, 0.5);
    }
    70% {
        box-shadow: 0 0 0 8px rgba(82, 196, 26, 0);
    }
    100% {
        box-shadow: 0 0 0 0 rgba(82, 196, 26, 0);
    }
}

@media (max-width: 1100px) {
    .bottom-layout {
        grid-template-columns: 1fr;
    }
}
</style>
