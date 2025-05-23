<template>
  <div class="main-layout">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="collapsed ? '64px' : '220px'" class="aside" :class="{ 'collapsed-menu': collapsed }">
        <div class="logo-container">
          <img src="/zcron-logo.svg" alt="ZCron Logo" class="logo-img animate__animated animate__pulse" />
          <h2 class="logo-text" v-show="!collapsed">ZCron</h2>
        </div>
        
        <el-menu
          router
          :default-active="activeMenu"
          class="menu"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#00ADD8"
          :collapse="collapsed">

          <!-- 菜单项 -->
          <el-menu-item v-for="item in routes" :key="item.path" :index="item.path">
            <el-icon><component :is="item.meta.icon" /></el-icon>
            <span>{{ item.meta.title }}</span>
          </el-menu-item>
        </el-menu>
        
        <div class="menu-footer" v-show="!collapsed">
          <div class="go-version">
            <img src="/zcron-logo.svg" alt="Gopher" class="mini-gopher" />
            <span>ZCron v1.0.0</span>
          </div>
        </div>
      </el-aside>
      
      <!-- 主区域 -->
      <el-container class="main-container" style="border: none !important;">
        <!-- 顶部导航 -->
        <el-header class="header">
          <div class="header-left">
            <el-button size="small" circle @click="toggleSidebar">
              <el-icon><Fold v-if="!collapsed" /><Expand v-else /></el-icon>
            </el-button>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ currentRoute.title }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          
          <div class="header-right">
            <span class="date-time">{{ currentDateTime }}</span>
            <el-dropdown trigger="click">
              <span class="user-info">
                <el-avatar size="small" src="https://avatars.githubusercontent.com/u/1479100?v=4" />
                <span class="username">Admin</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item>个人设置</el-dropdown-item>
                  <el-dropdown-item divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        
        <!-- 创建一个包含主内容和页脚的滚动容器 -->
        <div class="scroll-container">
          <!-- 内容区域 -->
          <el-main class="main">
            <router-view v-slot="{ Component }">
              <transition name="fade" mode="out-in">
                <component :is="Component" />
              </transition>
            </router-view>
          </el-main>
          
          <!-- 页脚，放在滚动容器内部的底部 -->
          <div class="footer" height="40px">
            <div class="footer-content">© 2025 ZCron 定时任务管理系统 | Powered by Go & Vue</div>
          </div>
        </div>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import moment from 'moment'

export default {
  name: 'MainLayout',
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    const collapsed = ref(false)
    const currentDateTime = ref('')
    let timer = null
    
    // 获取路由列表
    const routes = router.options.routes[0].children
    
    // 折叠/展开侧边栏
    const toggleSidebar = () => {
      collapsed.value = !collapsed.value
      // 保存用户的侧边栏偏好到 localStorage
      localStorage.setItem('sidebarCollapsed', collapsed.value)
    }
    
    // 从 localStorage 恢复侧边栏状态
    onMounted(() => {
      const savedState = localStorage.getItem('sidebarCollapsed')
      if (savedState !== null) {
        collapsed.value = savedState === 'true'
      }
    })
    
    // 当前激活的菜单项
    const activeMenu = computed(() => '/' + route.path.split('/')[1])
    
    // 当前路由信息
    const currentRoute = computed(() => {
      const path = '/' + route.path.split('/')[1]
      const found = routes.find(r => r.path === path.substring(1))
      return found ? found.meta : { title: '未知页面' }
    })
    
    // 更新当前时间
    const updateDateTime = () => {
      currentDateTime.value = moment().format('YYYY-MM-DD HH:mm:ss')
    }
    
    onMounted(() => {
      updateDateTime()
      timer = setInterval(updateDateTime, 1000)
    })
    
    onUnmounted(() => {
      if (timer) clearInterval(timer)
    })
    
    return {
      routes,
      activeMenu,
      currentRoute,
      collapsed,
      currentDateTime,
      toggleSidebar
    }
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
}

.aside {
  background-color: #304156;
  color: #bfcbd9;
  overflow: hidden;
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  transition: all 0.3s;
  z-index: 10;
}

.menu {
  border: none;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  border-bottom: 1px solid #283445;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #ffffff;
  white-space: nowrap;
  margin: 0;
}

.menu-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 16px;
  color: #909399;
  font-size: 12px;
  border-top: 1px solid #283445;
}

.go-version {
  display: flex;
  align-items: center;
}

.go-version i {
  margin-right: 5px;
  font-size: 14px;
}

.header {
  background-color: #ffffff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-left .el-button {
  margin-right: 16px;
}

.header-right {
  display: flex;
  align-items: center;
}

.date-time {
  margin-right: 20px;
  color: #606266;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 8px;
  height: 40px;
  border-radius: 4px;
  transition: all 0.3s;
}

.user-info:hover {
  background: #f2f2f2;
}

.username {
  margin: 0 8px;
  color: #606266;
}

.scroll-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto; /* 让这个容器可滚动 */
  overflow-x: hidden;
  height: calc(100vh - 64px);
}

.main {
  background-color: #f5f7fa;
  padding: 20px;
  flex: 1; /* 占用所有可用空间 */
  min-height: calc(100vh - 64px - 40px); /* 减去头部和页脚高度 */
  overflow: visible; /* 移除内部滚动条 */
}

.footer {
  background-color: #ffffff;
  color: #909399;
  text-align: center;
  line-height: 40px;
  height: 40px;
  font-size: 12px;
  border: none !important; /* 删除所有边线 */
  box-shadow: none !important; /* 移除所有阴影 */
  flex-shrink: 0; /* 防止页脚被压缩 */
}

.footer-content {
  border-top: none !important;
  box-shadow: none !important;
}

.main-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px);
  padding: 0;
  overflow: hidden; /* 防止外部容器滚动 */
}

.el-header {
  height: 64px;
}

.el-footer {
  height: 40px;
  margin-top: auto; /* 将页脚推到底部 */
  border: none !important;
  box-shadow: none !important;
}
</style>
