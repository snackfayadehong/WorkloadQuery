import axios from "axios";
import { ElLoading } from "element-plus";
import "element-plus/theme-chalk/el-loading.css";

// 自定义axios实例
const myAxios = axios.create({
    // 环境切换
    baseURL: import.meta.env.MODE === "production" ? import.meta.env.VITE_HTTP : import.meta.env.VITE_BASE_URL,
    // baseURL: "http://172.21.1.75:3007/api",
    // baseURL: "http://127.0.0.1:3007/api",
    timeout: 10000,
    withCredentialsL: false //跨域请求需要凭证
});

myAxios.defaults.headers.post["Content-Type"] = "application/json;charset=UTF-8";

let loadingInstance = null;
// 请求拦截器
myAxios.interceptors.request.use(
    req => {
        loadingInstance = ElLoading.service({ fullscreen: true });
        return req;
    },
    err => {
        loadingInstance.close();
        return Promise.reject(err);
    }
);
// 响应拦截器
myAxios.interceptors.response.use(
    res => {
        if (res.status !== 200) {
            loadingInstance.close();
            return "";
        }
        loadingInstance.close();
        return res.data;
    },
    err => {
        loadingInstance.close();
        if (err.message.includes("timeout")) {
            alert("网络异常,请重试！");
        }
        return Promise.reject(err);
    }
);

export default myAxios;
