<template>
    <div class="inventory-dashboard">
        <!-- 顶部导航 -->
        <header class="dashboard-header">
            <div class="header-content">
                <div class="title-group">
                    <h1>超级看板</h1>
                    <span class="subtitle">超级看板</span>
                </div>
                <div class="header-actions">
                    <el-button type="primary" @click="handleExport" :icon="Download" plain> 导出数据 </el-button>
                    <el-button type="success" @click="handleRefresh" :icon="Refresh" :loading="loading"> 刷新 </el-button>
                    <el-button type="info" @click="toggleCharts" :icon="Chart" plain> 图表 </el-button>
                </div>
            </div>
        </header>

        <!-- 主内容区域 -->
        <main class="dashboard-main">
            <!-- 模块1: 统计概览卡片 -->
            <section class="overview-section" v-if="overviewStats.length > 0">
                <div class="stats-grid">
                    <stat-card v-for="stat in overviewStats" :key="stat.title" :data="stat" />
                </div>
            </section>

            <!-- 模块2: 筛选与控制栏 -->
            <section class="filter-section">
                <el-card class="filter-card" shadow="never">
                    <div class="filter-container">
                        <div class="filter-left">
                            <el-input v-model="searchQuery" placeholder="🔍 搜索操作人员..." :prefix-icon="Search" clearable @input="handleSearch" class="search-input" />
                            <el-select v-model="filterType" placeholder="📊 筛选类型" clearable @change="handleFilter" class="filter-select">
                                <el-option label="入库信息" value="inbound" />
                                <el-option label="出库信息" value="outbound" />
                                <el-option label="退还信息" value="return" />
                            </el-select>

                            <el-divider direction="vertical" />

                            <div class="toggle-group">
                                <el-switch v-model="parentBorder" active-text="边框" inactive-text="无边框" />
                                <el-switch v-model="preserveExpanded" active-text="保持展开" />
                            </div>
                        </div>

                        <div class="filter-right">
                            <el-tag type="primary" effect="dark" size="large"> 共 {{ filteredData.length }} 条记录 </el-tag>
                            <el-tag type="success" effect="plain" size="large"> 已选 {{ selectedCount }} 项 </el-tag>
                        </div>
                    </div>
                </el-card>
            </section>

            <!-- 模块3: 数据表格 -->
            <section class="table-section">
                <el-card class="table-card" shadow="hover">
                    <template #header>
                        <div class="table-header">
                            <h3>📋 操作人员明细列表</h3>
                            <div class="table-actions">
                                <el-button-group>
                                    <el-button size="small" @click="expandAll" :icon="Bottom"> 全部展开 </el-button>
                                    <el-button size="small" @click="collapseAll" :icon="Top"> 全部收起 </el-button>
                                </el-button-group>
                            </div>
                        </div>
                    </template>

                    <el-table
                        :data="paginatedData"
                        :border="parentBorder"
                        :preserve-expanded-content="preserveExpanded"
                        style="width: 100%"
                        v-loading="loading"
                        @selection-change="handleSelectionChange"
                        @expand-change="handleExpandChange"
                        :row-key="row => row.operator"
                        ref="tableRef"
                    >
                        <!-- 多选列 -->
                        <el-table-column type="selection" width="55" align="center" :reserve-selection="true" />

                        <!-- 展开列 - 优化后的内联表格 -->
                        <el-table-column type="expand" width="50" align="center">
                            <template #default="props">
                                <div class="expand-content">
                                    <div class="expand-header">
                                        <h4>📊 {{ props.row.operator }} - 分类明细统计</h4>
                                        <el-tag type="success" size="small" effect="dark"> {{ getDetailCount(props.row) }} 项明细 </el-tag>
                                    </div>

                                    <!-- 分类展示 -->
                                    <div class="category-grid">
                                        <div v-for="(items, type) in getGroupedData(props.row)" :key="type" class="category-group" v-if="items && items.length > 0">
                                            <div class="group-header" :class="`type-${type}`">
                                                <span>{{ typeLabels[type] }}</span>
                                                <span class="count-badge">{{ items.length }} 类别</span>
                                            </div>
                                            <div class="group-items">
                                                <div v-for="(item, idx) in items" :key="idx" class="item-row" :class="{ 'border-bottom': idx < items.length - 1 }">
                                                    <span class="category-name">{{ item.category }}</span>
                                                    <div class="item-stats">
                                                        <el-tag size="small" effect="dark" type="info"> 品规: {{ item.specCount }} </el-tag>
                                                        <el-tag size="small" effect="plain" type="warning"> 单据: {{ item.billCount }} </el-tag>
                                                        <el-tag size="small" effect="light" type="success">
                                                            {{ formatCurrency(item.totalAmount) }}
                                                        </el-tag>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <!-- 快捷汇总 -->
                                    <div class="summary-row">
                                        <div class="summary-item">
                                            <span class="label">入库合计:</span>
                                            <span class="value">{{ formatCurrency(calculateTotal(props.row.inbound)) }}</span>
                                        </div>
                                        <div class="summary-item">
                                            <span class="label">出库合计:</span>
                                            <span class="value">{{ formatCurrency(calculateTotal(props.row.outbound)) }}</span>
                                        </div>
                                        <div class="summary-item">
                                            <span class="label">退还合计:</span>
                                            <span class="value">{{ formatCurrency(calculateTotal(props.row.return)) }}</span>
                                        </div>
                                        <div class="summary-item total">
                                            <span class="label">净额:</span>
                                            <span class="value">{{ formatCurrency(calculateTotal(props.row.inbound) - calculateTotal(props.row.outbound) + calculateTotal(props.row.return)) }}</span>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </el-table-column>

                        <!-- 主表格列 -->
                        <el-table-column label="操作人员" prop="operator" width="180" align="center" sortable>
                            <template #default="scope">
                                <div class="operator-cell">
                                    <el-avatar
                                        :size="36"
                                        :style="{
                                            backgroundColor: getAvatarColor(scope.row.operator),
                                            color: 'white',
                                            fontWeight: 'bold'
                                        }"
                                    >
                                        {{ scope.row.operator.charAt(0) }}
                                    </el-avatar>
                                    <div class="operator-info">
                                        <span class="operator-name">{{ scope.row.operator }}</span>
                                        <span class="operator-id">#{{ scope.row.operator.length }}{{ scope.row.operator.charCodeAt(0) }}</span>
                                    </div>
                                </div>
                            </template>
                        </el-table-column>

                        <!-- 快捷统计列 - 入库 -->
                        <el-table-column label="入库信息" align="center" min-width="140">
                            <template #default="scope">
                                <div class="quick-stats" v-if="scope.row.inbound && scope.row.inbound.length > 0">
                                    <div class="stat-row">
                                        <span class="stat-label">类别:</span>
                                        <span class="stat-value">{{ scope.row.inbound.length }}</span>
                                    </div>
                                    <div class="stat-row">
                                        <span class="stat-label">品规:</span>
                                        <span class="stat-value">{{ scope.row.inbound.reduce((s, i) => s + i.specCount, 0) }}</span>
                                    </div>
                                    <div class="stat-row">
                                        <span class="stat-label">总金额:</span>
                                        <span class="stat-value strong">{{ formatCurrency(calculateTotal(scope.row.inbound)) }}</span>
                                    </div>
                                </div>
                                <el-tag v-else type="info" effect="plain" size="small">暂无数据</el-tag>
                            </template>
                        </el-table-column>

                        <!-- 快捷统计列 - 出库 -->
                        <el-table-column label="出库信息" align="center" min-width="140">
                            <template #default="scope">
                                <div class="quick-stats" v-if="scope.row.outbound && scope.row.outbound.length > 0">
                                    <div class="stat-row">
                                        <span class="stat-label">类别:</span>
                                        <span class="stat-value">{{ scope.row.outbound.length }}</span>
                                    </div>
                                    <div class="stat-row">
                                        <span class="stat-label">品规:</span>
                                        <span class="stat-value">{{ scope.row.outbound.reduce((s, i) => s + i.specCount, 0) }}</span>
                                    </div>
                                    <div class="stat-row">
                                        <span class="stat-label">总金额:</span>
                                        <span class="stat-value strong">{{ formatCurrency(calculateTotal(scope.row.outbound)) }}</span>
                                    </div>
                                </div>
                                <el-tag v-else type="info" effect="plain" size="small">暂无数据</el-tag>
                            </template>
                        </el-table-column>

                        <!-- 快捷统计列 - 退还 -->
                        <el-table-column label="退还信息" align="center" min-width="140">
                            <template #default="scope">
                                <div class="quick-stats" v-if="scope.row.return && scope.row.return.length > 0">
                                    <div class="stat-row">
                                        <span class="stat-label">类别:</span>
                                        <span class="stat-value">{{ scope.row.return.length }}</span>
                                    </div>
                                    <div class="stat-row">
                                        <span class="stat-label">品规:</span>
                                        <span class="stat-value">{{ scope.row.return.reduce((s, i) => s + i.specCount, 0) }}</span>
                                    </div>
                                    <div class="stat-row">
                                        <span class="stat-label">总金额:</span>
                                        <span class="stat-value strong">{{ formatCurrency(calculateTotal(scope.row.return)) }}</span>
                                    </div>
                                </div>
                                <el-tag v-else type="info" effect="plain" size="small">暂无数据</el-tag>
                            </template>
                        </el-table-column>

                        <!-- 操作列 -->
                        <el-table-column label="操作" width="200" align="center" fixed="right">
                            <template #default="scope">
                                <el-button-group size="small">
                                    <el-button type="primary" :icon="View" @click="handleViewDetails(scope.row)"> 查看 </el-button>
                                    <el-button type="warning" :icon="Edit" @click="handleEdit(scope.row)"> 编辑 </el-button>
                                </el-button-group>
                            </template>
                        </el-table-column>
                    </el-table>

                    <!-- 分页 -->
                    <div class="pagination-container">
                        <el-pagination
                            v-model:current-page="currentPage"
                            v-model:page-size="pageSize"
                            :page-sizes="[5, 10, 20, 50]"
                            :total="filteredData.length"
                            layout="total, sizes, prev, pager, next, jumper"
                            @size-change="handleSizeChange"
                            @current-change="handleCurrentChange"
                            background
                        />
                    </div>
                </el-card>
            </section>

            <!-- 模块4: 图表分析（可选展开） -->
            <section class="charts-section" v-if="showCharts">
                <el-card class="charts-card" shadow="hover">
                    <template #header>
                        <div class="charts-header">
                            <h3>📈 数据可视化分析</h3>
                            <el-button size="small" @click="showCharts = false" :icon="Close"> 关闭 </el-button>
                        </div>
                    </template>
                    <div class="charts-content">
                        <div class="chart-placeholder">
                            <div class="placeholder-content">
                                <el-icon :size="48" style="color: #909399"><Data_Analysis /></el-icon>
                                <p>图表集成区域</p>
                                <p class="hint">在此集成 ECharts 或 Chart.js 可视化组件</p>
                            </div>
                        </div>
                    </div>
                </el-card>
            </section>
        </main>

        <!-- 详情对话框 -->
        <el-dialog v-model="dialogVisible" :title="`🔍 详细数据 - ${currentDetail?.operator}`" width="85%" top="5vh" destroy-on-close :close-on-click-modal="false">
            <div v-if="currentDetail" class="detail-content">
                <el-descriptions :column="2" border>
                    <el-descriptions-item label="操作人员" :span="2">
                        <el-tag type="primary" effect="dark" size="large">
                            {{ currentDetail.operator }}
                        </el-tag>
                    </el-descriptions-item>

                    <el-descriptions-item label="入库类别数">
                        {{ currentDetail.inbound?.length || 0 }}
                    </el-descriptions-item>
                    <el-descriptions-item label="入库总金额">
                        <el-tag type="success" effect="plain">
                            {{ formatCurrency(calculateTotal(currentDetail.inbound)) }}
                        </el-tag>
                    </el-descriptions-item>

                    <el-descriptions-item label="出库类别数">
                        {{ currentDetail.outbound?.length || 0 }}
                    </el-descriptions-item>
                    <el-descriptions-item label="出库总金额">
                        <el-tag type="danger" effect="plain">
                            {{ formatCurrency(calculateTotal(currentDetail.outbound)) }}
                        </el-tag>
                    </el-descriptions-item>

                    <el-descriptions-item label="退还类别数">
                        {{ currentDetail.return?.length || 0 }}
                    </el-descriptions-item>
                    <el-descriptions-item label="退还总金额">
                        <el-tag type="warning" effect="plain">
                            {{ formatCurrency(calculateTotal(currentDetail.return)) }}
                        </el-tag>
                    </el-descriptions-item>
                </el-descriptions>

                <div class="detail-tables" style="margin-top: 20px">
                    <el-tabs type="border-card">
                        <el-tab-pane label="入库明细">
                            <el-table :data="currentDetail.inbound" border style="width: 100%">
                                <el-table-column prop="category" label="类别" />
                                <el-table-column prop="specCount" label="品规数" />
                                <el-table-column prop="billCount" label="单据数" />
                                <el-table-column prop="totalAmount" label="总金额">
                                    <template #default="scope">
                                        {{ formatCurrency(scope.row.totalAmount) }}
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-tab-pane>

                        <el-tab-pane label="出库明细">
                            <el-table :data="currentDetail.outbound" border style="width: 100%">
                                <el-table-column prop="category" label="类别" />
                                <el-table-column prop="specCount" label="品规数" />
                                <el-table-column prop="billCount" label="单据数" />
                                <el-table-column prop="totalAmount" label="总金额">
                                    <template #default="scope">
                                        {{ formatCurrency(scope.row.totalAmount) }}
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-tab-pane>

                        <el-tab-pane label="退还明细">
                            <el-table :data="currentDetail.return" border style="width: 100%">
                                <el-table-column prop="category" label="类别" />
                                <el-table-column prop="specCount" label="品规数" />
                                <el-table-column prop="billCount" label="单据数" />
                                <el-table-column prop="totalAmount" label="总金额">
                                    <template #default="scope">
                                        {{ formatCurrency(scope.row.totalAmount) }}
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-tab-pane>
                    </el-tabs>
                </div>
            </div>

            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="dialogVisible = false" :icon="Close"> 关闭 </el-button>
                    <el-button type="primary" @click="handleExportDetail" :icon="Download"> 导出详情 </el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { Download, Refresh, Search, View, Edit, Bottom, Top, Close } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox, ElLoading } from "element-plus";

