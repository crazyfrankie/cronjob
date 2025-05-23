import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import './assets/styles/main.css'
import './assets/styles/animations.css'
import './assets/styles/gopher.css'
import './assets/styles/form-animations.css'
import './assets/styles/responsive.css'
import './assets/styles/fix-styles.css'
import './assets/styles/footer-fix.css'
import './assets/styles/footer-scroll.css'

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 添加全局错误处理
window.addEventListener('error', (event) => {
  // 忽略ResizeObserver循环错误
  if (event && event.message && (
    event.message.includes('ResizeObserver loop') || 
    event.message.includes('ResizeObserver Loop')
  )) {
    event.stopImmediatePropagation();
    event.preventDefault();
    return false;
  }
});

// 控制台错误处理
const originalConsoleError = console.error;
console.error = (...args) => {
  // 忽略ResizeObserver循环错误和相关错误
  if (args[0] && typeof args[0] === 'string' && (
    args[0].includes('ResizeObserver loop') || 
    args[0].includes('ResizeObserver Loop') ||
    args[0].includes('ResizeObserver completed with undelivered notifications')
  )) {
    return;
  }
  originalConsoleError(...args);
};

app.use(ElementPlus)
app.use(createPinia())
app.use(router)

app.mount('#app')
