<template>
  <div class="page-container dashboard">
    <GopherIntro />
    
    <div class="page-header">
      <h2>系统概览</h2>
      <p class="subtitle">实时监控任务执行状态和系统运行情况</p>
    </div>
    
    <el-row :gutter="20">
      <!-- 系统概览卡片 -->
      <el-col :span="6" v-for="(card, index) in statsCards" :key="index">
        <el-card 
          shadow="hover" 
          class="stats-card custom-card hover-lift animate__animated animate__fadeInUp"
          :style="{'animation-delay': index * 0.1 + 's'}"
        >
          <div class="stats-content">
            <el-icon class="stats-icon" :style="{ color: card.color }">
              <component :is="card.icon"></component>
            </el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ card.value }}</div>
              <div class="stats-label">{{ card.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" class="mt-20">
      <!-- 任务执行状态图表 -->
      <el-col :span="16">
        <el-card shadow="hover" class="custom-card animate__animated animate__fadeInLeft" style="animation-delay: 0.2s">
          <template #header>
            <div class="card-header">
              <span><el-icon><DataAnalysis /></el-icon> 任务执行状态</span>
              <el-button-group>
                <el-button size="small" type="primary" plain class="btn-pulse">今日</el-button>
                <el-button size="small" plain>本周</el-button>
                <el-button size="small" plain>本月</el-button>
              </el-button-group>
            </div>
          </template>
          <div class="chart-container">
            <div class="chart-placeholder animate__animated animate__fadeIn" style="animation-delay: 0.4s">
              <el-icon class="float-animation"><DataAnalysis /></el-icon>
              <span>任务执行趋势图</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 系统状态 -->
      <el-col :span="8">
        <el-card shadow="hover" class="custom-card animate__animated animate__fadeInRight" style="animation-delay: 0.2s">
          <template #header>
            <div class="card-header">
              <span><el-icon><Monitor /></el-icon> 系统状态</span>
              <el-tag size="small" type="success" effect="dark">正常运行</el-tag>
            </div>
          </template>
          <div class="system-stats">
            <div class="stat-item" v-for="(item, index) in systemStats" :key="index">
              <div class="stat-label">
                <el-icon class="stat-icon" v-if="getSystemStatsIcon(item.label)">
                  <component :is="getSystemStatsIcon(item.label)"></component>
                </el-icon>
                {{ item.label }}
              </div>
              <div class="stat-progress">
                <el-progress
                  :percentage="item.value"
                  :stroke-width="10"
                  :color="getProgressColor(item.value)"
                  :format="(val) => `${val}%`"
                  class="animate__animated animate__fadeIn"
                  :style="{'animation-delay': (index * 0.2 + 0.3) + 's'}"
                ></el-progress>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" class="mt-20">
      <!-- 最近任务列表 -->
      <el-col :span="24">
        <el-card shadow="hover" class="custom-card animate__animated animate__fadeInUp" style="animation-delay: 0.3s">
          <template #header>
            <div class="card-header">
              <span><el-icon><Operation /></el-icon> 最近任务执行</span>
              <el-button size="small" type="primary" plain @click="refreshRecentJobs" class="btn-pulse">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>
          </template>
          <el-table 
            :data="recentJobs" 
            stripe 
            style="width: 100%" 
            :row-class-name="tableRowClassName"
            highlight-current-row
          >
            <el-table-column prop="jobName" label="任务名称" min-width="180">
              <template #default="scope">
                <div class="job-name">
                  <el-icon><Calendar /></el-icon>
                  <span>{{ scope.row.jobName }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="command" label="命令" min-width="200" show-overflow-tooltip>
              <template #default="scope">
                <el-tooltip 
                  :content="scope.row.command" 
                  placement="top" 
                  :show-after="800"
                  class="tooltip-trigger"
                >
                  <div class="command-display">
                    <el-icon><Terminal /></el-icon>
                    <code>{{ scope.row.command }}</code>
                  </div>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column prop="startTime" label="开始时间" min-width="160" :formatter="formatTime" />
            <el-table-column prop="endTime" label="结束时间" min-width="160" :formatter="formatTime" />
            <el-table-column label="状态" width="100" align="center">
              <template #default="scope">
                <el-tag
                  :type="scope.row.err ? 'danger' : 'success'"
                  size="small"
                  effect="dark"
                  :class="{'pulse-animation': scope.row.err}"
                >
                  <el-icon v-if="scope.row.err"><CircleCloseFilled /></el-icon>
                  <el-icon v-else><CircleCheckFilled /></el-icon>
                  {{ scope.row.err ? '失败' : '成功' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" align="center">
              <template #default="scope">
                <el-button link type="primary" size="small" @click="viewJobDetail(scope.row)">
                  查看详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 详情对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="任务执行详情"
      width="650px"
    >
      <div v-if="selectedJob">
        <div class="job-detail-item">
          <span class="detail-label">任务名称：</span>
          <span>{{ selectedJob.jobName }}</span>
        </div>
        <div class="job-detail-item">
          <span class="detail-label">执行命令：</span>
          <span>{{ selectedJob.command }}</span>
        </div>
        <div class="job-detail-item">
          <span class="detail-label">开始时间：</span>
          <span>{{ formatTimeStr(selectedJob.startTime) }}</span>
        </div>
        <div class="job-detail-item">
          <span class="detail-label">结束时间：</span>
          <span>{{ formatTimeStr(selectedJob.endTime) }}</span>
        </div>
        <div class="job-detail-item">
          <span class="detail-label">执行耗时：</span>
          <span>{{ calculateDuration(selectedJob) }}</span>
        </div>
        <div class="job-detail-item">
          <span class="detail-label">执行结果：</span>
          <el-tag
            :type="selectedJob.err ? 'danger' : 'success'"
            size="small"
            effect="dark"
          >
            {{ selectedJob.err ? '失败' : '成功' }}
          </el-tag>
        </div>
        
        <el-divider content-position="left">输出内容</el-divider>
        <div class="job-output">
          <pre>{{ selectedJob.output || '无输出内容' }}</pre>
        </div>
        
        <div v-if="selectedJob.err" class="job-error">
          <el-divider content-position="left">错误信息</el-divider>
          <pre>{{ selectedJob.err }}</pre>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { getJobLogs, getJobList, getWorkerList } from '@/api/job'
import moment from 'moment'
import GopherIntro from '@/components/GopherIntro.vue'

export default {
  name: 'DashboardView',
  components: { 
    GopherIntro
  },
  
  setup() {
    const recentJobs = ref([])
    const dialogVisible = ref(false)
    const selectedJob = ref(null)
    const workerCount = ref(0)
    const jobCount = ref(0)
    
    // 统计卡片数据
    const statsCards = ref([
      {
        icon: 'Calendar',
        label: '总任务数',
        value: 0,
        color: '#409EFF'
      },
      {
        icon: 'Connection',
        label: 'Worker节点',
        value: 0,
        color: '#67C23A'
      },
      {
        icon: 'Check',
        label: '今日成功',
        value: 0,
        color: '#67C23A'
      },
      {
        icon: 'Warning',
        label: '今日失败',
        value: 0,
        color: '#F56C6C'
      }
    ])
    
    // 系统状态数据
    const systemStats = [
      {
        label: 'CPU使用率',
        value: 35,
        color: '#409EFF'
      },
      {
        label: '内存使用率',
        value: 65,
        color: '#67C23A'
      },
      {
        label: '磁盘使用率',
        value: 42,
        color: '#E6A23C'
      }
    ]
    
    // 获取最近任务执行列表
    const fetchRecentJobs = async () => {
      try {
        const res = await getJobLogs('', 0, 10)
        recentJobs.value = res.data || []
        
        // 计算成功失败数
        const today = moment().startOf('day').valueOf()
        const todayLogs = recentJobs.value.filter(job => job.startTime >= today)
        const successCount = todayLogs.filter(job => !job.err).length
        const failedCount = todayLogs.filter(job => job.err).length
        
        statsCards.value[2].value = successCount
        statsCards.value[3].value = failedCount
      } catch (error) {
        console.error('获取最近任务执行列表失败:', error)
      }
    }
    
    // 获取任务数量
    const fetchJobCount = async () => {
      try {
        const res = await getJobList()
        jobCount.value = (res.data || []).length
        statsCards.value[0].value = jobCount.value
      } catch (error) {
        console.error('获取任务列表失败:', error)
      }
    }
    
    // 获取Worker节点数量
    const fetchWorkerCount = async () => {
      try {
        const res = await getWorkerList()
        workerCount.value = (res.data || []).length
        statsCards.value[1].value = workerCount.value
      } catch (error) {
        console.error('获取Worker节点列表失败:', error)
      }
    }
    
    // 刷新最近任务列表
    const refreshRecentJobs = () => {
      fetchRecentJobs()
    }
    
    // 查看任务详情
    const viewJobDetail = (job) => {
      selectedJob.value = job
      dialogVisible.value = true
    }
    
    // 获取系统状态图标
    const getSystemStatsIcon = (label) => {
      if (label.includes('CPU')) {
        return 'CPU'
      } else if (label.includes('内存')) {
        return 'Histogram'
      } else if (label.includes('磁盘')) {
        return 'DataLine'
      }
      return null
    }
    
    // 获取进度条颜色
    const getProgressColor = (percentage) => {
      if (percentage < 50) {
        return '#67C23A' // 绿色
      } else if (percentage < 80) {
        return '#E6A23C' // 黄色
      } else {
        return '#F56C6C' // 红色
      }
    }
    
    // 格式化时间戳为字符串
    const formatTime = (row, column, cellValue) => {
      if (!cellValue) return '-'
      return moment(cellValue).format('YYYY-MM-DD HH:mm:ss')
    }
    
    // 格式化时间戳为字符串 (直接使用)
    const formatTimeStr = (timestamp) => {
      if (!timestamp) return '-'
      return moment(timestamp).format('YYYY-MM-DD HH:mm:ss')
    }
    
    // 计算执行耗时
    const calculateDuration = (job) => {
      if (!job.startTime || !job.endTime) return '-'
      const duration = moment.duration(job.endTime - job.startTime)
      return `${duration.minutes()}分${duration.seconds()}秒`
    }
    
    // 设置表格行样式
    const tableRowClassName = ({ row }) => {
      if (row.err) {
        return 'error-row'
      }
      return ''
    }
    
    onMounted(() => {
      fetchRecentJobs()
      fetchJobCount()
      fetchWorkerCount()
    })
    
    return {
      recentJobs,
      statsCards,
      systemStats,
      dialogVisible,
      selectedJob,
      refreshRecentJobs,
      viewJobDetail,
      formatTime,
      formatTimeStr,
      calculateDuration,
      getSystemStatsIcon,
      getProgressColor,
      tableRowClassName
    }
  }
}
</script>

<style scoped>
.dashboard .subtitle {
  color: #909399;
  margin-top: 0;
}

.mt-20 {
  margin-top: 20px;
}

.stats-card {
  height: 100px;
}

.stats-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.stats-icon {
  font-size: 48px;
  margin-right: 20px;
}

.stats-info {
  display: flex;
  flex-direction: column;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  line-height: 1.2;
}

.stats-label {
  font-size: 14px;
  color: #909399;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chart-container {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #909399;
}

.chart-placeholder .el-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.system-stats {
  padding: 10px 0;
}

.stat-item {
  margin-bottom: 24px;
}

.stat-label {
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
  display: flex;
  align-items: center;
}

.stat-icon {
  margin-right: 6px;
  color: var(--primary-color);
}

.job-detail-item {
  margin-bottom: 12px;
  line-height: 24px;
}

.detail-label {
  font-weight: bold;
  color: #606266;
  margin-right: 8px;
  display: inline-block;
  width: 80px;
  text-align: right;
}

.job-output, .job-error {
  margin-top: 16px;
  margin-bottom: 16px;
}

.job-output pre, .job-error pre {
  background-color: #f9f9f9;
  border: 1px solid #e0e0e0;
  padding: 12px;
  border-radius: 4px;
  max-height: 200px;
  overflow-y: auto;
  line-height: 1.4;
  font-family: monospace;
  white-space: pre-wrap;
  word-break: break-all;
}

.job-error pre {
  color: #f56c6c;
  background-color: #fef0f0;
  border-color: #fde2e2;
}

.job-name {
  display: flex;
  align-items: center;
}

.job-name .el-icon {
  margin-right: 8px;
  color: var(--primary-color);
}

.command-display {
  display: flex;
  align-items: center;
  font-family: Monaco, Menlo, Consolas, "Courier New", monospace;
}

.command-display .el-icon {
  margin-right: 8px;
  color: #606266;
}

.command-display code {
  background-color: rgba(0, 0, 0, 0.05);
  padding: 2px 4px;
  border-radius: 3px;
  font-size: 12px;
}

:deep(.el-table .error-row) {
  background-color: rgba(245, 108, 108, 0.05);
}

:deep(.el-table .el-table__row:hover) {
  background-color: rgba(0, 173, 216, 0.05);
  transition: background-color 0.3s;
}
</style>