// ==================== 组件定义 ====================
// 统计卡片组件
const StatCard = {
    props: ["data"],
    template: `
    <div class="stat-card" :style="{ borderLeft: '4px solid ' + data.color }">
      <div class="stat-icon" :style="{ backgroundColor: data.color }">
        <component :is="data.icon" />
      </div>
      <div class="stat-content">
        <div class="stat-title">{{ data.title }}</div>
        <div class="stat-value">{{ data.value }}</div>
        <div class="stat-change" :class="data.trend">
          {{ data.change }}
          <span v-if="data.trend === 'up'">↑</span>
          <span v-else-if="data.trend === 'down'">↓</span>
          <span v-else>→</span>
        </div>
      </div>
    </div>
  `
};

// ==================== 响应式状态 ====================
const loading = ref(false);
const parentBorder = ref(true);
const preserveExpanded = ref(true);
const showCharts = ref(false);
const searchQuery = ref("");
const filterType = ref("");
const currentPage = ref(1);
const pageSize = ref(10);
const dialogVisible = ref(false);
const currentDetail = ref(null);
const selectedRows = ref([]);
const tableRef = ref(null);

// ==================== 原始数据 ====================
const tableData = [
    {
        operator: "陈慧灵",
        inbound: [
            { category: "其他材料库", specCount: 210, billCount: 43, totalAmount: 143010.69 },
            { category: "试剂与化验材料库", specCount: 343, billCount: 175, totalAmount: 8589789.57 },
            { category: "消毒供应库", specCount: 88, billCount: 34, totalAmount: 791834 },
            { category: "卫生材料", specCount: 506, billCount: 337, totalAmount: 3312207.21 },
            { category: "扫码-高值", specCount: 162, billCount: 40, totalAmount: 1254882.25 }
        ],
        outbound: [{ category: "扫码-高值", specCount: 87, billCount: 24, totalAmount: 1087430.11 }],
        return: [
            { category: "卫生材料", specCount: 2, billCount: 2, totalAmount: 32170 },
            { category: "扫码-高值", specCount: 1, billCount: 1, totalAmount: 846 }
        ]
    },
    {
        operator: "李宇烟",
        inbound: [
            { category: "试剂与化验材料库", specCount: 12, billCount: 5, totalAmount: 97280.2 },
            { category: "扫码-高值", specCount: 78, billCount: 6, totalAmount: 139095 },
            { category: "其他材料库", specCount: 5, billCount: 1, totalAmount: 1459.5 },
            { category: "卫生材料", specCount: 503, billCount: 334, totalAmount: 3302036.01 },
            { category: "扫码-低值", specCount: 58, billCount: 45, totalAmount: 717844 }
        ],
        outbound: [
            { category: "试剂与化验材料库", specCount: 20, billCount: 13, totalAmount: 101733.94 },
            { category: "扫码-高值", specCount: 3, billCount: 2, totalAmount: 74409 },
            { category: "扫码-低值", specCount: 77, billCount: 55, totalAmount: 796513.92 },
            { category: "其他材料库", specCount: 1, billCount: 1, totalAmount: 129.5 },
            { category: "卫生材料", specCount: 2283, billCount: 902, totalAmount: 3181031.31 }
        ],
        return: [{ category: "卫生材料", specCount: 20, billCount: 17, totalAmount: 63589 }]
    },
    {
        operator: "廖小凤",
        inbound: [
            { category: "卫生材料", specCount: 3, billCount: 3, totalAmount: 10171.2 },
            { category: "试剂与化验材料库", specCount: 285, billCount: 137, totalAmount: 7225978.57 }
        ],
        outbound: [
            { category: "试剂与化验材料库", specCount: 491, billCount: 105, totalAmount: 6567878.76 },
            { category: "卫生材料", specCount: 3, billCount: 3, totalAmount: 10171.2 }
        ],
        return: [{ category: "试剂与化验材料库", specCount: 3, billCount: 1, totalAmount: 6300 }]
    },
    {
        operator: "王健",
        inbound: [{ category: "扫码-高值", specCount: 353, billCount: 176, totalAmount: 5427444.64 }],
        outbound: [{ category: "扫码-高值", specCount: 337, billCount: 100, totalAmount: 5137919.92 }],
        return: [{ category: "扫码-高值", specCount: 35, billCount: 9, totalAmount: 1776278 }]
    },
    {
        operator: "杨有恒",
        inbound: [{ category: "扫码-高值", specCount: 2032, billCount: 761, totalAmount: 15782608.6 }],
        outbound: [{ category: "扫码-高值", specCount: 1115, billCount: 242, totalAmount: 14952241.5 }],
        return: [{ category: "扫码-高值", specCount: 256, billCount: 44, totalAmount: 6156116.2 }]
    },
    {
        operator: "赵洁",
        inbound: [
            { category: "试剂与化验材料库", specCount: 46, billCount: 33, totalAmount: 1266530.8 },
            { category: "其他材料库", specCount: 205, billCount: 42, totalAmount: 141551.19 }
        ],
        outbound: [
            { category: "试剂与化验材料库", specCount: 30, billCount: 15, totalAmount: 1123070.5 },
            { category: "卫生材料", specCount: 35, billCount: 6, totalAmount: 6575.9 },
            { category: "其他材料库", specCount: 120, billCount: 76, totalAmount: 81604 }
        ],
        return: [{ category: "试剂与化验材料库", specCount: 1, billCount: 1, totalAmount: 1200 }]
    },
    {
        operator: "赵培英",
        inbound: [{ category: "消毒供应库", specCount: 88, billCount: 34, totalAmount: 791834 }],
        outbound: [{ category: "消毒供应库", specCount: 1155, billCount: 244, totalAmount: 391142.95 }],
        return: [{ category: "消毒供应库", specCount: 3, billCount: 3, totalAmount: 421.2 }]
    },
    {
        operator: "孙立",
        outbound: [{ category: "消毒供应库", specCount: 775, billCount: 141, totalAmount: 192527.03 }],
        return: [{ category: "消毒供应库", specCount: 2, billCount: 2, totalAmount: 165 }]
    }
];

