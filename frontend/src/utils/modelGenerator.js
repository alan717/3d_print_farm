// 3D模型工具类 - 用于生成演示模型
export class ModelGenerator {
  // 生成一个简单的立方体STL数据
  static generateCubeSTL(size = 20) {
    const vertices = [
      // 前面
      [-size/2, -size/2, size/2],   // 0
      [size/2, -size/2, size/2],    // 1
      [size/2, size/2, size/2],     // 2
      [-size/2, size/2, size/2],    // 3
      
      // 后面
      [-size/2, -size/2, -size/2],  // 4
      [size/2, -size/2, -size/2],   // 5
      [size/2, size/2, -size/2],    // 6
      [-size/2, size/2, -size/2]    // 7
    ]
    
    // 立方体的12个三角形面
    const faces = [
      // 前面
      [0, 1, 2], [0, 2, 3],
      // 后面
      [4, 7, 6], [4, 6, 5],
      // 左面
      [0, 3, 7], [0, 7, 4],
      // 右面
      [1, 5, 6], [1, 6, 2],
      // 上面
      [3, 2, 6], [3, 6, 7],
      // 下面
      [0, 4, 5], [0, 5, 1]
    ]
    
    return this.generateSTLFromMesh(vertices, faces)
  }
  
  // 生成STL文件格式的二进制数据
  static generateSTLFromMesh(vertices, faces) {
    const numTriangles = faces.length
    const bufferSize = 80 + 4 + (numTriangles * 50) // STL文件格式
    const buffer = new ArrayBuffer(bufferSize)
    const view = new DataView(buffer)
    
    let offset = 0
    
    // 80字节的头部（可以留空）
    offset += 80
    
    // 三角形数量
    view.setUint32(offset, numTriangles, true)
    offset += 4
    
    // 为每个三角形写入数据
    faces.forEach(face => {
      const v1 = vertices[face[0]]
      const v2 = vertices[face[1]]
      const v3 = vertices[face[2]]
      
      // 计算法向量
      const normal = this.calculateNormal(v1, v2, v3)
      
      // 法向量 (12 bytes)
      view.setFloat32(offset, normal[0], true); offset += 4
      view.setFloat32(offset, normal[1], true); offset += 4
      view.setFloat32(offset, normal[2], true); offset += 4
      
      // 顶点1 (12 bytes)
      view.setFloat32(offset, v1[0], true); offset += 4
      view.setFloat32(offset, v1[1], true); offset += 4
      view.setFloat32(offset, v1[2], true); offset += 4
      
      // 顶点2 (12 bytes)
      view.setFloat32(offset, v2[0], true); offset += 4
      view.setFloat32(offset, v2[1], true); offset += 4
      view.setFloat32(offset, v2[2], true); offset += 4
      
      // 顶点3 (12 bytes)
      view.setFloat32(offset, v3[0], true); offset += 4
      view.setFloat32(offset, v3[1], true); offset += 4
      view.setFloat32(offset, v3[2], true); offset += 4
      
      // 属性字节计数 (2 bytes)
      view.setUint16(offset, 0, true); offset += 2
    })
    
    return new Uint8Array(buffer)
  }
  
  // 计算三角形法向量
  static calculateNormal(v1, v2, v3) {
    const edge1 = [v2[0] - v1[0], v2[1] - v1[1], v2[2] - v1[2]]
    const edge2 = [v3[0] - v1[0], v3[1] - v1[1], v3[2] - v1[2]]
    
    // 叉积
    const normal = [
      edge1[1] * edge2[2] - edge1[2] * edge2[1],
      edge1[2] * edge2[0] - edge1[0] * edge2[2],
      edge1[0] * edge2[1] - edge1[1] * edge2[0]
    ]
    
    // 归一化
    const length = Math.sqrt(normal[0] * normal[0] + normal[1] * normal[1] + normal[2] * normal[2])
    if (length > 0) {
      normal[0] /= length
      normal[1] /= length
      normal[2] /= length
    }
    
    return normal
  }
  
  // 创建演示文件
  static createDemoFile(type = 'cube', size = 20) {
    let data
    let filename
    
    switch (type) {
      case 'cube':
        data = this.generateCubeSTL(size)
        filename = 'demo_cube.stl'
        break
      default:
        data = this.generateCubeSTL(size)
        filename = 'demo_model.stl'
    }
    
    const blob = new Blob([data], { type: 'application/octet-stream' })
    return new File([blob], filename, { type: 'application/octet-stream' })
  }
  
  // 创建文件URL用于预览
  static createModelURL(file) {
    return URL.createObjectURL(file)
  }
  
  // 清理URL资源
  static revokeModelURL(url) {
    URL.revokeObjectURL(url)
  }
}