import axios from 'axios';
import store from '@/store';
import { getToken } from '@/utils/auth';
import { ElMessageBox } from 'element-plus';

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
});

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent

    if (store.user().token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['Authorization'] = 'Bearer ' + getToken()
    }
    return config;
  },
  error => {
    // do something with request error
    console.log(error); // for debug
    return Promise.reject(error);
  }
);

// response interceptor
service.interceptors.response.use(
  response => {
    const res = response.data;
    
    // 直接返回响应数据，不进行成功与否的判断
    return res;
  },
  error => {
    console.log(error);
    
    if (error.response) {
      // 根据HTTP状态码进行不同的错误处理
      switch (error.response.status) {
        case 401:
          if (error.response.data.Message.indexOf('JWT认证失败') !== -1) {
            // 弹出确认框，让用户选择重新登录或留在当前页面
            ElMessageBox.confirm(
              '登录超时, 重新登录或继续停留在当前页？',
              '登录状态已失效',
              {
                confirmButtonText: '重新登录',
                cancelButtonText: '继续停留',
                type: 'warning'
              }
            ).then(async () => {
              // 重置令牌并刷新页面
             await store.user().resetToken()
             location.reload();
            });
          } else {
            // 显示错误消息
            ElMessage({
              message: error.response.data.Message,
              type: 'error',
              duration: 5 * 1000
            });
          }
          break;
        case 403:
          // 重定向到403错误页面
          router.push({ path: '/401' });
          break;
        default:
          // 显示通用错误消息
          ElMessage({
            message: error.response.data.Message || error.message,
            type: 'error',
            duration: 5 * 1000
          });
      }
    } else {
      // 处理非HTTP错误（如网络错误）
      ElMessage({
        message: error.message,
        type: 'error',
        duration: 5 * 1000
      });
    }
    // 返回Promise.reject以拒绝错误的响应
    return Promise.reject(error);
  }
);

export default service