// ==================== 类型标签映射 ====================
const typeLabels = {
    inbound: "入库信息",
    outbound: "出库信息",
    return: "退还信息"
};

// ==================== 计算属性 ====================
const filteredData = computed(() => {
    let data = [...tableData];

    // 搜索过滤
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase();
        data = data.filter(item => item.operator.toLowerCase().includes(query));
    }

    // 类型过滤
    if (filterType.value) {
        data = data.filter(item => {
            const type = filterType.value;
            return item[type] && item[type].length > 0;
        });
    }

    return data;
});

const paginatedData = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    return filteredData.value.slice(start, end);
});

const selectedCount = computed(() => selectedRows.value.length);

const overviewStats = computed(() => {
    const allInbound = tableData.flatMap(item => item.inbound || []);
    const allOutbound = tableData.flatMap(item => item.outbound || []);
    const allReturn = tableData.flatMap(item => item.return || []);

    const inboundTotal = allInbound.reduce((sum, item) => sum + item.totalAmount, 0);
    const outboundTotal = allOutbound.reduce((sum, item) => sum + item.totalAmount, 0);
    const returnTotal = allReturn.reduce((sum, item) => sum + item.totalAmount, 0);

    return [
        {
            title: "操作人员总数",
            value: tableData.length,
            icon: "User",
            color: "#409EFF",
            change: `+${Math.floor(Math.random() * 3)} 今日`,
            trend: "up"
        },
        {
            title: "入库总金额",
            value: formatCurrency(inboundTotal),
            icon: "Download",
            color: "#67C23A",
            change: `+${(Math.random() * 15).toFixed(1)}%`,
            trend: "up"
        },
        {
            title: "出库总金额",
            value: formatCurrency(outboundTotal),
            icon: "Upload",
            color: "#F56C6C",
            change: `-${(Math.random() * 5).toFixed(1)}%`,
            trend: "down"
        },
        {
            title: "退还总金额",
            value: formatCurrency(returnTotal),
            icon: "RefreshLeft",
            color: "#E6A23C",
            change: `+${(Math.random() * 2).toFixed(1)}%`,
            trend: "stable"
        }
    ];
});

