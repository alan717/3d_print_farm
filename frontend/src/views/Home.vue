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
    <el-row :gutter="20" class="features-section">
      <el-col :xs="24" :sm="12" :md="8" :lg="8" :xl="8" v-for="feature in features" :key="feature.title">
        <el-card class="feature-card">
          <div class="feature-icon">
            <el-icon size="40" :class="feature.icon"></el-icon>
          </div>
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.description }}</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- Materials Section -->
    <div class="materials-section">
      <h2>支持的打印材料</h2>
      <el-row :gutter="20">
        <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="6" v-for="material in materials" :key="material.name">
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
import { Upload, View, PriceTag } from '@element-plus/icons-vue'

export default {
  name: 'Home',
  components: {
    Upload,
    View,
    PriceTag
  },
  setup() {
    const features = ref([
      {
        title: '简单上传',
        description: '支持STL、OBJ格式，最大50MB，拖拽即可上传',
        icon: 'Upload'
      },
      {
        title: '3D预览',
        description: '实时3D预览，360度查看您的模型效果',
        icon: 'View'
      },
      {
        title: '即时报价',
        description: '根据体积和材料自动计算价格，透明收费',
        icon: 'PriceTag'
      }
    ])

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
      features,
      materials
    }
  }
}
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.hero-section {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 10px;
  margin-bottom: 40px;
  box-sizing: border-box;
}

.hero-section h1 {
  font-size: 2.5rem;
  margin-bottom: 20px;
  font-weight: 300;
  line-height: 1.2;
  word-break: keep-all;
  overflow-wrap: break-word;
}

.hero-section p {
  font-size: 1.1rem;
  margin-bottom: 30px;
  opacity: 0.9;
  line-height: 1.6;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 10px;
}

.features-section {
  margin-top: 60px;
}

.feature-card {
  text-align: center;
  height: 220px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-bottom: 20px;
  box-sizing: border-box;
}

.feature-icon {
  margin-bottom: 20px;
  color: #409eff;
}

.feature-card h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 1.2rem;
  word-break: keep-all;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
  padding: 0 10px;
  font-size: 0.95rem;
  word-break: keep-all;
  overflow-wrap: break-word;
}

.materials-section {
  margin-top: 80px;
  text-align: center;
}

.materials-section h2 {
  margin-bottom: 40px;
  color: #333;
  font-size: 1.8rem;
  word-break: keep-all;
}

.material-card {
  height: 150px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-bottom: 20px;
  box-sizing: border-box;
}

.material-card h4 {
  margin: 0 0 10px 0;
  color: #409eff;
  font-size: 1.2rem;
  word-break: keep-all;
}

.material-card p {
  color: #666;
  margin: 0 0 15px 0;
  font-size: 0.9rem;
  line-height: 1.4;
  padding: 0 8px;
  word-break: keep-all;
  overflow-wrap: break-word;
}

.price {
  font-size: 1.1rem;
  font-weight: bold;
  color: #e6a23c;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home {
    padding: 15px;
  }

  .hero-section {
    padding: 40px 15px;
    margin-bottom: 30px;
  }

  .hero-section h1 {
    font-size: 2rem;
    margin-bottom: 15px;
  }

  .hero-section p {
    font-size: 1rem;
    margin-bottom: 25px;
    padding: 0 5px;
  }

  .features-section {
    margin-top: 40px;
  }

  .feature-card {
    height: 180px;
    margin-bottom: 15px;
  }

  .feature-card h3 {
    font-size: 1.1rem;
  }

  .feature-card p {
    font-size: 0.9rem;
    padding: 0 8px;
  }

  .materials-section {
    margin-top: 60px;
  }

  .materials-section h2 {
    font-size: 1.5rem;
    margin-bottom: 30px;
  }

  .material-card {
    height: 140px;
    margin-bottom: 15px;
  }

  .material-card h4 {
    font-size: 1.1rem;
  }

  .material-card p {
    font-size: 0.85rem;
    padding: 0 5px;
  }

  .price {
    font-size: 1rem;
  }
}

@media (max-width: 480px) {
  .home {
    padding: 10px;
  }

  .hero-section {
    padding: 30px 10px;
    border-radius: 5px;
  }

  .hero-section h1 {
    font-size: 1.6rem;
    margin-bottom: 12px;
  }

  .hero-section p {
    font-size: 0.9rem;
    margin-bottom: 20px;
    line-height: 1.5;
  }

  .features-section {
    margin-top: 30px;
  }

  .feature-card {
    height: 160px;
    margin-bottom: 12px;
  }

  .feature-card h3 {
    font-size: 1rem;
    margin-bottom: 10px;
  }

  .feature-card p {
    font-size: 0.8rem;
    padding: 0 5px;
    line-height: 1.4;
  }

  .materials-section {
    margin-top: 50px;
  }

  .materials-section h2 {
    font-size: 1.3rem;
    margin-bottom: 25px;
  }

  .material-card {
    height: 120px;
    margin-bottom: 12px;
  }

  .material-card h4 {
    font-size: 1rem;
    margin-bottom: 8px;
  }

  .material-card p {
    font-size: 0.8rem;
    padding: 0 3px;
    margin-bottom: 10px;
  }

  .price {
    font-size: 0.9rem;
  }
}

@media (max-width: 360px) {
  .home {
    padding: 8px;
  }

  .hero-section {
    padding: 25px 8px;
  }

  .hero-section h1 {
    font-size: 1.4rem;
    line-height: 1.3;
  }

  .hero-section p {
    font-size: 0.85rem;
    line-height: 1.4;
  }

  .feature-card h3 {
    font-size: 0.95rem;
  }

  .feature-card p {
    font-size: 0.75rem;
  }

  .material-card h4 {
    font-size: 0.95rem;
  }

  .material-card p {
    font-size: 0.75rem;
  }

  .price {
    font-size: 0.85rem;
  }
}

/* 高分辨率屏幕优化 */
@media (min-width: 1200px) {
  .hero-section h1 {
    font-size: 3rem;
  }

  .hero-section p {
    font-size: 1.2rem;
  }

  .feature-card h3 {
    font-size: 1.3rem;
  }

  .feature-card p {
    font-size: 1rem;
  }

  .materials-section h2 {
    font-size: 2rem;
  }

  .material-card h4 {
    font-size: 1.3rem;
  }

  .material-card p {
    font-size: 0.95rem;
  }

  .price {
    font-size: 1.2rem;
  }
}
</style>