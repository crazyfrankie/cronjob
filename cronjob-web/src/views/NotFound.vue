<template>
  <div class="not-found-container">
    <div class="not-found-content animate__animated animate__fadeIn">
      <div class="gopher-container animate__animated animate__bounce">
        <img src="@/assets/images/gopher-404.png" alt="404" class="not-found-image" />
      </div>
      <h1 class="not-found-title animate__animated animate__fadeInDown">404</h1>
      <h2 class="not-found-subtitle animate__animated animate__fadeInUp" style="animation-delay: 0.2s">页面未找到</h2>
      <p class="not-found-desc animate__animated animate__fadeInUp" style="animation-delay: 0.3s">抱歉，您访问的页面不存在或已被移除</p>
      
      <div class="not-found-actions animate__animated animate__fadeInUp" style="animation-delay: 0.5s">
        <el-button type="primary" @click="goHome" class="btn-hover-lift">
          <el-icon><HomeFilled /></el-icon>
          返回首页
        </el-button>
        
        <el-button plain @click="goBack" class="btn-hover-lift">
          <el-icon><Back /></el-icon>
          返回上一页
        </el-button>
      </div>
      
      <div class="cron-animation">
        <div class="clock-hand"></div>
      </div>
    </div>
    
    <!-- 背景动画元素 -->
    <div class="background-elements">
      <div class="bg-element e1"></div>
      <div class="bg-element e2"></div>
      <div class="bg-element e3"></div>
      <div class="bg-element e4"></div>
      <div class="bg-element e5"></div>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'NotFound',

  setup() {
    const router = useRouter()

    const goHome = () => {
      router.push('/')
    }
    
    const goBack = () => {
      router.back()
    }

    return {
      goHome,
      goBack
    }
  }
}
</script>

<style scoped>
.not-found-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
  position: relative;
  overflow: hidden;
}

.not-found-content {
  text-align: center;
  padding: 40px;
  max-width: 600px;
  background-color: white;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  position: relative;
  z-index: 10;
}

.gopher-container {
  display: inline-block;
  animation-duration: 2s;
  animation-iteration-count: infinite;
  transform-origin: bottom center;
}

.not-found-image {
  width: 220px;
  filter: drop-shadow(0px 4px 8px rgba(0, 0, 0, 0.2));
  transition: all 0.3s ease;
}

.not-found-image:hover {
  transform: rotate(5deg) scale(1.05);
}

.not-found-title {
  font-size: 92px;
  font-weight: 800;
  background: linear-gradient(135deg, var(--primary-color) 0%, #67c23a 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
  line-height: 1;
  letter-spacing: -2px;
}

.not-found-subtitle {
  font-size: 32px;
  font-weight: 600;
  margin: 16px 0;
  color: #333;
}

.not-found-desc {
  color: #606266;
  font-size: 18px;
  margin-bottom: 30px;
}

.not-found-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
}

.btn-hover-lift {
  transition: all 0.3s ease;
}

.btn-hover-lift:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
}

/* 时钟动画 */
.cron-animation {
  position: absolute;
  bottom: 20px;
  right: 20px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid var(--primary-color);
  opacity: 0.6;
}

.clock-hand {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 18px;
  height: 2px;
  background-color: var(--primary-color);
  transform-origin: left center;
  animation: rotate 5s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 背景元素动画 */
.background-elements {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  pointer-events: none;
}

.bg-element {
  position: absolute;
  background-color: var(--primary-color);
  opacity: 0.05;
  border-radius: 50%;
}

.e1 {
  width: 300px;
  height: 300px;
  top: -150px;
  left: -150px;
  animation: float 15s ease-in-out infinite alternate;
}

.e2 {
  width: 200px;
  height: 200px;
  top: 70%;
  right: -50px;
  animation: float 12s ease-in-out infinite alternate-reverse;
  background-color: #67c23a;
}

.e3 {
  width: 100px;
  height: 100px;
  bottom: 10%;
  left: 15%;
  animation: float 8s ease-in-out infinite alternate;
  background-color: #e6a23c;
}

.e4 {
  width: 120px;
  height: 120px;
  top: 20%;
  right: 15%;
  animation: float 10s ease-in-out infinite alternate-reverse;
  background-color: #f56c6c;
}

.e5 {
  width: 80px;
  height: 80px;
  top: 40%;
  left: 10%;
  animation: float 6s ease-in-out infinite alternate;
  background-color: #909399;
}

@keyframes float {
  0% {
    transform: translateY(0) translateX(0);
  }
  50% {
    transform: translateY(30px) translateX(20px) rotate(5deg);
  }
  100% {
    transform: translateY(-20px) translateX(-15px) rotate(-5deg);
  }
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .not-found-content {
    padding: 30px 20px;
    max-width: 90%;
  }
  
  .not-found-image {
    width: 150px;
  }
  
  .not-found-title {
    font-size: 72px;
  }
  
  .not-found-subtitle {
    font-size: 24px;
  }
  
  .not-found-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .not-found-actions .el-button {
    width: 100%;
    margin-bottom: 10px;
  }
}
</style>
