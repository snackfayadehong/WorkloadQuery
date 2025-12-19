import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "@/styles/global.css";
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
// 注意：由于 vite 配置了自动导入，这里不需要手动 use(ElementPlus) 也能使用组件
// 但如果你需要全局配置（如中文），仍可保留
// 1. 导入所有图标
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate)
// 2. 循环注册所有图标，这一步是图标能显示的关键
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(pinia);
app.use(router);
app.mount("#app");