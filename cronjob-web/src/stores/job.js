import { defineStore } from 'pinia'
import { getJobList, saveJob, deleteJob, killJob } from '@/api/job'
import { ElMessage } from 'element-plus'

export const useJobStore = defineStore('job', {
  state: () => ({
    jobs: [],
    loading: false,
    currentJob: null
  }),
  
  actions: {
    // 获取任务列表
    async fetchJobs() {
      this.loading = true
      try {
        const response = await getJobList()
        this.jobs = response.data || []
      } catch (error) {
        console.error('获取任务列表失败:', error)
        ElMessage.error('获取任务列表失败')
      } finally {
        this.loading = false
      }
    },
    
    // 保存或更新任务
    async saveOrUpdateJob(jobData) {
      try {
        await saveJob(jobData)
        ElMessage.success('保存任务成功')
        this.fetchJobs()
        return true
      } catch (error) {
        console.error('保存任务失败:', error)
        ElMessage.error('保存任务失败')
        return false
      }
    },
    
    // 删除任务
    async removeJob(name) {
      try {
        await deleteJob(name)
        ElMessage.success('删除任务成功')
        this.fetchJobs()
        return true
      } catch (error) {
        console.error('删除任务失败:', error)
        ElMessage.error('删除任务失败')
        return false
      }
    },
    
    // 终止任务
    async stopJob(name) {
      try {
        await killJob(name)
        ElMessage.success('终止任务成功')
        return true
      } catch (error) {
        console.error('终止任务失败:', error)
        ElMessage.error('终止任务失败')
        return false
      }
    },
    
    // 设置当前编辑的任务
    setCurrentJob(job) {
      this.currentJob = job ? { ...job } : null
    },
    
    // 重置当前任务
    resetCurrentJob() {
      this.currentJob = null
    }
  }
})
