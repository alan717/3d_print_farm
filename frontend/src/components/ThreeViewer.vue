<template>
  <div class="three-viewer" ref="containerRef">
    <div v-if="!modelLoaded" class="loading-overlay">
      <el-icon class="loading-icon" size="40">
        <Loading />
      </el-icon>
      <p>加载3D模型中...</p>
    </div>
    
    <div v-if="error" class="error-overlay">
      <el-icon class="error-icon" size="40">
        <Warning />
      </el-icon>
      <p>{{ error }}</p>
      <el-button @click="retryLoad">重试</el-button>
    </div>
    
    <div class="controls" v-if="modelLoaded">
      <el-button-group>
        <el-button @click="resetView" size="small">重置视角</el-button>
        <el-button @click="toggleWireframe" size="small">
          {{ wireframe ? '实体模式' : '线框模式' }}
        </el-button>
        <el-button @click="toggleAutoRotate" size="small">
          {{ autoRotate ? '停止旋转' : '自动旋转' }}
        </el-button>
      </el-button-group>
      
      <div class="model-info">
        <span>顶点数: {{ modelStats.vertices }}</span>
        <span>面数: {{ modelStats.faces }}</span>
        <span>体积: {{ modelStats.volume }} cm³</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Loading, Warning } from '@element-plus/icons-vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js'
import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader.js'

