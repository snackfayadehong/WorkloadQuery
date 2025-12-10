<template>
    <el-table :data="tableData" :border="parentBorder" :preserve-expanded-content="preserveExpanded" style="width: 100%">
        <el-table-column type="expand">
            <template #default="props">
                <div m="4">
                    <el-table :data="format(props.row)" style="width: 100%" border>
                        <!-- 入库信息 -->
                        <el-table-column label="入库信息" align="center" header-align="center">
                            <el-table-column prop="inCategory" label="类别" width="100" align="center" />
                            <el-table-column prop="inSpecCount" label="品规数" width="100" align="center" />
                            <el-table-column prop="inBillCount" label="单据数" width="100" align="center" />
                            <el-table-column prop="inTotalAmount" label="总金额" width="120" align="center" />
                        </el-table-column>

                        <!-- 出库信息 -->
                        <el-table-column label="出库信息" align="center" header-align="center">
                            <el-table-column prop="outCategory" label="类别" width="100" align="center" />
                            <el-table-column prop="outSpecCount" label="品规数" width="100" align="center" />
                            <el-table-column prop="outBillCount" label="单据数" width="100" align="center" />
                            <el-table-column prop="outTotalAmount" label="总金额" width="120" align="center" />
                        </el-table-column>

                        <!-- 退还信息 -->
                        <el-table-column label="退还信息" align="center" header-align="center">
                            <el-table-column prop="returnCategory" label="类别" width="100" align="center" />
                            <el-table-column prop="returnSpecCount" label="品规数" width="100" align="center" />
                            <el-table-column prop="returnBillCount" label="单据数" width="100" align="center" />
                            <el-table-column prop="returnTotalAmount" label="总金额" width="120" align="center" />
                        </el-table-column>
                    </el-table>
                </div>
            </template>
        </el-table-column>
        <el-table-column label="操作人员" prop="operator" width="150" align="center" />
    </el-table>
</template>

<script setup>
import { ref } from "vue";

