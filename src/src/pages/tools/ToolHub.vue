<template>
    <div class="tool-hub-page">
        <header class="hero-section">
            <div class="hero-content">
                <h1 class="welcome-text">智能工具箱</h1>
                <p class="date-text">集成高效的小工具，助力 HIS 系统维护与数据对账</p>
            </div>
        </header>

        <main class="tool-grid">
            <div v-for="(tool, index) in tools" :key="tool.path" class="tool-card glass-panel"
                :style="{ '--delay': index * 0.1 + 's' }" @click="navigateTo(tool)">
                <div v-if="tool.status === 'coming-soon'" class="status-tag">开发中</div>

                <div class="tool-icon-wrapper" :style="{ background: tool.color + '15', color: tool.color }">
                    <el-icon :size="32">
                        <component :is="tool.icon" />
                    </el-icon>
                </div>

                <div class="tool-info">
                    <h3 class="tool-name">{{ tool.title }}</h3>
                    <p class="tool-desc">{{ tool.description }}</p>
                </div>

                <div class="tool-footer">
                    <el-button :type="tool.status === 'active' ? 'primary' : 'info'" link
                        :icon="tool.status === 'active' ? 'Right' : 'Lock'">
                        {{ tool.status === 'active' ? '立即进入' : '暂未开放' }}
                    </el-button>
                </div>

                <el-icon class="bg-icon" :size="100">
                    <component :is="tool.icon" />
                </el-icon>
            </div>
        </main>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { Collection, Setting, Connection, Right, Lock } from "@element-plus/icons-vue";
import { ElMessage } from 'element-plus';

const router = useRouter();

const tools = [
    {
        title: "字典对比工具",
        description: "快速对比本地系统与 HIS 系统字典差异，支持名称、编码和拼音码检索。",
        icon: Collection,
        color: "#409EFF",
        path: "/tools/dict-compare",
        status: "active" // 已上线
    },
    {
        title: "接口重试管理",
        description: "监控配送与退库接口状态，支持手动触发失败任务的补偿重试。",
        icon: Connection,
        color: "#67C23A",
        path: "/tools/retry-manager",
        status: "coming-soon" // 开发中
    },
    {
        title: "系统参数校对",
        description: "核对本地运行参数与 HIS 全局配置的一致性，防止配置冲突。",
        icon: Setting,
        color: "#E6A23C",
        path: "/tools/config-check",
        status: "coming-soon" // 开发中
    }
];

const navigateTo = (tool) => {
    if (tool.status === 'coming-soon') {
        ElMessage.info(`工具「${tool.title}」正在研发中，敬请期待！`);
        return;
    }
    router.push(tool.path);
};
</script>

<style scoped>
.tool-hub-page {
    --hero-gradient: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    padding: 32px;
    background: var(--el-bg-color-page);
    min-height: 100%;
    display: flex;
    flex-direction: column;
    gap: 32px;
    box-sizing: border-box;
}

:global(html.dark) .tool-hub-page {
    --hero-gradient: linear-gradient(135deg, #1f2d5a 0%, #004e92 100%);
}

.hero-section {
    padding: 60px 48px;
    background: var(--hero-gradient);
    border-radius: 24px;
    color: white;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.welcome-text {
    font-size: 36px;
    font-weight: 800;
    margin: 0;
    text-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.date-text {
    font-size: 16px;
    opacity: 0.9;
    margin-top: 12px;
}

.tool-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 24px;
}

/* 玻璃拟态卡片优化 */
.tool-card {
    position: relative;
    padding: 36px;
    background: var(--el-bg-color);
    backdrop-filter: blur(12px);
    border: 1px solid var(--el-border-color-lighter) !important;
    border-radius: 24px;
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    gap: 20px;
    /* 入场动画 */
    animation: fadeInUp 0.6s ease forwards;
    animation-delay: var(--delay);
    opacity: 0;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(30px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.tool-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.08);
    border-color: var(--el-color-primary-light-5) !important;
}

/* 状态标签 */
.status-tag {
    position: absolute;
    top: 20px;
    right: -30px;
    background: #909399;
    color: white;
    font-size: 12px;
    padding: 4px 35px;
    transform: rotate(45deg);
    z-index: 10;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.tool-icon-wrapper {
    width: 68px;
    height: 68px;
    border-radius: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
    transition: transform 0.3s ease;
}

.tool-card:hover .tool-icon-wrapper {
    transform: scale(1.1) rotate(5deg);
}

.tool-name {
    font-size: 22px;
    font-weight: 700;
    color: var(--el-text-color-primary);
    margin: 0 0 10px;
}

.tool-desc {
    font-size: 14px;
    color: var(--el-text-color-secondary);
    line-height: 1.6;
}

.tool-footer {
    margin-top: auto;
    z-index: 1;
}

.bg-icon {
    position: absolute;
    right: -20px;
    bottom: -20px;
    opacity: 0.03;
    transform: rotate(-15deg);
    pointer-events: none;
    transition: all 0.5s ease;
}

.tool-card:hover .bg-icon {
    opacity: 0.08;
    transform: scale(1.2) rotate(0deg);
}

:global(html.dark) .tool-card {
    background: rgba(255, 255, 255, 0.03) !important;
}
</style>