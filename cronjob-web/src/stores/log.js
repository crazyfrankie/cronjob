import { defineStore } from 'pinia'
import { getJobLogs } from '@/api/job'
import { ElMessage } from 'element-plus'

export const useLogStore = defineStore('log', {
  state: () => ({
    logs: [],
    loading: false,
    currentLog: null,
    pagination: {
            limit: 20,
      total: 0
    },
    filterJobName: ''
  }),
  
  actions: {
    // 获取任务日志
    async fetchLogs(jobName = '', skip = 0, limit = 20) {
      this.loading = true
      this.filterJobName = jobName
      
      try {
        const response = await getJobLogs(jobName, skip, limit)
        this.logs = response.data || []
        this.pagination.skip = skip
        this.pagination.limit = limit
        
        // 由于后端API没有提供总数，这里做一个估计
        if (this.logs.length < limit) {
          this.pagination.total = skip + this.logs.length
        } else {
          this.pagination.total = skip + limit + 20 // 预估还有更多
        }
        
        return this.logs
      } catch (error) {
        console.error('获取日志失败:', error)
        ElMessage.error('获取日志失败')
        return []
      } finally {
        this.loading = false
      }
    },
    
    // 设置当前日志
    setCurrentLog(log) {
      this.currentLog = log ? { ...log } : null
    },
    
    // 重置当前日志
    resetCurrentLog() {
      this.currentLog = null
    },
    
    // 重置分页
    resetPagination() {
      this.pagination.skip = 0
      this.pagination.limit = 20
    }
  }
})