export default {
  name: 'ThreeViewer',
  components: {
    Loading,
    Warning
  },
  props: {
    modelUrl: {
      type: String,
      default: ''
    },
    modelType: {
      type: String,
      default: 'stl', // stl or obj
      validator: (value) => ['stl', 'obj'].includes(value)
    },
    width: {
      type: Number,
      default: 400
    },
    height: {
      type: Number,
      default: 300
    }
  },
  emits: ['model-loaded', 'model-error', 'model-stats'],
  setup(props, { emit }) {
    const containerRef = ref()
    const modelLoaded = ref(false)
    const error = ref('')
    const wireframe = ref(false)
    const autoRotate = ref(false)
    
    const modelStats = ref({
      vertices: 0,
      faces: 0,
      volume: 0
    })
    
    let scene, camera, renderer, controls, model
    let animationId
    
    const initThree = () => {
      // 创建场景
      scene = new THREE.Scene()
      scene.background = new THREE.Color(0xf0f0f0)
      
      // 创建相机
      camera = new THREE.PerspectiveCamera(
        75, 
        props.width / props.height, 
        0.1, 
        1000
      )
      camera.position.set(5, 5, 5)
      
      // 创建渲染器
      renderer = new THREE.WebGLRenderer({ antialias: true })
      renderer.setSize(props.width, props.height)
      renderer.shadowMap.enabled = true
      renderer.shadowMap.type = THREE.PCFSoftShadowMap
      
      // 添加到容器
      containerRef.value.appendChild(renderer.domElement)
      
      // 创建控制器
      controls = new OrbitControls(camera, renderer.domElement)
      controls.enableDamping = true
      controls.dampingFactor = 0.05
      
      // 添加灯光
      addLights()
      
      // 开始渲染循环
      animate()
    }
    
    const addLights = () => {
      // 环境光
      const ambientLight = new THREE.AmbientLight(0x404040, 0.6)
      scene.add(ambientLight)
      
      // 方向光
      const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8)
      directionalLight.position.set(10, 10, 5)
      directionalLight.castShadow = true
      directionalLight.shadow.mapSize.width = 2048
      directionalLight.shadow.mapSize.height = 2048
      scene.add(directionalLight)
      
      // 点光源
      const pointLight = new THREE.PointLight(0xffffff, 0.5)
      pointLight.position.set(-10, 10, -10)
      scene.add(pointLight)
    }
    
    const loadModel = async () => {
      if (!props.modelUrl) return
      
      try {
        modelLoaded.value = false
        error.value = ''
        
        // 移除之前的模型
        if (model) {
          scene.remove(model)
        }
        
        let loader
        if (props.modelType === 'stl') {
          loader = new STLLoader()
        } else if (props.modelType === 'obj') {
          loader = new OBJLoader()
        }
        
        const geometry = await new Promise((resolve, reject) => {
          loader.load(
            props.modelUrl,
            (loadedGeometry) => {
              if (props.modelType === 'obj') {
                // OBJ加载器返回的是Group对象
                const mesh = loadedGeometry.children[0]
                resolve(mesh.geometry)
              } else {
                // STL加载器直接返回geometry
                resolve(loadedGeometry)
              }
            },
            (progress) => {
              console.log('Loading progress:', progress)
            },
            (error) => {
              reject(error)
            }
          )
        })
        
        // 计算模型统计信息
        calculateModelStats(geometry)
        
        // 创建材质
        const material = new THREE.MeshLambertMaterial({
          color: 0x00aaff,
          transparent: true,
          opacity: 0.9
        })
        
        // 创建网格
        model = new THREE.Mesh(geometry, material)
        model.castShadow = true
        model.receiveShadow = true
        
        // 居中模型
        centerModel(geometry)
        
        scene.add(model)
        
        modelLoaded.value = true
        emit('model-loaded', model)
        emit('model-stats', modelStats.value)
        
      } catch (err) {
        error.value = '模型加载失败: ' + err.message
        emit('model-error', err)
        console.error('Model loading error:', err)
      }
    }
    
    const calculateModelStats = (geometry) => {
      geometry.computeBoundingBox()
      
      // 计算顶点和面数
      const positions = geometry.attributes.position
      modelStats.value.vertices = positions.count
      modelStats.value.faces = positions.count / 3
      
      // 计算体积（简化方法）
      const box = geometry.boundingBox
      const size = new THREE.Vector3()
      box.getSize(size)
      modelStats.value.volume = (size.x * size.y * size.z * 1000).toFixed(2) // 转换为cm³
    }
    
    const centerModel = (geometry) => {
      geometry.computeBoundingBox()
      const box = geometry.boundingBox
      const center = new THREE.Vector3()
      box.getCenter(center)
      
      // 将模型移动到原点
      geometry.translate(-center.x, -center.y, -center.z)
      
      // 调整相机位置以适应模型大小
      const size = new THREE.Vector3()
      box.getSize(size)
      const maxDim = Math.max(size.x, size.y, size.z)
      
      camera.position.set(maxDim * 1.5, maxDim * 1.5, maxDim * 1.5)
      controls.target.set(0, 0, 0)
      controls.update()
    }
    
    const animate = () => {
      animationId = requestAnimationFrame(animate)
      
      if (autoRotate.value && model) {
        model.rotation.y += 0.01
      }
      
      controls.update()
      renderer.render(scene, camera)
    }
    
    const resetView = () => {
      if (model) {
        model.rotation.set(0, 0, 0)
        camera.position.set(5, 5, 5)
        controls.target.set(0, 0, 0)
        controls.update()
      }
    }
    
    const toggleWireframe = () => {
      if (model) {
        wireframe.value = !wireframe.value
        model.material.wireframe = wireframe.value
      }
    }
    
    const toggleAutoRotate = () => {
      autoRotate.value = !autoRotate.value
    }
    
    const retryLoad = () => {
      loadModel()
    }
    
    const handleResize = () => {
      if (camera && renderer) {
        camera.aspect = props.width / props.height
        camera.updateProjectionMatrix()
        renderer.setSize(props.width, props.height)
      }
    }
    
    // 监听模型URL变化
    watch(() => props.modelUrl, () => {
      if (props.modelUrl) {
        loadModel()
      }
    }, { immediate: true })
    
    // 监听尺寸变化
    watch([() => props.width, () => props.height], () => {
      handleResize()
    })
    
    onMounted(() => {
      initThree()
    })
    
    onBeforeUnmount(() => {
      if (animationId) {
        cancelAnimationFrame(animationId)
      }
      if (renderer) {
        renderer.dispose()
      }
      if (controls) {
        controls.dispose()
      }
    })
    
    return {
      containerRef,
      modelLoaded,
      error,
      wireframe,
      autoRotate,
      modelStats,
      resetView,
      toggleWireframe,
      toggleAutoRotate,
      retryLoad
    }
  }
}
</script>

<style scoped>
.three-viewer {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.loading-overlay,
.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.loading-icon {
  animation: spin 1s linear infinite;
  color: #409eff;
  margin-bottom: 10px;
}

.error-icon {
  color: #f56c6c;
  margin-bottom: 10px;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.controls {
  position: absolute;
  top: 10px;
  left: 10px;
  right: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.9);
  padding: 10px;
  border-radius: 5px;
  z-index: 5;
}

.model-info {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #666;
}

.model-info span {
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 3px;
}
</style>