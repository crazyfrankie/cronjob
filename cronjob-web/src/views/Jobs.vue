<template>
  <div class="page-container jobs">
    <div class="page-header flex-between animate__animated animate__fadeIn">
      <div>
        <h2><el-icon><Calendar /></el-icon> 任务管理</h2>
        <p class="subtitle">管理所有定时任务，支持创建、编辑和删除操作</p>
      </div>
      <div>
        <el-button type="primary" @click="openJobDialog(null)" class="btn-pulse">
          <el-icon><Plus /></el-icon> 创建任务
        </el-button>
      </div>
    </div>
    
    <!-- 搜索栏 -->
    <div class="search-bar animate__animated animate__fadeInDown" style="animation-delay: 0.2s">
      <el-input
        v-model="searchQuery"
        placeholder="搜索任务名称"
        clearable
        style="width: 240px"
        @input="handleSearchInput"
        class="search-input"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
        <template #append>
          <el-button>
            <el-icon><Search /></el-icon>
          </el-button>
        </template>
      </el-input>
    </div>
    
    <!-- 任务表格 -->
    <el-card shadow="hover" class="job-table-card custom-card animate__animated animate__fadeIn" style="animation-delay: 0.3s">
      <el-table
        v-loading="loading"
        :data="filteredJobs"
        stripe
        style="width: 100%"
        :empty-text="loading ? '加载中...' : '暂无任务数据'"
        @row-click="handleRowClick"
      >
        <el-table-column prop="name" label="任务名称" min-width="100" max-width="140" show-overflow-tooltip>
          <template #default="scope">
            <div class="job-name-cell">
              <el-icon class="job-icon"><Calendar /></el-icon>
              <span>{{ scope.row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="command" label="命令" min-width="200" max-width="300" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small" type="info" class="command-tag" style="max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
              <el-icon><Terminal /></el-icon>
              <span class="command-text" @click.stop="showDetail('命令', scope.row.command)">{{ scope.row.command }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cronExpr" label="Cron表达式" min-width="180" max-width="240" show-overflow-tooltip>
          <template #default="scope">
            <el-tooltip :content="describeCron(scope.row.cronExpr)" placement="top">
              <el-tag size="small" class="cron-tag" effect="plain" style="max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" @click.stop="showDetail('Cron表达式', scope.row.cronExpr)">{{ scope.row.cronExpr }}</el-tag>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button type="primary" plain size="small" @click.stop="openJobDialog(scope.row)" class="btn-hover-lift">
                <el-icon><Edit /></el-icon> 编辑
              </el-button>
              <el-button type="danger" plain size="small" @click.stop="confirmDeleteJob(scope.row)" class="btn-hover-lift">
                <el-icon><Delete /></el-icon> 删除
              </el-button>
              <el-button type="warning" plain size="small" @click.stop="confirmKillJob(scope.row)" class="btn-hover-lift">
                <el-icon><CircleClose /></el-icon> 终止
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="table-pagination" v-if="filteredJobs.length > 0">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="filteredJobs.length"
          :page-size="10"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
    
    <!-- 任务表单对话框 -->
    <el-dialog
      v-model="jobDialogVisible"
      :title="currentJob.id ? '编辑任务' : '创建任务'"
      width="600px"
      @closed="resetJobForm"
      custom-class="job-dialog animate__animated animate__fadeInUp fancy-dialog"
      destroy-on-close
    >
      <div class="dialog-content-wrapper">
        <el-form
          ref="jobFormRef"
          :model="jobForm"
          :rules="jobRules"
          label-width="100px"
          class="job-form"
          @submit.prevent="submitJobForm"
        >
          <div class="form-group animate__animated animate__fadeInUp" style="animation-delay: 0.1s">
            <div class="form-group-title">
              <el-icon><InfoFilled /></el-icon>
              <span>基本信息</span>
            </div>
            <el-form-item label="任务名称" prop="name">
              <div class="fancy-input">
                <el-input 
                  v-model="jobForm.name" 
                  placeholder="请输入任务名称"
                  :class="{'form-error': formErrors.name}" 
                  clearable
                >
                  <template #prefix>
                    <el-icon><Document /></el-icon>
                  </template>
                </el-input>
              </div>
              <div class="form-tip" v-if="!formErrors.name">输入一个易于识别的任务名称</div>
            </el-form-item>
          </div>
          
          <div class="form-group animate__animated animate__fadeInUp" style="animation-delay: 0.2s">
            <div class="form-group-title">
              <el-icon><Terminal /></el-icon>
              <span>执行配置</span>
            </div>
            <el-form-item label="执行命令" prop="command">
              <div class="fancy-input">
                <el-input
                  v-model="jobForm.command"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入要执行的命令"
                  :class="{'form-error': formErrors.command}"
                  style="font-family: 'Courier New', monospace;"
                />
              </div>
              <div class="form-tip" v-if="!formErrors.command">支持任何shell命令或脚本路径</div>
            </el-form-item>
            
            <el-form-item label="Cron表达式" prop="cronExpr">
              <div class="fancy-input cron-input-container">
                <el-input 
                  v-model="jobForm.cronExpr" 
                  placeholder="例如: */5 * * * * 表示每5分钟执行一次" 
                  class="cron-input"
                  :class="{'form-error': formErrors.cronExpr}"
                  @blur="validateCronExpression"
                >
                  <template #append>
                    <el-tooltip
                      content="Cron表达式格式: 分 时 日 月 周，如 0 0 * * * 表示每天0点执行"
                      placement="top"
                    >
                      <el-icon><QuestionFilled /></el-icon>
                    </el-tooltip>
                  </template>
                </el-input>
              </div>
            </el-form-item>
          </div>
        </el-form>
        
        <!-- 使用说明 -->
        <div class="usage-tips animate__animated animate__fadeIn" style="animation-delay: 0.4s">
          <el-alert
            title="Cron表达式使用说明"
            type="info"
            description="Cron表达式由5个字段组成: 分钟 小时 日期 月份 星期。例如 '*/10 * * * *' 表示每10分钟执行一次。"
            show-icon
            :closable="false"
            class="tip-alert"
          />
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="jobDialogVisible = false" class="btn-cancel btn-hover-lift">
            <el-icon><Close /></el-icon> 取消
          </el-button>
          <el-button 
            type="primary" 
            @click="submitJobForm" 
            :loading="submitLoading" 
            class="submit-btn btn-hover-lift btn-submit"
          >
            <el-icon><Check /></el-icon>
            {{ currentJob.id ? '保存' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <el-dialog v-model="detailDialog" :title="detailTitle" width="500px">
      <pre style="white-space: pre-wrap; word-break: break-all;">{{ detailContent }}</pre>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted, reactive } from 'vue'
import { useJobStore } from '@/stores/job'
import { ElMessage, ElMessageBox } from 'element-plus'
// import { Search } from '@element-plus/icons-vue'
import moment from 'moment'
import { storeToRefs } from 'pinia'
import cronstrue from 'cronstrue/i18n'
import '@/assets/styles/jobs.css'

export default {
  name: 'JobsView',
  
  setup() {
    const jobStore = useJobStore()
    const { jobs, loading } = storeToRefs(jobStore)
    
    const jobDialogVisible = ref(false)
    const searchQuery = ref('')
    const currentJob = ref({})
    const submitLoading = ref(false)
    const jobFormRef = ref(null)
    const searchTimer = ref(null)
    const isValidCron = ref(true)
    const detailDialog = ref(false)
    const detailTitle = ref('')
    const detailContent = ref('')
    
    // 提供Search图标
    // const SearchIcon = Search
    
    // 表单验证错误标记
    const formErrors = reactive({
      name: false,
      command: false,
      cronExpr: false
    })
    
    // 任务表单
    const jobForm = reactive({
      name: '',
      command: '',
      cronExpr: ''
    })
    
    // 表单验证规则
    const jobRules = {
      name: [
        { required: true, message: '请输入任务名称', trigger: 'blur' },
        { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
      ],
      command: [
        { required: true, message: '请输入执行命令', trigger: 'blur' }
      ],
      cronExpr: [
        { required: true, message: '请输入Cron表达式', trigger: 'blur' }
      ]
    }
    
    // 过滤后的任务列表
    const filteredJobs = computed(() => {
      if (!searchQuery.value) return jobs.value
      
      return jobs.value.filter(job => 
        job.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        job.command.toLowerCase().includes(searchQuery.value.toLowerCase())
      )
    })
    
    // 打开任务表单对话框
    const openJobDialog = (job) => {
      currentJob.value = job || {}
      
      if (job) {
        jobForm.name = job.name
        jobForm.command = job.command
        jobForm.cronExpr = job.cronExpr
      }
      
      jobDialogVisible.value = true
    }
    
    // 重置表单
    const resetJobForm = () => {
      if (jobFormRef.value) {
        jobFormRef.value.resetFields()
      }
      
      jobForm.name = ''
      jobForm.command = ''
      jobForm.cronExpr = ''
      currentJob.value = {}
      
      // 重置表单错误状态
      Object.keys(formErrors).forEach(key => {
        formErrors[key] = false
      })
      
      isValidCron.value = true
    }
    
    // 验证Cron表达式
    const validateCronExpression = () => {
      if (!jobForm.cronExpr) {
        formErrors.cronExpr = true;
        isValidCron.value = false;
        return false;
      }
      // 只做基本的格式检查：确保是6个字段
      const parts = jobForm.cronExpr.trim().split(/\s+/);
      if (parts.length !== 6) {
        formErrors.cronExpr = true;
        isValidCron.value = false;
        return false;
      }
      
      // 基本格式正确，具体验证交给后端
      formErrors.cronExpr = false;
      isValidCron.value = true;
      return true;
    };
    
    // 提交表单
    const submitJobForm = async () => {
      if (!jobFormRef.value) return
      
      await jobFormRef.value.validate(async (valid) => {
        if (!valid) {
          // 标记错误字段
          const fields = jobFormRef.value.fields
          fields.forEach(field => {
            if (field.validateState === 'error') {
              formErrors[field.prop] = true
            }
          })
          
          ElMessage({
            message: '请完善表单信息',
            type: 'warning',
            duration: 2000,
            showClose: true
          })
          return
        }
        
        submitLoading.value = true
        
        try {
          const jobData = {
            name: jobForm.name,
            command: jobForm.command,
            cronExpr: jobForm.cronExpr
          }
          
          if (currentJob.value.id) {
            jobData.id = currentJob.value.id
          }
          
          const success = await jobStore.saveOrUpdateJob(jobData)
          
          if (success) {
            ElMessage({
              message: currentJob.value.id ? '任务更新成功！' : '任务创建成功！',
              type: 'success',
              duration: 2000
            })
            jobDialogVisible.value = false
            resetJobForm()
          }
        } catch (error) {
          ElMessage({
            message: '操作失败：' + (error.message || '未知错误'),
            type: 'error',
            duration: 2000
          })
        } finally {
          submitLoading.value = false
        }
      })
    }
    
    // 确认删除任务
    const confirmDeleteJob = (job) => {
      ElMessageBox.confirm(
        `确定要删除任务 "${job.name}" 吗？此操作不可恢复！`,
        '删除确认',
        {
          confirmButtonText: '删除',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(async () => {
          await jobStore.removeJob(job.name)
        })
        .catch(() => {
          ElMessage({
            type: 'info',
            message: '已取消删除操作'
          })
        })
    }
    
    // 确认终止任务
    const confirmKillJob = (job) => {
      ElMessageBox.confirm(
        `确定要终止任务 "${job.name}" 的当前执行实例吗？`,
        '终止确认',
        {
          confirmButtonText: '终止',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
        .then(async () => {
          await jobStore.stopJob(job.name)
        })
        .catch(() => {
          ElMessage({
            type: 'info',
            message: '已取消终止操作'
          })
        })
    }
    
    // 处理搜索输入
    const handleSearchInput = () => {
      // 添加节流逻辑
      clearTimeout(searchTimer.value)
      searchTimer.value = setTimeout(() => {
        // 搜索操作
      }, 300)
    }
    
    // 获取下次执行时间
    const getNextRunTime = (cronExpr) => {
      if (!cronExpr) return '-'
      
      try {
        const parser = require('cron-parser')
        const interval = parser.parseExpression(cronExpr)
        return moment(interval.next().toDate()).format('YYYY-MM-DD HH:mm:ss')
      } catch (err) {
        return '表达式无效'
      }
    }
    
    // 设置表格行的class
    const tableRowClassName = ({ rowIndex }) => {
      return rowIndex % 2 === 0 ? 'row-even' : 'row-odd'
    }
    
    // 处理行点击
    const handleRowClick = (row) => {
      // 可以添加行选中效果
      console.log('Row clicked:', row)
    }
    
    // 处理分页变化
    const handlePageChange = (currentPage) => {
      console.log('Current page:', currentPage)
    }
    
    // 描述Cron表达式
    const describeCron = (cronExpr) => {
      if (!cronExpr) return '无效的 cron 表达式';
      try {
        // 使用 cronstrue 解析 6 位 cron 表达式
        return cronstrue.toString(cronExpr, { locale: 'zh_CN' });
      } catch (error) {
        return '无效的 cron 表达式';
      }
    };
    
    const showDetail = (title, content) => {
      detailTitle.value = title
      detailContent.value = content
      detailDialog.value = true
    }
    
    onMounted(() => {
      jobStore.fetchJobs()
    })
    
    return {
      jobs,
      loading,
      filteredJobs,
      jobDialogVisible,
      searchQuery,
      jobForm,
      jobRules,
      currentJob,
      submitLoading,
      jobFormRef,
      formErrors,
      isValidCron,
      openJobDialog,
      resetJobForm,
      submitJobForm,
      confirmDeleteJob,
      confirmKillJob,
      handleSearchInput,
      getNextRunTime,
      describeCron,
      tableRowClassName,
      handleRowClick,
      handlePageChange,
      validateCronExpression,
      showDetail,
      detailDialog,
      detailTitle,
      detailContent
    }
  }
}
</script>

<style scoped>
.jobs .subtitle {
  color: #909399;
  margin-top: 0;
}

.search-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-start;
}

.job-table-card {
  margin-left: 0 !important;
  padding-left: 0 !important;
  width: 100%;
  box-sizing: border-box;
}

.job-table-card .el-table {
  width: 100% !important;
  margin-left: 0 !important;
}

.job-name-cell {
  display: flex;
  align-items: center;
}

.job-icon {
  color: #409eff;
  margin-right: 8px;
}

.command-tag {
  border-radius: 16px;
}

.command-text {
  margin-left: 4px;
}

.cron-tag {
  border-radius: 16px;
}

.next-run-time {
  display: flex;
  align-items: center;
}

.next-run-time el-icon {
  color: #67c23a;
  margin-right: 4px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: flex-start;
  min-width: 240px;
  max-width: 270px;
}

.action-buttons .el-button {
  min-width: 70px;
  max-width: 100px;
  flex: 0 0 auto;
  padding: 0 10px;
  font-size: 14px;
  height: 32px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-hover-lift {
  transition: transform 0.2s;
}

.btn-hover-lift:hover {
  transform: translateY(-2px);
}

.table-pagination {
  margin-left: 0 !important;
  padding-left: 0 !important;
  width: 100%;
  justify-content: flex-end;
}

.table-pagination .el-pagination {
  min-width: 340px;
  display: inline-flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.submit-btn {
  position: relative;
}

.pulse-animation {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
  100% {
    transform: scale(1);
  }
}

.text-muted {
  color: #909399;
  font-size: 14px;
}

.row-even {
  background-color: #f9f9f9;
}

.row-odd {
  background-color: #ffffff;
}

.fancy-dialog .el-dialog__header {
  background-color: #f5f7fa;
  border-bottom: 1px solid #e4e7ec;
  padding: 16px;
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
}

.fancy-dialog .el-dialog__title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.fancy-dialog .el-dialog__close {
  color: #909399;
}

.dialog-content-wrapper {
  padding: 24px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group-title {
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.form-group-title el-icon {
  color: #409eff;
  margin-right: 8px;
}

.fancy-input {
  position: relative;
}

.fancy-input .el-input {
  padding-right: 40px;
}

.fancy-input .el-input__prefix,
.fancy-input .el-input__suffix {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
}

.fancy-input .el-input__prefix {
  left: 10px;
}

.fancy-input .el-input__suffix {
  right: 10px;
}

.form-error .el-input {
  border-color: #f56c6c;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.cron-input-container {
  position: relative;
}

.cron-input-container .el-input {
  padding-right: 40px;
}

.cron-input-container .el-input__suffix {
  position: absolute;
  top: 50%;
  right: 10px;
  transform: translateY(-50%);
}

.cron-description-box {
  margin-top: 8px;
  padding: 12px;
  background-color: #f9f9f9;
  border-radius: 4px;
  border: 1px solid #e4e7ec;
}

.cron-description-title {
  display: flex;
  align-items: center;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.cron-next-times {
  display: flex;
  flex-direction: column;
}

.next-time-item {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.tip-alert {
  margin-top: 16px;
}
</style>