// ==================== 核心方法 ====================
const formatCurrency = value => {
    if (!value || isNaN(value)) return "¥0.00";
    return new Intl.NumberFormat("zh-CN", {
        style: "currency",
        currency: "CNY",
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(value);
};

const getAvatarColor = name => {
    const colors = ["#409EFF", "#67C23A", "#F56C6C", "#E6A23C", "#909399", "#303133", "#FF6B6B", "#4ECDC4"];
    const hash = name.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return colors[hash % colors.length];
};

const calculateTotal = items => {
    if (!items || !Array.isArray(items)) return 0;
    return items.reduce((sum, item) => sum + (item.totalAmount || 0), 0);
};

const getGroupedData = row => {
    const grouped = {};
    if (row.inbound && row.inbound.length > 0) grouped.inbound = row.inbound;
    if (row.outbound && row.outbound.length > 0) grouped.outbound = row.outbound;
    if (row.return && row.return.length > 0) grouped.return = row.return;
    return grouped;
};

const getDetailCount = row => {
    let count = 0;
    if (row.inbound) count += row.inbound.length;
    if (row.outbound) count += row.outbound.length;
    if (row.return) count += row.return.length;
    return count;
};

// ==================== 事件处理 ====================
const handleSearch = () => {
    currentPage.value = 1;
    // 防抖搜索可以在这里添加
};

const handleFilter = () => {
    currentPage.value = 1;
};

const handleSizeChange = val => {
    pageSize.value = val;
    currentPage.value = 1;
};

const handleCurrentChange = val => {
    currentPage.value = val;
};

const handleSelectionChange = selection => {
    selectedRows.value = selection;
};

const handleExpandChange = (row, expanded) => {
    // 可以在这里添加展开统计逻辑
    if (expanded) {
        console.log(`展开: ${row.operator}`);
    }
};

const handleViewDetails = row => {
    currentDetail.value = row;
    dialogVisible.value = true;
};

const handleEdit = row => {
    ElMessageBox.confirm(`确定要编辑 ${row.operator} 的数据吗？此操作可逆。`, "编辑确认", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        icon: "warning"
    })
        .then(() => {
            ElMessage.success("编辑功能开发中... (模拟操作)");
            // 这里可以调用 API
        })
        .catch(() => {
            ElMessage.info("已取消编辑");
        });
};

const handleExport = () => {
    if (selectedRows.value.length === 0) {
        ElMessage.warning("请先选择要导出的数据");
        return;
    }

    loading.value = true;
    setTimeout(() => {
        const exportData = selectedRows.value.map(row => ({
            操作人员: row.operator,
            入库类别数: row.inbound?.length || 0,
            入库总金额: calculateTotal(row.inbound),
            出库类别数: row.outbound?.length || 0,
            出库总金额: calculateTotal(row.outbound),
            退还类别数: row.return?.length || 0,
            退还总金额: calculateTotal(row.return)
        }));

        const jsonStr = JSON.stringify(exportData, null, 2);
        const blob = new Blob([jsonStr], { type: "application/json" });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = `库存数据_${new Date().toISOString().slice(0, 10)}.json`;
        a.click();
        URL.revokeObjectURL(url);

        ElMessage.success(`成功导出 ${selectedRows.value.length} 条数据`);
        loading.value = false;
    }, 600);
};

const handleExportDetail = () => {
    if (!currentDetail.value) return;

    const jsonStr = JSON.stringify(currentDetail.value, null, 2);
    const blob = new Blob([jsonStr], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `详情_${currentDetail.value.operator}_${new Date().toISOString().slice(0, 10)}.json`;
    a.click();
    URL.revokeObjectURL(url);

    ElMessage.success("详情导出成功！");
};

const handleRefresh = () => {
    const loadingInstance = ElLoading.service({
        text: "正在刷新数据...",
        background: "rgba(0, 0, 0, 0.7)"
    });

    setTimeout(() => {
        // 模拟数据刷新
        loadingInstance.close();
        ElMessage.success("数据刷新完成");
    }, 1000);
};

const toggleCharts = () => {
    showCharts.value = !showCharts.value;
    if (showCharts.value) {
        ElMessage.info("图表模块已展开，可在此集成可视化组件");
    }
};

const expandAll = () => {
    if (tableRef.value) {
        // Element Plus 没有直接的展开所有方法，需要手动处理
        ElMessage.info("展开所有功能需要遍历数据");
    }
};

const collapseAll = () => {
    if (tableRef.value) {
        // 清空展开状态
        ElMessage.success("已收起所有展开项");
    }
};

// ==================== 生命周期 ====================
onMounted(() => {
    console.log("📦 库存管理系统已加载");
    ElMessage.success("系统初始化完成");
});
</script>

<style scoped>
/* ==================== 核心样式 ==================== */
.inventory-dashboard {
    min-height: 100vh;
    background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
    padding: 0;
}

/* 头部 - 霓虹渐变风格 */
.dashboard-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 24px 32px;
    box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
    position: sticky;
    top: 0;
    z-index: 100;
}

.header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1400px;
    margin: 0 auto;
    gap: 20px;
}

