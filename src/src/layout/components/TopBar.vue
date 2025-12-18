<template>
    <div class="top-bar">
        <div class="left">
            <!-- 折叠按钮 -->
            <el-button text @click="$emit('toggle-menu')">
                <el-icon><Menu /></el-icon>
            </el-button>
            <!-- 左侧：面包屑 -->
            <el-breadcrumb separator="/">
                <el-breadcrumb-item v-for="item in breadcrumb" :key="item.path">
                    {{ item.label }}
                </el-breadcrumb-item>
            </el-breadcrumb>
        </div>

        <!-- 右侧功能区 -->
        <div class="right">
            <!-- 系统状态 -->
            <div class="status">
                <span class="dot" :class="status"></span>
                {{ statusText }}
            </div>

            <!-- 当前时间 -->
            <div class="time">
                {{ currentTime }}
            </div>

            <!-- 全屏 -->
            <el-tooltip content="全屏">
                <el-button text @click="toggleFullscreen"> ⛶ </el-button>
            </el-tooltip>

            <!-- 用户 -->
            <el-dropdown>
                <span class="user-info">
                    <el-avatar :size="32" :src="user.avatar" />
                    <span class="name">{{ user.name }}</span>
                </span>

                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>个人信息</el-dropdown-item>
                        <el-dropdown-item divided @click="logout"> 退出登录 </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </div>
</template>

<script>
import { getBreadcrumb } from "../utils/breadcrumb";
import { Menu } from "@element-plus/icons-vue";

export default {
    name: "TopBar",
    components: { Menu },
    data() {
        return {
            user: {
                name: "Admin",
                avatar: "https://dummyimage.com/100x100"
            },
            status: "ok", // ok | warn | error
            currentTime: ""
        };
    },
    computed: {
        breadcrumb() {
            return getBreadcrumb(this.$route.path);
        },
        statusText() {
            return {
                ok: "系统正常",
                warn: "存在告警",
                error: "系统异常"
            }[this.status];
        }
    },
    mounted() {
        this.updateTime();
        this.timer = setInterval(this.updateTime, 1000);
    },
    beforeUnmount() {
        clearInterval(this.timer);
    },
    methods: {
        updateTime() {
            const now = new Date();
            const pad = n => String(n).padStart(2, "0");
            this.currentTime = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`;
        },
        toggleFullscreen() {
            if (!document.fullscreenElement) {
                document.documentElement.requestFullscreen();
            } else {
                document.exitFullscreen();
            }
        },
        logout() {
            console.log("logout");
            // 以后：清 token + router.push('/login')
        }
    }
};
</script>

<style scoped>
.left {
    display: flex;
    align-items: center;
    gap: 12px;
}

.top-bar {
    height: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.right {
    display: flex;
    align-items: center;
    gap: 12px;
}

.status {
    display: flex;
    align-items: center;
    font-size: 13px;
}

.dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-right: 4px;
}

.dot.ok {
    background: #67c23a;
}
.dot.warn {
    background: #e6a23c;
}
.dot.error {
    background: #f56c6c;
}

.time {
    font-size: 13px;
    color: #666;
}

.user-info {
    display: flex;
    align-items: center;
    cursor: pointer;
}
.name {
    margin-left: 6px;
}
</style>
