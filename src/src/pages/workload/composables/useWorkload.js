import { ref, computed, onMounted, watch } from "vue";
import myAxios from "@/services/myAxios";
import { ElMessage } from "element-plus";
import dayjs from "dayjs";

export function useWorkload() {
    // --- 1. 数据状态 ---
    const rawData = ref([]); // 存储从后端获取的原始聚合数据
    const loading = ref(false); // 全局加载状态控制

    // --- 2. 筛选条件 (与页面 UI 绑定) ---
    const searchQuery = ref(""); // 搜索操作员名称
    const filterType = ref("all"); // 业务类型筛选: all | inbound | outbound | return
    const dateRange = ref([]); // 时间范围筛选

    // --- 3. 分页配置 ---
    const currentPage = ref(1); // 当前页码
    const pageSize = ref(10); // 每页条数

    /**
     * 获取数据：调用后端 API
     */
    const fetchList = async () => {
        if (!dateRange.value || dateRange.value.length < 2) {
            ElMessage.error("请选择完整时间范围!")
            return
        }
        loading.value = true;
        try {
            // 注意：这里对接后端聚合接口
            const res = await myAxios.post("/getWorkload", {
                startTime: dateRange.value?.[0] ? dayjs(dateRange.value[0]).startOf('day').format("YYYY-MM-DD HH:mm:ss") : "",
                endTime: dateRange.value?.[1] ? dayjs(dateRange.value[1]).endOf('day').format("YYYY-MM-DD HH:mm:ss") : ""
            });

            // 这里的 res 已经是 myAxios 拦截器处理后的 res.data
            if (res.code === 0) {
                rawData.value = res.data || [];
            } else {
                ElMessage.error(res.message || "获取数据失败");
            }
        } catch (error) {
            console.error("请求异常:", error);
        } finally {
            loading.value = false;
        }
    };

    /**
     * 计算属性：根据 姓名 和 业务类型 进行实时前端过滤
     */
    const filteredData = computed(() => {
        return rawData.value.filter(item => {
            // 1. 姓名匹配
            const matchName = item.operator.toLowerCase().includes(searchQuery.value.toLowerCase());

            // 2. 业务类型匹配
            // 如果筛选了具体类型（入库/出库/退还），则只显示该项数据不为空的人员
            let matchType = true;
            if (filterType.value !== "all") {
                const typeData = item[filterType.value];
                matchType = Array.isArray(typeData) && typeData.length > 0;
            }

            return matchName && matchType;
        });
    });

    /**
     * 计算属性：对过滤后的结果进行前端切片（分页）
     */
    const paginatedData = computed(() => {
        const start = (currentPage.value - 1) * pageSize.value;
        const end = start + pageSize.value;
        return filteredData.value.slice(start, end);
    });

    /**
     * 监听器：当搜索条件或筛选类型变化时，重置页码到第一页
     */
    watch([searchQuery, filterType], () => {
        currentPage.value = 1;
    });

    /**
     * 监听器：当日期变化时，重新从后端拉取原始数据
     */
    watch(dateRange, (newVal) => {
        fetchList();
    });

    // 初始化加载
    onMounted(() => {
        // fetchList();
    });

    return {
        // 原始状态
        rawData,
        loading,
        searchQuery,
        filterType,
        dateRange,
        currentPage,
        pageSize,
        // 处理后的数据
        filteredData,
        paginatedData,
        // 核心方法
        fetchList
    };
}