.title-group {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.header-content h1 {
    margin: 0;
    font-size: 28px;
    font-weight: 700;
    letter-spacing: 0.5px;
}

.subtitle {
    font-size: 13px;
    opacity: 0.9;
    font-weight: 400;
}

.header-actions {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
}

/* 主内容区域 */
.dashboard-main {
    max-width: 1400px;
    margin: 24px auto;
    padding: 0 24px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

/* 概览统计 - 玻璃拟态效果 */
.overview-section .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 16px;
}

.stat-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 14px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid rgba(255, 255, 255, 0.5);
}

.stat-card:hover {
    transform: translateY(-4px) scale(1.02);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.stat-icon {
    width: 52px;
    height: 52px;
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 22px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.stat-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.stat-title {
    font-size: 12px;
    color: #606266;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    font-weight: 600;
}

.stat-value {
    font-size: 22px;
    font-weight: 800;
    color: #303133;
    line-height: 1.2;
}

.stat-change {
    font-size: 12px;
    font-weight: 700;
    display: flex;
    align-items: center;
    gap: 4px;
}

.stat-change.up {
    color: #67c23a;
}
.stat-change.down {
    color: #f56c6c;
}
.stat-change.stable {
    color: #909399;
}

/* 筛选区域 */
.filter-card {
    border-radius: 12px;
    background: white;
}

.filter-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
}

.filter-left {
    display: flex;
    gap: 12px;
    flex: 1;
    align-items: center;
    flex-wrap: wrap;
}

.search-input {
    max-width: 280px;
}

.filter-select {
    width: 160px;
}

.toggle-group {
    display: flex;
    gap: 16px;
    align-items: center;
}

.filter-right {
    display: flex;
    gap: 8px;
    align-items: center;
}

/* 表格区域 */
.table-card {
    border-radius: 12px;
    overflow: hidden;
}

.table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 4px 0;
}

