<template>
    <div class="workload-page">
        <!-- <QueryBar /> -->
        <WorkLoadDataShow2></WorkLoadDataShow2>
        <!-- <UnCheckBillsShow /> -->
    </div>
</template>

<script>
// import QueryBar from "@/components/QueryBarTop/QueryBar.vue";
import WorkLoadDataShow2 from "./components/WorkLoadDataShow2.vue";
// import UnCheckBillsShow from "@/components/UnCheckBillsShow.vue";

export default {
    name: "WorkloadPage",
    components: {
        QueryBar,
        WorkLoadDataShow2
        // UnCheckBillsShow
    }
};
</script>

<style scoped>
.workload-page {
    padding: 16px;
}
</style>
