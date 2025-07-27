<template>
  <div class="home">
    <el-row>
      <el-col :span="24">
        <!-- Hero Section -->
        <div class="hero-section">
          <h1>专业的云端3D打印服务</h1>
          <p>上传您的3D模型，选择材料和质量，我们为您提供高品质的打印服务</p>
          <el-button type="primary" size="large" @click="$router.push('/upload')">
            开始打印
          </el-button>
        </div>
      </el-col>
    </el-row>

    <!-- Features Section -->
    <el-row :gutter="20" style="margin-top: 60px;">
      <el-col :span="8">
        <el-card class="feature-card">
          <div class="feature-icon">
            <el-icon size="40"><upload /></el-icon>
          </div>
          <h3>简单上传</h3>
          <p>支持STL、OBJ格式，最大50MB，拖拽即可上传</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="feature-card">
          <div class="feature-icon">
            <el-icon size="40"><view /></el-icon>
          </div>
          <h3>3D预览</h3>
          <p>实时3D预览，360度查看您的模型效果</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="feature-card">
          <div class="feature-icon">
            <el-icon size="40"><price-tag /></el-icon>
          </div>
          <h3>即时报价</h3>
          <p>根据体积和材料自动计算价格，透明收费</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- Materials Section -->
    <div class="materials-section">
      <h2>支持的打印材料</h2>
      <el-row :gutter="20">
        <el-col :span="6" v-for="material in materials" :key="material.name">
          <el-card class="material-card">
            <h4>{{ material.name }}</h4>
            <p>{{ material.description }}</p>
            <div class="price">￥{{ material.price }}/克</div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'Home',
  setup() {
    const materials = ref([
      {
        name: 'PLA',
        description: '环保易用，适合初学者',
        price: '0.8'
      },
      {
        name: 'ABS',
        description: '耐用性强，适合功能件',
        price: '1.0'
      },
      {
        name: 'PETG',
        description: '透明度好，化学稳定',
        price: '1.2'
      },
      {
        name: 'TPU',
        description: '柔性材料，适合弹性件',
        price: '2.0'
      }
    ])

    return {
      materials
    }
  }
}
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.hero-section {
  text-align: center;
  padding: 80px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 10px;
  margin-bottom: 40px;
}

.hero-section h1 {
  font-size: 3rem;
  margin-bottom: 20px;
  font-weight: 300;
}

.hero-section p {
  font-size: 1.2rem;
  margin-bottom: 30px;
  opacity: 0.9;
}

.feature-card {
  text-align: center;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.feature-icon {
  margin-bottom: 20px;
  color: #409eff;
}

.feature-card h3 {
  margin: 0 0 15px 0;
  color: #333;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}

.materials-section {
  margin-top: 80px;
  text-align: center;
}

.materials-section h2 {
  margin-bottom: 40px;
  color: #333;
}

.material-card {
  height: 150px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.material-card h4 {
  margin: 0 0 10px 0;
  color: #409eff;
  font-size: 1.2rem;
}

.material-card p {
  color: #666;
  margin: 0 0 15px 0;
}

.price {
  font-size: 1.1rem;
  font-weight: bold;
  color: #e6a23c;
}
</style>