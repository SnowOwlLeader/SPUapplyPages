<template>
  <div class="callback-container">
    <el-card class="callback-card">
      <template #header>
        <div class="callback-header">
          <el-icon><Loading /></el-icon>
          <span>正在处理登录请求...</span>
        </div>
      </template>
      <div class="callback-content">
        <el-progress type="circle" :percentage="progress" />
        <p>{{ message }}</p>
        <div v-if="userInfo" class="user-info">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
            <el-descriptions-item label="姓名">{{ userInfo.name }}</el-descriptions-item>
            <el-descriptions-item label="信任等级">
              <el-tag :type="getTrustLevelType(userInfo.trust_level)">
                Level {{ userInfo.trust_level }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const progress = ref(0)
const message = ref('正在验证授权信息...')
const userInfo = ref(null)

const getTrustLevelType = (level) => {
  if (level <= 1) {
    return 'danger'
  } else {
    return 'success'
  }
}

onMounted(async () => {
  const code = new URLSearchParams(window.location.search).get('code')
  const state = new URLSearchParams(window.location.search).get('state')
  const savedState = localStorage.getItem('oauth_state')
  
  if (!code) {
    message.value = '授权失败：未获取到授权码'
    ElMessage.error('授权失败：未获取到授权码')
    setTimeout(() => router.push('/'), 2000)
    return
  }

  if (!state || state !== savedState) {
    message.value = '授权失败：state不匹配'
    ElMessage.error('授权失败：state不匹配')
    setTimeout(() => router.push('/'), 2000)
    return
  }

  localStorage.removeItem('oauth_state')
  progress.value = 30
  message.value = '正在获取访问令牌...'

  try {
    const response = await fetch('/api/oauth/callback?' + new URLSearchParams({
      code: code
    }), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    if (!response.ok) {
      throw new Error('服务器响应错误')
    }

    progress.value = 60
    message.value = '正在保存登录信息...'
    
    const data = await response.json()
    if (data.error) {
      throw new Error(data.error)
    }

    // 保存token
    localStorage.setItem('token', data.token)
    localStorage.setItem('token_type', data.type)
    
    // 只在当前页面使用用户信息，不保存到localStorage
    userInfo.value = data.user
    
    progress.value = 100
    
    // 判断信任等级
    if (data.user.trust_level < 2) {
      message.value = '信任等级不足，需要 Level 2 以上才能申请'
      ElMessage.error('信任等级不足，需要 Level 2 以上才能申请')
      setTimeout(() => {
        localStorage.removeItem('token')
        localStorage.removeItem('token_type')
        router.push('/')
      }, 3000)
    } else {
      message.value = '登录成功！正在跳转到注册页面...'
      ElMessage.success('登录成功')
      // 将用户信息作为路由参数传递
      setTimeout(() => {
        router.push({
          name: 'register',
          state: { user: data.user }
        })
      }, 1500)
    }
    
  } catch (error) {
    message.value = '登录失败：' + error.message
    ElMessage.error('登录失败：' + error.message)
  }
})
</script>

<style lang="scss" scoped>
.callback-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .callback-card {
    width: 90%;
    max-width: 400px;
    
    .callback-header {
      display: flex;
      align-items: center;
      gap: 10px;
      
      .el-icon {
        font-size: 20px;
        animation: spin 1s linear infinite;
      }
    }
    
    .callback-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 20px;
      padding: 20px;
      
      p {
        margin: 0;
        text-align: center;
        color: #606266;
      }

      .user-info {
        width: 100%;
        margin-top: 20px;

        :deep(.el-tag) {
          margin: 0;
          min-width: 80px;
          text-align: center;
        }
      }
    }
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style> 