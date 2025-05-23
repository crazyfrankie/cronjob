<template>
  <div class="page-container workers">
    <div class="page-header flex-between">
      <div>
        <h2>节点管理</h2>
        <p class="subtitle">查看和管理所有Worker节点状态</p>
      </div>
      <div>
        <el-button type="primary" @click="refreshWorkers" class="custom-button">
          <el-icon><Refresh /></el-icon> 刷新节点
        </el-button>
      </div>
    </div>
    
    <!-- 节点数据概览 -->
    <el-row :gutter="20" class="metrics-row animate__animated animate__fadeIn">
      <el-col :span="6">
        <div class="metric-card" :style="{ backgroundColor: 'rgba(0, 173, 216, 0.05)' }">
          <div class="metric-value animate__animated animate__fadeInUp">
            {{ workers.length }}
          </div>
          <div class="metric-name animate__animated animate__fadeInUp animate__delay-1s">总节点数</div>
          <el-icon class="background-icon"><Connection /></el-icon>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card" :style="{ backgroundColor: 'rgba(103, 194, 58, 0.05)' }">
          <div class="metric-value animate__animated animate__fadeInUp">
            {{ workers.length }}
          </div>
          <div class="metric-name animate__animated animate__fadeInUp animate__delay-1s">在线节点</div>
          <el-icon class="background-icon"><CircleCheckFilled /></el-icon>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="metric-card" :style="{ backgroundColor: 'rgba(255, 203, 45, 0.05)' }">
          <div class="metric-value animate__animated animate__fadeInUp">
            {{ getLastUpdateTimeText() }}
          </div>
          <div class="metric-name animate__animated animate__fadeInUp animate__delay-1s">最近更新时间</div>
          <el-icon class="background-icon"><Timer /></el-icon>
        </div>
      </el-col>
    </el-row>
    
    <h3 class="section-title animate__animated animate__fadeIn"><el-icon><Monitor /></el-icon> 节点列表</h3>
    
    <!-- 节点卡片列表 -->
    <el-row :gutter="20">
      <el-col :span="8" v-for="(worker, index) in workers" :key="index">
        <el-card 
          shadow="hover" 
          class="worker-card custom-card" 
          :body-style="{ padding: '0px' }"
          :class="{'animate__animated animate__fadeInUp': true}"
          :style="{'animation-delay': index * 0.1 + 's'}"
        >
          <div class="worker-header">
            <div class="worker-info">
              <el-avatar :size="40" class="worker-avatar">
                <el-icon><Monitor /></el-icon>
              </el-avatar>
              <div class="worker-name-container">
                <span class="worker-name">{{ worker || 'Unknown' }}</span>
                <el-tag 
                  size="small" 
                  type="success"
                >
                  <el-icon>
                    <CircleCheck />
                  </el-icon>
                  在线
                </el-tag>
              </div>
            </div>
          </div>
          <div class="worker-body">
            <div class="info-row">
              <span class="info-label"><el-icon><Location /></el-icon> IP地址:</span>
              <span class="info-value">{{ worker || 'N/A' }}</span>
            </div>
          </div>
          <div class="worker-footer">
            <el-button 
              type="primary" 
              @click="showWorkerDetail(worker)"
            >
              <el-icon><View /></el-icon> 查看节点详情
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <!-- 无节点时的提示 -->
      <el-col :span="24" v-if="workers.length === 0 && !loading">
        <el-empty 
          description="暂无节点数据" 
          :image-size="200"
        >
          <template #description>
            <p>当前没有活跃的Worker节点，请检查节点服务是否正常运行</p>
          </template>
          <el-button type="primary" @click="refreshWorkers">刷新节点</el-button>
        </el-empty>
      </el-col>
      
      <!-- 加载中 -->
      <el-col :span="24" v-if="loading">
        <div class="loading-container">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <span>正在加载节点数据...</span>
        </div>
      </el-col>
    </el-row>
    
    <!-- 节点详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="节点详细信息"
      direction="rtl"
      size="400px"
      class="worker-detail-drawer"
    >
      <div v-if="currentWorker">
        <div class="worker-detail-header animate__animated animate__fadeInRight">
          <el-avatar :size="60" class="worker-detail-avatar">
            <el-icon><Monitor /></el-icon>
          </el-avatar>
          <div class="worker-detail-info">
            <h3 class="worker-detail-name">{{ currentWorker || 'Unknown' }}</h3>
            <el-tag 
              type="success"
            >
              <el-icon><CircleCheck /></el-icon>
              在线
            </el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="worker-detail-content">
          <div class="detail-section animate__animated animate__fadeInRight" style="animation-delay: 0.2s">
            <h4 class="section-title"><el-icon><Document /></el-icon> 基本信息</h4>
            <div class="detail-item">
              <span class="detail-label"><el-icon><Location /></el-icon> IP地址:</span>
              <span class="detail-value">{{ currentWorker || 'N/A' }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <template v-slot:footer>
        <div class="drawer-footer" v-if="currentWorker">
          <el-button @click="drawerVisible = false" class="btn-pulse">
            <el-icon><Close /></el-icon> 关闭
          </el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useWorkerStore } from '@/stores/worker'
import { storeToRefs } from 'pinia'
import moment from 'moment'
import {
  Refresh,
  Monitor,
  View,
  Close,
  Loading,
  Document,
  Location,
  CircleCheckFilled,
  CircleCheck,
  Timer,
  Connection
} from '@element-plus/icons-vue'

