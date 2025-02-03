<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="login-header">
        <img src="../assets/BigLOGO.png" alt="Serbian Polytechnic University Logo" class="logo">
        <h2>学籍申请系统</h2>
        <p class="subtitle">国立坎特洛特大学学籍申请系统</p>
      </div>
      
      <div class="login-content">
        <p class="description">
          请使用您的 Linux Do 论坛账号登录以开始申请流程（低于二级将被拒绝）
        </p>

        <el-button
          type="primary"
          class="oauth-button"
          @click="handleOAuthLogin"
          size="large"
        >
          <el-icon class="icon"><Position /></el-icon>
          使用 Linux Do 账号登录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { Position } from '@element-plus/icons-vue'

const handleOAuthLogin = () => {
  const clientId = 'Nes7cKLZZn2tfdrcBJLun0EKCywfmdzr'
  const redirectUri = 'http://127.0.0.1:5173/oauth/callback'
  const state = Math.random().toString(36).substring(7)
  localStorage.setItem('oauth_state', state)
  const oauthUrl = `https://connect.linux.do/oauth2/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&response_type=code&scope=read write&state=${state}`
  window.location.href = oauthUrl
}
</script>

<style lang="scss" scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .login-card {
    width: 90%;
    max-width: 600px;
    padding: 30px;
    
    .login-header {
      text-align: center;
      margin-bottom: 40px;
      
      .logo {
        width: 400px;
        height: auto;
        margin-bottom: 20px;
      }
      
      h2 {
        color: #409EFF;
        font-size: 28px;
        margin: 0 0 10px 0;
      }

      .subtitle {
        color: #606266;
        font-size: 16px;
        margin: 0;
      }
    }
    
    .login-content {
      .description {
        color: #606266;
        text-align: center;
        margin-bottom: 30px;
        line-height: 1.6;
      }

      .oauth-button {
        width: 100%;
        height: 50px;
        font-size: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;

        .icon {
          font-size: 20px;
        }
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 480px) {
  .login-container {
    .login-card {
      width: 95%;
      padding: 20px;
      
      .login-header {
        margin-bottom: 30px;
        
        .logo {
          width: 100%;
          max-width: 300px;
          height: auto;
        }
        
        h2 {
          font-size: 24px;
        }

        .subtitle {
          font-size: 14px;
        }
      }

      .login-content {
        .description {
          font-size: 14px;
        }

        .oauth-button {
          height: 44px;
          font-size: 15px;
        }
      }
    }
  }
}
</style> 