const parentBorder = ref(false);
const childBorder = ref(false);
const preserveExpanded = ref(false);
const tableData = [
    {
        operator: "陈慧灵",
        inbound: [
            {
                category: "其他材料库",
                specCount: 210,
                billCount: 43,
                totalAmount: 143010.69
            },
            {
                category: "试剂与化验材料库",
                specCount: 343,
                billCount: 175,
                totalAmount: 8589789.57
            },
            {
                category: "消毒供应库",
                specCount: 88,
                billCount: 34,
                totalAmount: 791834
            },
            {
                category: "卫生材料",
                specCount: 506,
                billCount: 337,
                totalAmount: 3312207.21
            },
            {
                category: "扫码-高值",
                specCount: 162,
                billCount: 40,
                totalAmount: 1254882.25
            }
        ],
        outbound: [
            {
                category: "扫码-高值",
                specCount: 87,
                billCount: 24,
                totalAmount: 1087430.11
            }
        ],
        return: [
            {
                category: "卫生材料",
                specCount: 2,
                billCount: 2,
                totalAmount: 32170
            },
            {
                category: "扫码-高值",
                specCount: 1,
                billCount: 1,
                totalAmount: 846
            }
        ]
    },
    {
        operator: "李宇烟",
        inbound: [
            {
                category: "试剂与化验材料库",
                specCount: 12,
                billCount: 5,
                totalAmount: 97280.2
            },
            {
                category: "扫码-高值",
                specCount: 78,
                billCount: 6,
                totalAmount: 139095
            },
            {
                category: "其他材料库",
                specCount: 5,
                billCount: 1,
                totalAmount: 1459.5
            },
            {
                category: "卫生材料",
                specCount: 503,
                billCount: 334,
                totalAmount: 3302036.01
            },
            {
                category: "扫码-低值",
                specCount: 58,
                billCount: 45,
                totalAmount: 717844
            }
        ],
        outbound: [
            {
                category: "试剂与化验材料库",
                specCount: 20,
                billCount: 13,
                totalAmount: 101733.94
            },
            {
                category: "扫码-高值",
                specCount: 3,
                billCount: 2,
                totalAmount: 74409
            },
            {
                category: "扫码-低值",
                specCount: 77,
                billCount: 55,
                totalAmount: 796513.92
            },
            {
                category: "其他材料库",
                specCount: 1,
                billCount: 1,
                totalAmount: 129.5
            },
            {
                category: "卫生材料",
                specCount: 2283,
                billCount: 902,
                totalAmount: 3181031.31
            }
        ],
        return: [
            {
                category: "卫生材料",
                specCount: 20,
                billCount: 17,
                totalAmount: 63589
            }
        ]
    },
    {
        operator: "廖小凤",
        inbound: [
            {
                category: "卫生材料",
                specCount: 3,
                billCount: 3,
                totalAmount: 10171.2
            },
            {
                category: "试剂与化验材料库",
                specCount: 285,
                billCount: 137,
                totalAmount: 7225978.57
            }
        ],
        outbound: [
            {
                category: "试剂与化验材料库",
                specCount: 491,
                billCount: 105,
                totalAmount: 6567878.76
            },
            {
                category: "卫生材料",
                specCount: 3,
                billCount: 3,
                totalAmount: 10171.2
            }
        ],
        return: [
            {
                category: "试剂与化验材料库",
                specCount: 3,
                billCount: 1,
                totalAmount: 6300
            }
        ]
    },
    {
        operator: "王健",
        inbound: [
            {
                category: "扫码-高值",
                specCount: 353,
                billCount: 176,
                totalAmount: 5427444.64
            }
        ],
        outbound: [
            {
                category: "扫码-高值",
                specCount: 337,
                billCount: 100,
                totalAmount: 5137919.92
            }
        ],
        return: [
            {
                category: "扫码-高值",
                specCount: 35,
                billCount: 9,
                totalAmount: 1776278
            }
        ]
    },
    {
        operator: "杨有恒",
        inbound: [
            {
                category: "扫码-高值",
                specCount: 2032,
                billCount: 761,
                totalAmount: 15782608.6
            }
        ],
        outbound: [
            {
                category: "扫码-高值",
                specCount: 1115,
                billCount: 242,
                totalAmount: 14952241.5
            }
        ],
        return: [
            {
                category: "扫码-高值",
                specCount: 256,
                billCount: 44,
                totalAmount: 6156116.2
            }
        ]
    },
    {
        operator: "赵洁",
        inbound: [
            {
                category: "试剂与化验材料库",
                specCount: 46,
                billCount: 33,
                totalAmount: 1266530.8
            },
            {
                category: "其他材料库",
                specCount: 205,
                billCount: 42,
                totalAmount: 141551.19
            }
        ],
        outbound: [
            {
                category: "试剂与化验材料库",
                specCount: 30,
                billCount: 15,
                totalAmount: 1123070.5
            },
            {
                category: "卫生材料",
                specCount: 35,
                billCount: 6,
                totalAmount: 6575.9
            },
            {
                category: "其他材料库",
                specCount: 120,
                billCount: 76,
                totalAmount: 81604
            }
        ],
        return: [
            {
                category: "试剂与化验材料库",
                specCount: 1,
                billCount: 1,
                totalAmount: 1200
            }
        ]
    },
    {
        operator: "赵培英",
        inbound: [
            {
                category: "消毒供应库",
                specCount: 88,
                billCount: 34,
                totalAmount: 791834
            }
        ],
        outbound: [
            {
                category: "消毒供应库",
                specCount: 1155,
                billCount: 244,
                totalAmount: 391142.95
            }
        ],
        return: [
            {
                category: "消毒供应库",
                specCount: 3,
                billCount: 3,
                totalAmount: 421.2
            }
        ]
    },
    {
        operator: "孙立",
        outbound: [
            {
                category: "消毒供应库",
                specCount: 775,
                billCount: 141,
                totalAmount: 192527.03
            }
        ],
        return: [
            {
                category: "消毒供应库",
                specCount: 2,
                billCount: 2,
                totalAmount: 165
            }
        ]
    }
];
</script>
