import { defineStore } from 'pinia'
import { getWorkerList } from '@/api/job'
import { ElMessage } from 'element-plus'

export const useWorkerStore = defineStore('worker', {
  state: () => ({
    workers: [],
    loading: false
  }),
  
  actions: {
    // 获取Worker节点列表
    async fetchWorkers() {
      this.loading = true
      try {
        const response = await getWorkerList()
        
        // 检查response和data字段是否存在
        if (response && response.data) {
          this.workers = response.data
        } else {
          // 如果响应格式不正确，使用空数组
          console.warn('Worker列表响应数据格式不正确:', response)
          this.workers = []
        }
        
        return this.workers
      } catch (error) {
        console.error('获取Worker节点列表失败:', error)
        // 避免重复显示错误消息
        // 只有在响应拦截器没有显示错误时才显示
        if (!error.handled) {
          ElMessage.error('获取Worker节点列表失败')
        }
        return []
      } finally {
        this.loading = false
      }
    }
  }
})
