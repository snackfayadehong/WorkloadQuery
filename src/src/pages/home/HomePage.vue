<template>
    <div class="home">
        <!-- 统计卡片 -->
        <el-row :gutter="16">
            <el-col :span="6">
                <el-card class="stat-card">
                    <div class="stat-title">今日查询</div>
                    <div class="stat-value">128</div>
                </el-card>
            </el-col>

            <el-col :span="6">
                <el-card class="stat-card">
                    <div class="stat-title">本月查询</div>
                    <div class="stat-value">3,542</div>
                </el-card>
            </el-col>

            <el-col :span="6">
                <el-card class="stat-card">
                    <div class="stat-title">失败次数</div>
                    <div class="stat-value error">4</div>
                </el-card>
            </el-col>

            <el-col :span="6">
                <el-card class="stat-card">
                    <div class="stat-title">系统状态</div>
                    <div class="stat-value success">正常</div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 内容区 -->
        <el-card class="main-card">
            <div class="placeholder">这里以后放：图表 / 列表 / 查询结果</div>
        </el-card>
    </div>
</template>

<script>
export default {
    name: "HomePage"
};
</script>

<style scoped>
.home {
    width: 100%;
}

.stat-card {
    text-align: center;
}

.stat-title {
    font-size: 14px;
    color: #909399;
}

.stat-value {
    margin-top: 8px;
    font-size: 28px;
    font-weight: bold;
    color: #303133;
}

.stat-value.success {
    color: #67c23a;
}

.stat-value.error {
    color: #f56c6c;
}

.main-card {
    margin-top: 16px;
}

.placeholder {
    height: 300px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #909399;
}
</style>