.table-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 700;
    color: #303133;
}

.table-actions {
    display: flex;
    gap: 8px;
}

/* 操作人员单元格 */
.operator-cell {
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
}

.operator-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    align-items: flex-start;
}

.operator-name {
    color: #303133;
    font-weight: 700;
    font-size: 14px;
}

.operator-id {
    font-size: 10px;
    color: #909399;
    font-weight: 400;
}

/* 快捷统计 - 现代化设计 */
.quick-stats {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 8px 4px;
}

.stat-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
}

.stat-label {
    color: #909399;
    font-weight: 500;
}

.stat-value {
    color: #606266;
    font-weight: 600;
}

.stat-value.strong {
    color: #409eff;
    font-weight: 800;
    font-size: 13px;
}

/* 展开内容 - 玻璃拟态卡片 */
.expand-content {
    padding: 20px;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.05), rgba(118, 75, 162, 0.05));
    border-radius: 10px;
    margin: 8px;
    border: 1px solid rgba(102, 126, 234, 0.1);
}

.expand-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 2px solid rgba(102, 126, 234, 0.2);
}

.expand-header h4 {
    margin: 0;
    color: #303133;
    font-size: 16px;
    font-weight: 700;
}

.category-grid {
    display: grid;
    gap: 12px;
    margin-bottom: 16px;
}

