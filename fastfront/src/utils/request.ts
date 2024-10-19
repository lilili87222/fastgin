import axios from "axios";
import type {
  AxiosResponse,
  AxiosInstance,
  InternalAxiosRequestConfig,
} from "axios";
import store from "@/store";
import { getToken } from "@/utils/auth";
import { ElMessageBox, ElMessage } from "element-plus";
import router from "@/router"; // 确保导入router

// 创建一个axios实例
const service: AxiosInstance = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  timeout: 5000, // 请求超时
});

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    // 使用 InternalAxiosRequestConfig
    // 在请求发送之前做一些处理
    if (store.user().token) {
      // 让每个请求都携带token
      config.headers["Authorization"] = "Bearer " + getToken();
    }
    return config;
  },
  (error: any) => {
    // 处理请求错误
    console.log(error); // 调试用
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data;
    // 直接返回响应数据
    return res;
  },
  (error: any) => {
    console.log(error);

    if (error.response) {
      // 根据HTTP状态码进行不同的错误处理
      switch (error.response.status) {
        case 401:
          if (error.response.data.Message.indexOf("JWT认证失败") !== -1) {
            // 弹出确认框，让用户选择重新登录或留在当前页面
            ElMessageBox.confirm(
              "登录超时, 重新登录或继续停留在当前页？",
              "登录状态已失效",
              {
                confirmButtonText: "重新登录",
                cancelButtonText: "继续停留",
                type: "warning",
              }
            ).then(async () => {
              // 重置令牌并刷新页面
              await store.user().resetToken();
              location.reload();
            });
          } else {
            // 显示错误消息
            ElMessage({
              message: error.response.data.Message,
              type: "error",
              duration: 5 * 1000,
            });
          }
          break;
        case 403:
          // 重定向到403错误页面
          router.push({ path: "/401" });
          break;
        default:
          // 显示通用错误消息
          ElMessage({
            message: error.response.data.Message || error.message,
            type: "error",
            duration: 5 * 1000,
          });
      }
    } else {
      // 处理非HTTP错误（如网络错误）
      ElMessage({
        message: error.message,
        type: "error",
        duration: 5 * 1000,
      });
    }
    // 返回Promise.reject以拒绝错误的响应
    return Promise.reject(error);
  }
);

export default service;
