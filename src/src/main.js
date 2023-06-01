import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import ElementPlus from "element-plus";
import zhCn from "element-plus/dist/locale/zh-cn.mjs"; // 组件中文化

const app = createApp(App);

// 中文
app.use(ElementPlus, {
    locale: zhCn
});

app.mount("#app");

// createApp(App).mount("#app");