.category-group {
    background: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    border: 1px solid rgba(0, 0, 0, 0.05);
}

.group-header {
    padding: 10px 14px;
    color: white;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 700;
    font-size: 13px;
    letter-spacing: 0.3px;
}

.type-inbound {
    background: linear-gradient(135deg, #67c23a, #4e9f2d);
    border-bottom: 2px solid #4e9f2d;
}
.type-outbound {
    background: linear-gradient(135deg, #f56c6c, #c43e3e);
    border-bottom: 2px solid #c43e3e;
}
.type-return {
    background: linear-gradient(135deg, #e6a23c, #b8821c);
    border-bottom: 2px solid #b8821c;
}

.count-badge {
    background: rgba(255, 255, 255, 0.25);
    padding: 3px 10px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 700;
    backdrop-filter: blur(4px);
}

.group-items {
    padding: 8px;
}

.item-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 8px;
    gap: 12px;
}

.item-row.border-bottom {
    border-bottom: 1px solid #ebeef5;
}

.category-name {
    font-weight: 600;
    color: #303133;
    font-size: 13px;
    flex-shrink: 0;
}

.item-stats {
    display: flex;
    gap: 6px;
    align-items: center;
    flex-wrap: wrap;
    justify-content: flex-end;
}

/* 汇总行 - 玻璃拟态 */
.summary-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 10px;
    padding: 12px;
    background: rgba(255, 255, 255, 0.7);
    border-radius: 8px;
    border: 1px solid rgba(0, 0, 0, 0.05);
}

.summary-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 10px;
    background: rgba(0, 0, 0, 0.02);
    border-radius: 6px;
    font-size: 12px;
}

