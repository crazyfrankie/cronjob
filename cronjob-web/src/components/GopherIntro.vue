<template>
  <div class="gopher-intro animate__animated animate__fadeIn">
    <div class="gopher-logo-container float-animation">
      <img src="/zcron-logo.svg" alt="ZCron Logo" class="gopher-logo" />
    </div>
    <div class="gopher-text animate__animated animate__fadeInUp">
      <h1>欢迎使用 ZCron</h1>
      <p>基于 Go 的分布式定时任务管理系统</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GopherIntro'
}
</script>

<style scoped>
.gopher-intro {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  text-align: center;
  margin-bottom: 30px;
}

.gopher-logo-container {
  margin-bottom: 20px;
}

.gopher-logo {
  width: 120px;
  height: 120px;
}

.gopher-text h1 {
  font-size: 32px;
  color: var(--primary-color);
  margin-bottom: 10px;
}

.gopher-text p {
  font-size: 18px;
  color: #666;
}
</style>
