<template>
  <div class="page-container logs">
    <div class="page-header flex-between animate__animated animate__fadeIn">
      <div>
        <h2><el-icon><Document /></el-icon> 日志查看</h2>
        <p class="subtitle">查看所有任务的执行日志和执行状态</p>
      </div>
      <div>
        <el-button type="primary" @click="fetchLogData" class="btn-pulse">
          <el-icon><Refresh /></el-icon> 刷新数据
        </el-button>
      </div>
    </div>
    
    <!-- 搜索和过滤 -->
    <div class="filter-bar animate__animated animate__fadeInDown" style="animation-delay: 0.2s">
      <el-form :inline="true" class="filter-form">
        <el-form-item label="任务名称">
          <el-input
            v-model="jobNameFilter"
            placeholder="请输入任务名称"
            clearable
            prefix-icon="Search"
            @clear="onJobNameFilterClear"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onFilterSubmit">
            <el-icon><Search /></el-icon> 查询
          </el-button>
          <el-button @click="resetFilter">
            <el-icon><RefreshRight /></el-icon> 重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 日志表格 -->
    <el-card shadow="hover" class="log-table-card custom-card animate__animated animate__fadeIn" style="animation-delay: 0.3s">
      <template #header>
        <div class="card-header">
          <span><el-icon><Document /></el-icon> 日志列表</span>
          <div>
            <el-tag type="success" v-if="logs.length > 0">共 {{ logs.length }} 条记录</el-tag>
          </div>
        </div>
      </template>
      <el-table
        v-loading="loading"
        :data="logs"
        stripe
        style="width: 100%"
        :empty-text="loading ? '加载中...' : '暂无日志数据'"
        :row-class-name="tableRowClassName"
        highlight-current-row
        border
        class="log-table"
      >
        <el-table-column prop="jobName" label="任务名称" min-width="180" show-overflow-tooltip>
          <template #header>
            <span class="column-header"><el-icon><Calendar /></el-icon> 任务名称</span>
          </template>
          <template #default="scope">
            <div class="job-name">
              <el-icon><Calendar /></el-icon>
              <span>{{ scope.row.jobName }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="command" label="命令" min-width="220" show-overflow-tooltip>
          <template #header>
            <span class="column-header"><el-icon><Terminal /></el-icon> 命令</span>
          </template>
          <template #default="scope">
            <div class="command-display">
              <code>{{ scope.row.command }}</code>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="开始时间" min-width="180">
          <template #header>
            <span class="column-header"><el-icon><Timer /></el-icon> 开始时间</span>
          </template>
          <template #default="scope">
            <span class="time-display">{{ formatDate(scope.row.startTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="结束时间" min-width="180">
          <template #header>
            <span class="column-header"><el-icon><Timer /></el-icon> 结束时间</span>
          </template>
          <template #default="scope">
            <span class="time-display">{{ formatDate(scope.row.endTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="执行耗时" width="140">
          <template #header>
            <span class="column-header"><el-icon><Clock /></el-icon> 执行耗时</span>
          </template>
          <template #default="scope">
            {{ calculateDuration(scope.row) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag
              :type="scope.row.err ? 'danger' : 'success'"
              effect="dark"
              size="small"
            >
              {{ scope.row.err ? '失败' : '成功' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="viewLogDetail(scope.row)">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="logDetailVisible"
      title="任务执行详情"
      width="650px"
    >
      <div v-if="currentLog">
        <div class="log-detail-item">
          <span class="detail-label">任务名称：</span>
          <span>{{ currentLog.jobName }}</span>
        </div>
        <div class="log-detail-item">
          <span class="detail-label">执行命令：</span>
          <span>{{ currentLog.command }}</span>
        </div>
        <div class="log-detail-item">
          <span class="detail-label">开始时间：</span>
          <span>{{ formatDate(currentLog.startTime) }}</span>
        </div>
        <div class="log-detail-item">
          <span class="detail-label">结束时间：</span>
          <span>{{ formatDate(currentLog.endTime) }}</span>
        </div>
        <div class="log-detail-item">
          <span class="detail-label">执行耗时：</span>
          <span>{{ calculateDuration(currentLog) }}</span>
        </div>
        <div class="log-detail-item">
          <span class="detail-label">执行结果：</span>
          <el-tag
            :type="currentLog.err ? 'danger' : 'success'"
            size="small"
            effect="dark"
          >
            {{ currentLog.err ? '失败' : '成功' }}
          </el-tag>
        </div>
        
        <el-divider content-position="left">输出内容</el-divider>
        <div class="log-output">
          <pre>{{ currentLog.output || '无输出内容' }}</pre>
        </div>
        
        <div v-if="currentLog.err" class="log-error">
          <el-divider content-position="left">错误信息</el-divider>
          <pre>{{ currentLog.err }}</pre>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="logDetailVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useLogStore } from '@/stores/log'
import { storeToRefs } from 'pinia'
import moment from 'moment'

export default {
  name: 'LogsView',
  
  setup() {
    const logStore = useLogStore()
    const { logs, loading } = storeToRefs(logStore)
    
    const logDetailVisible = ref(false)
    const jobNameFilter = ref('')
    const currentPage = ref(1)
    const pageSize = ref(20)
    
    // 计算总数
    const total = computed(() => logStore.pagination.total)
    
    // 当前选中的日志
    const currentLog = computed(() => logStore.currentLog)
    
    // 查询日志数据
    const fetchLogData = async (reset = true) => {
      if (reset) {
        currentPage.value = 1
      }
      
      const skip = (currentPage.value - 1) * pageSize.value
      await logStore.fetchLogs(jobNameFilter.value, skip, pageSize.value)
    }
    
    // 重置筛选条件
    const resetFilter = () => {
      jobNameFilter.value = ''
      currentPage.value = 1
      logStore.resetPagination()
      fetchLogData()
    }
    
    // 清除任务名称筛选
    const onJobNameFilterClear = () => {
      fetchLogData()
    }
    
    // 提交筛选
    const onFilterSubmit = () => {
      fetchLogData()
    }
    
    // 查看日志详情
    const viewLogDetail = (log) => {
      logStore.setCurrentLog(log)
      logDetailVisible.value = true
    }
    
    // 处理页码变更
    const handleCurrentChange = (page) => {
      currentPage.value = page
      fetchLogData(false)
    }
    
    // 处理每页条数变更
    const handleSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1
      fetchLogData(false)
    }
    
    // 格式化日期（北京时间）
    const formatDate = (timestamp) => {
      if (!timestamp) return '-'
      // 兼容后端返回的Unix秒
      if (timestamp > 1e12) {
        // 已经是毫秒
        return moment(timestamp).utcOffset(8).format('YYYY-MM-DD HH:mm:ss')
      } else {
        // 秒
        return moment.unix(timestamp).utcOffset(8).format('YYYY-MM-DD HH:mm:ss')
      }
    }
    
    // 计算执行耗时
    const calculateDuration = (log) => {
      if (!log || !log.startTime || !log.endTime) return '-'
      const duration = moment.duration(log.endTime - log.startTime)
      const minutes = duration.minutes()
      const seconds = duration.seconds()
      
      if (minutes > 0) {
        return `${minutes}分${seconds}秒`
      } else {
        return `${seconds}秒`
      }
    }
    
    // 监听页面大小变化，更新每页条数
    watch(() => pageSize.value, (newVal) => {
      logStore.pagination.limit = newVal
    })
    
    onMounted(() => {
      fetchLogData()
    })
    
    return {
      logs,
      loading,
      jobNameFilter,
      currentPage,
      pageSize,
      total,
      currentLog,
      logDetailVisible,
      fetchLogData,
      resetFilter,
      onJobNameFilterClear,
      onFilterSubmit,
      viewLogDetail,
      handleCurrentChange,
      handleSizeChange,
      formatDate,
      calculateDuration
    }
  }
}
</script>

<style scoped>
.logs .subtitle {
  color: #909399;
  margin-top: 0;
}

.filter-bar {
  margin-bottom: 20px;
}

.filter-form {
  background-color: #fff;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.log-detail-item {
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

.log-output, .log-error {
  margin-top: 16px;
  margin-bottom: 16px;
}

.log-output pre, .log-error pre {
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

.log-error pre {
  color: #f56c6c;
  background-color: #fef0f0;
  border-color: #fde2e2;
}
</style>
