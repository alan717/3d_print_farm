<template>
  <div class="app-footer">
    <el-row>
      <el-col :span="24" style="text-align: center;">
        <p>&copy; 2025 3D打印农场. 保留所有权利.</p>
        <p>联系我们: info@3dprint-farm.com | 客服热线: 400-123-4567</p>
        <div class="footer-links">
          <router-link to="/admin/login" class="admin-link">管理员登录</router-link>
        </div>
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

.footer-links {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #e0e0e0;
}

.admin-link {
  color: #999;
  text-decoration: none;
  font-size: 12px;
  transition: color 0.3s;
}

.admin-link:hover {
  color: #409eff;
}
</style>