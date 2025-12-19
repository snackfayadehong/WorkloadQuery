<template>
    <el-config-provider :locale="zhCn">
        <router-view />
    </el-config-provider>
</template>

<script setup>
import { zhCn } from 'element-plus/es/locale/index' // 引入中文语言包

</script>

<style>
/* 全局样式可以放在这里 */
</style>