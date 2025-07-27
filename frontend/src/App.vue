<template>
  <div id="app">
    <el-container>
      <el-header>
        <nav-header />
      </el-header>
      <el-main>
        <router-view />
      </el-main>
      <el-footer>
        <app-footer />
      </el-footer>
    </el-container>
  </div>
</template>

<script>
import NavHeader from '@/components/NavHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

export default {
  name: 'App',
  components: {
    NavHeader,
    AppFooter
  }
}
</script>

<style>
#app {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  min-height: 100vh;
}

.el-header {
  padding: 0;
  height: 60px !important;
}

.el-footer {
  height: 80px !important;
  background-color: #f5f5f5;
}

.el-main {
  padding: 20px;
  min-height: calc(100vh - 140px);
}
</style>