export default {
  name: 'WorkersView',
  components: {
    Refresh,
    Monitor,
    View,
    Close,
    Loading,
    Document,
    Location,
    CircleCheckFilled,
    CircleCheck,
    Timer,
    Connection
  },
  setup() {
    const workerStore = useWorkerStore()
    const { workers, loading } = storeToRefs(workerStore)
    
    const drawerVisible = ref(false)
    const currentWorker = ref(null)
    
    // 获取Worker节点列表
    const fetchWorkers = async () => {
      await workerStore.fetchWorkers()
    }
    
    // 刷新节点列表
    const refreshWorkers = () => {
      fetchWorkers()
    }
    
    // 显示节点详情
    const showWorkerDetail = (worker) => {
      currentWorker.value = worker
      drawerVisible.value = true
    }
    
    // 获取节点状态类型
    const getStatusType = () => {
      return 'success'
    }
    
    // 获取节点状态文本
    const getStatusText = () => {
      return '在线'
    }
    
    // 判断节点是否在线
    const isOnline = (lastSeen) => {
      if (!lastSeen) return false;
      // lastSeen 是时间戳（秒或毫秒），统一转为毫秒
      const now = Date.now();
      const last = lastSeen > 1e12 ? lastSeen : lastSeen * 1000;
      // 1分钟内视为在线
      return now - last < 60 * 1000;
    }
    
    // 获取节点状态图标
    const getStatusIcon = () => {
      return 'CircleCheck'
    }
    
    // 获取节点状态样式类名
    const getStatusClass = () => {
      return 'status-online'
    }
    
    // 获取节点状态详细描述
    const getStatusDetail = () => {
      return '正常运行中'
    }
    
    // 获取最近更新时间文本
    const getLastUpdateTimeText = () => {
      return moment().format('HH:mm:ss')
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
    
    onMounted(() => {
      fetchWorkers()
      
      // 设置定时刷新
      const timer = setInterval(() => {
        fetchWorkers()
      }, 30000) // 每30秒刷新一次
      
      // 组件卸载时清除定时器
      return () => {
        clearInterval(timer)
      }
    })
    
    return {
      workers,
      loading,
      drawerVisible,
      currentWorker,
      refreshWorkers,
      showWorkerDetail,
      getStatusType,
      getStatusText,
      isOnline,
      getStatusIcon,
      getStatusClass,
      getStatusDetail,
      getLastUpdateTimeText,
      getProgressColor
    }
  }
}
</script>

<style scoped>
.workers .subtitle {
  color: #909399;
  margin-top: 0;
}

.worker-card {
  margin-bottom: 20px;
  overflow: hidden;
}

.worker-header {
  background-color: #f9fafc;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.worker-info {
  display: flex;
  align-items: center;
}

.worker-avatar {
  color: #fff;
  background-color: var(--primary-color);
  margin-right: 12px;
}

.worker-name-container {
  display: flex;
  flex-direction: column;
}

.worker-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.worker-body {
  padding: 16px;
}

.info-row {
  display: flex;
  margin-bottom: 8px;
}

.info-label {
  width: 80px;
  color: #909399;
}

.info-value {
  flex: 1;
  font-weight: 500;
}

.worker-footer {
  border-top: 1px solid var(--border-color);
  padding: 10px 16px;
  display: flex;
  justify-content: flex-end;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: #909399;
}

.loading-icon {
  font-size: 32px;
  margin-bottom: 16px;
}

.worker-detail-header {
  display: flex;
  align-items: center;
  padding: 16px;
}

.worker-detail-avatar {
  color: #fff;
  background-color: var(--primary-color);
  margin-right: 16px;
}

.worker-detail-info {
  flex: 1;
}

.worker-detail-name {
  margin: 0 0 8px 0;
}

.worker-detail-content {
  padding: 16px;
}

.detail-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--border-color);
}

.detail-item {
  margin-bottom: 16px;
}

.detail-label {
  display: flex;
  align-items: center;
  color: #909399;
  margin-bottom: 8px;
}

.detail-label .el-icon {
  margin-right: 6px;
  color: var(--primary-color);
}

.detail-value {
  font-weight: 500;
}

/* 修复进度条黑线问题 */
.progress-wrapper {
  width: 100%;
  padding: 2px;
  background-color: transparent;
  /* 防止ResizeObserver问题 */
  position: relative;
  contain: content;
  min-height: 24px;
}

/* 消除进度条中的黑线 */
.progress-wrapper :deep(.el-progress-bar__outer) {
  background-color: #EBEEF5;
  border-radius: 100px;
  border: none;
  outline: none;
  overflow: hidden; /* 防止溢出 */
}

.progress-wrapper :deep(.el-progress-bar__inner) {
  border-radius: 100px;
  transition: width 0.3s ease;
}

/* 确保进度条两边无边框 */
.progress-wrapper :deep(.el-progress) {
  line-height: normal;
  width: 100%;
  display: block;
}

/* 移除可能导致黑线的边框和阴影 */
.progress-wrapper :deep(.el-progress-bar) {
  padding: 0;
  margin: 0;
  border: none;
  background-color: transparent;
  transform: translateZ(0); /* 启用GPU加速，减少浏览器重绘 */
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
}

/* 节点卡片hover效果增强 */
.worker-card:hover .worker-avatar {
  transform: scale(1.1);
  transition: transform 0.3s ease;
}

/* 节点列表标题 */
.section-title {
  font-size: 18px;
  font-weight: 500;
  margin: 20px 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
}

.section-title .el-icon {
  margin-right: 8px;
  color: var(--primary-color);
}

/* 节点状态标签 */
.status-online {
  color: var(--success-color);
}

.status-offline {
  color: var(--danger-color);
}
</style>
