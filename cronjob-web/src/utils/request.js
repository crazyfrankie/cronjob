import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API,
  timeout: 15000,
  // 解决跨域请求问题
  withCredentials: true,
  // 确保正确处理不同类型的响应
  responseType: 'json'
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 可以在这里添加认证信息
    // 添加统一的请求头
    config.headers = {
      ...config.headers,
      'Content-Type': 'application/json',
      'X-Requested-With': 'XMLHttpRequest'
    }
    return config
  },
  error => {
    console.error('请求错误：', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    // 如果响应状态码不是200，也直接返回，不作处理
    if (response.status !== 200) {
      return response;
    }
    
    const res = response.data;
    
    // 如果没有code字段，或者直接是数组数据，直接返回数据
    if (res.code === undefined || Array.isArray(res)) {
      return { code: 0, data: res };
    }
    
    // 根据API返回的状态码判断请求是否成功
    if (res.code !== 20000) {
      // 避免重复显示错误消息
      const errMsg = res.msg || '请求失败';
      
      // 错误加上handled标记，避免重复处理
      const error = new Error(errMsg);
      error.handled = true;
      
      ElMessage({
        message: errMsg,
        type: 'error',
        duration: 5000
      });
      
      return Promise.reject(error);
    } else {
      return res;
    }
  },
  error => {
    console.error('响应错误：', error);
    
    // 如果请求被取消，不显示错误信息
    if (axios.isCancel(error)) {
      return Promise.reject(error);
    }

    // 网络错误特殊处理
    if (error.message && error.message.includes('Network Error')) {
      ElMessage({
        message: '网络连接错误，请检查服务器是否运行',
        type: 'error',
        duration: 5000
      })
    } else {
      ElMessage({
        message: error.message || '请求失败，请检查网络连接',
        type: 'error',
        duration: 5000
      })
    }
    
    return Promise.reject(error)
  }
)

export default service
