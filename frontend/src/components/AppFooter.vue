<template>
  <div class="app-footer">
    <el-row>
      <el-col :span="24" style="text-align: center;">
        <p>&copy; 2024 3D打印农场. 保留所有权利.</p>
        <p>联系我们: info@3dprint-farm.com | 客服热线: 400-123-4567</p>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'AppFooter'
}
</script>

<style scoped>
.app-footer {
  padding: 20px;
  background-color: #f5f5f5;
  color: #666;
}

.app-footer p {
  margin: 5px 0;
  font-size: 14px;
}
</style>