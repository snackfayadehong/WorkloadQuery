// src/vite.config.js 优化版
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";
import VueSetupExtend from "vite-plugin-vue-setup-extend";
import { resolve } from "path";

export default defineConfig({
    envDir: "env",
    plugins: [
        vue(),
        VueSetupExtend(), // 支持 <script setup name="MyComponent">
        AutoImport({
            // 自动导入 Vue, Vue Router, Pinia 的 API
            imports: ["vue", "vue-router", "pinia"],
            resolvers: [ElementPlusResolver()],
            dts: "src/auto-imports.d.ts" // 生成类型声明文件
        }),
        Components({
            // 自动导入 src/components 下的组件以及 Element Plus 组件
            resolvers: [ElementPlusResolver()],
            dts: "src/components.d.ts"
        })
    ],
    server: {
        // host: "0.0.0.0", // 允许局域网访问
        host: "127.0.0.1",
        port: 5173, // 固定端口
        open: false // 启动后自动打开浏览器
    },
    build: {
        outDir: "dist",
        reportCompressedSize: false, // 禁用大小报告，提高构建速度
        rollupOptions: {
            output: {
                // 优化：更加细致的代码分割
                manualChunks(id) {
                    if (id.includes("node_modules")) {
                        // 将每个库打包成独立文件
                        return id.toString().split("node_modules/")[1].split("/")[0].toString();
                    }
                }
            }
        }
    },
    resolve: {
        alias: {
            "@": resolve(__dirname, "src") //
        }
    }
});