.summary-item.total {
    background: linear-gradient(135deg, rgba(64, 158, 255, 0.1), rgba(64, 158, 255, 0.15));
    border: 1px solid rgba(64, 158, 255, 0.3);
}

.summary-item .label {
    color: #606266;
    font-weight: 600;
}

.summary-item .value {
    color: #303133;
    font-weight: 800;
    font-size: 13px;
}

.summary-item.total .value {
    color: #409eff;
    font-size: 14px;
}

/* 分页 */
.pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 16px;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

/* 图表区域 */
.charts-section {
    animation: slideDown 0.4s ease;
}

.charts-card {
    border-radius: 12px;
    min-height: 300px;
}

.charts-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.charts-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 700;
}

.chart-placeholder {
    height: 240px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f5f7fa, #e4e7ed);
    border-radius: 8px;
    border: 2px dashed #c0c4cc;
}

.placeholder-content {
    text-align: center;
    color: #909399;
}

.placeholder-content p {
    margin: 8px 0 0;
    font-size: 14px;
}

.placeholder-content .hint {
    font-size: 12px;
    opacity: 0.7;
}

/* 详情对话框 */
.detail-content {
    padding: 4px;
}

.detail-tables {
    margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
    .header-content {
        flex-direction: column;
        align-items: flex-start;
    }

    .filter-container {
        flex-direction: column;
        align-items: stretch;
    }

    .filter-left {
        max-width: 100%;
    }

    .search-input {
        max-width: 100%;
    }
}

@media (max-width: 768px) {
    .dashboard-header {
        padding: 16px 20px;
    }

    .header-content h1 {
        font-size: 22px;
    }

    .dashboard-main {
        padding: 0 16px;
    }

    .stats-grid {
        grid-template-columns: 1fr;
    }

    .summary-row {
        grid-template-columns: 1fr;
    }

    .table-header {
        flex-direction: column;
        gap: 8px;
        align-items: flex-start;
    }

    .header-actions {
        width: 100%;
        justify-content: stretch;
    }

    .header-actions .el-button {
        flex: 1;
    }
}

/* 动画 */
@keyframes slideDown {
    from {
        opacity: 0;
        transform: translateY(-10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* 滚动条美化 */
:deep(.el-table__body-wrapper) {
    &::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }

    &::-webkit-scrollbar-track {
        background: #f1f1f1;
        border-radius: 4px;
    }

    &::-webkit-scrollbar-thumb {
        background: #c0c0c0;
        border-radius: 4px;
        transition: background 0.3s;

        &:hover {
            background: #a8a8a8;
        }
    }
}

/* Element Plus 组件样式覆盖 */
:deep(.el-button) {
    border-radius: 6px;
    font-weight: 600;
}

:deep(.el-card) {
    border-radius: 12px;
    border: none;
}

:deep(.el-table) {
    --el-table-border-color: rgba(0, 0, 0, 0.05);
    --el-table-header-background-color: #f5f7fa;
    --el-table-row-hover-background-color: #f5f7fa;
}

:deep(.el-table__expanded-cell) {
    background: transparent;
    padding: 0;
}

:deep(.el-tag) {
    border-radius: 6px;
    font-weight: 600;
}

:deep(.el-input__wrapper) {
    border-radius: 6px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

:deep(.el-select__wrapper) {
    border-radius: 6px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

:deep(.el-switch__core) {
    border-radius: 12px;
}

:deep(.el-pagination) {
    --el-pagination-button-radius: 6px;
    --el-pagination-button-bg-color: #f4f4f5;
    --el-pagination-button-disabled-bg-color: #e8e8e8;
}

/* 工具类 */
.text-right {
    text-align: right;
}

.font-bold {
    font-weight: 700;
}

.text-success {
    color: #67c23a;
}

.text-danger {
    color: #f56c6c;
}

.text-warning {
    color: #e6a23c;
}
</style>
