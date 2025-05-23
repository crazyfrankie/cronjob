<template>
  <router-view v-slot="{ Component }">
    <transition name="page-slide" mode="out-in">
      <component :is="Component" />
    </transition>
  </router-view>
</template>

<script>
export default {
  name: 'App',
  mounted() {
    // 添加页面加载完成动画
    document.body.classList.add('app-loaded');
  }
}
</script>

<style>
/* 全局样式在 assets/styles/main.css 中 */
.app-loaded {
  animation: fadeIn 0.5s ease-in;
}

/* 页面过渡效果 */
.page-slide-enter-active,
.page-slide-leave-active {
  transition: all 0.3s ease;
}

.page-slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.page-slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #00ADD8;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #007D9C;
}
</style>
