<template>
  <div class="register-container">
    <el-card class="register-card">
      <div class="register-header">
        <img src="../assets/BigLOGO.png" alt="Serbian Polytechnic University Logo" class="logo">
        <h2>邮箱注册</h2>
        <p class="subtitle">请完善您的个人信息</p>
      </div>
      
      <div class="register-content">
        <el-form 
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          class="register-form"
        >
          <div class="form-section">
            <h3 class="section-title">基本信息</h3>
            <div class="name-row">
              <el-form-item label="姓氏" prop="lastName" class="name-item">
                <el-input
                  v-model="form.lastName"
                  placeholder="请输入您的姓氏"
                  clearable
                >
                  <template #prefix>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </el-form-item>

              <el-form-item label="名字" prop="firstName" class="name-item">
                <el-input
                  v-model="form.firstName"
                  placeholder="请输入您的名字"
                  clearable
                >
                  <template #prefix>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
            </div>
          </div>

          <div class="form-section">
            <h3 class="section-title">邮箱信息</h3>
            <el-form-item label="学校邮箱" prop="schoolEmail">
              <div class="school-email-input">
                <el-input
                  v-model="form.schoolEmailPrefix"
                  placeholder="请输入邮箱前缀"
                  clearable
                >
                  <template #prefix>
                    <el-icon><Message /></el-icon>
                  </template>
                </el-input>
                <span class="email-suffix">@polyu.edu.rs</span>
              </div>
              <div class="email-tip">系统将自动添加 @polyu.edu.rs 后缀</div>
            </el-form-item>

            <el-form-item label="辅助邮箱" prop="backupEmail">
              <el-input
                v-model="form.backupEmail"
                placeholder="请输入您的辅助邮箱"
                clearable
              >
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </div>

          <el-form-item>
            <el-button type="primary" class="submit-button" @click="handleSubmit">
              提交注册
            </el-button>
          </el-form-item>
        </el-form>

        <div v-if="userInfo" class="user-info">
          <h3 class="section-title">账号信息</h3>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="Linux Do 用户名">{{ userInfo.username }}</el-descriptions-item>
            <el-descriptions-item label="信任等级">
              <el-tag :type="getTrustLevelType(userInfo.trust_level)">
                Level {{ userInfo.trust_level }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-card>

    <!-- 注册成功对话框 -->
    <el-dialog
      v-model="successDialogVisible"
      title="注册成功"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <div class="success-dialog-content">
        <el-icon class="success-icon" color="#67C23A"><CircleCheckFilled /></el-icon>
        <h3>您的账号已创建成功！</h3>
        <div class="account-info">
          <p><strong>邮箱：</strong>{{ registrationResult?.email }}</p>
          <p><strong>密码：</strong>{{ registrationResult?.password }}</p>
          <div class="password-tip">请保存好您的密码，首次登录后请及时修改</div>
        </div>
        <div class="dialog-actions">
          <el-button type="primary" @click="handleGoogleLogin">
            <el-icon><Position /></el-icon>
            前往 Google 登录
          </el-button>
          <el-button @click="router.push('/')">返回首页</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Message, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const formRef = ref(null)
const userInfo = ref(null)
const successDialogVisible = ref(false)
const registrationResult = ref(null)

const form = ref({
  lastName: '',
  firstName: '',
  schoolEmailPrefix: '',
  backupEmail: ''
})

// 计算完整的学校邮箱
const schoolEmail = computed(() => {
  return form.value.schoolEmailPrefix ? form.value.schoolEmailPrefix + '@polyu.edu.rs' : ''
})

const validateName = (rule, value, callback) => {
  if (!value) {
    callback(new Error('此项为必填项'))
  } else if (value.length < 2) {
    callback(new Error('长度至少为2个字符'))
  } else {
    callback()
  }
}

const validateSchoolEmail = (rule, value, callback) => {
  if (!form.value.schoolEmailPrefix) {
    callback(new Error('请输入邮箱前缀'))
  } else {
    // 验证前缀是否只包含允许的字符
    const prefixRegex = /^[a-zA-Z0-9._-]+$/
    if (!prefixRegex.test(form.value.schoolEmailPrefix)) {
      callback(new Error('邮箱前缀只能包含字母、数字、点、下划线和横线'))
    } else {
      callback()
    }
  }
}

const validateBackupEmail = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请输入邮箱'))
  } else {
    const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/
    if (!emailRegex.test(value)) {
      callback(new Error('请输入有效的邮箱地址'))
    } else if (value.toLowerCase().endsWith('polyu.edu.rs')) {
      callback(new Error('备用邮箱不能使用学校邮箱'))
    } else {
      callback()
    }
  }
}

const rules = {
  lastName: [
    { required: true, validator: validateName, trigger: 'blur' }
  ],
  firstName: [
    { required: true, validator: validateName, trigger: 'blur' }
  ],
  schoolEmailPrefix: [
    { required: true, validator: validateSchoolEmail, trigger: 'blur' }
  ],
  backupEmail: [
    { required: true, validator: validateBackupEmail, trigger: 'blur' }
  ]
}

const getTrustLevelType = (level) => {
  if (level <= 1) {
    return 'danger'
  } else {
    return 'success'
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid, fields) => {
    if (valid) {
      try {
        const token = localStorage.getItem('token')
        const tokenType = localStorage.getItem('token_type')
        
        if (!token || !tokenType) {
          throw new Error('未登录')
        }

        const registerData = {
          lastName: form.value.lastName,
          firstName: form.value.firstName,
          schoolEmail: schoolEmail.value,
          backupEmail: form.value.backupEmail
        }

        const response = await fetch('/api/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `${tokenType} ${token}`
          },
          body: JSON.stringify(registerData)
        })

        const data = await response.json()

        if (!response.ok) {
          throw new Error(data.error || '注册失败')
        }

        ElMessage.success('注册成功')
        registrationResult.value = {
          email: schoolEmail.value,
          password: data.password
        }
        successDialogVisible.value = true
      } catch (error) {
        console.error('注册失败:', error)
        ElMessage.error(error.message)
      }
    } else {
      console.log('验证失败:', fields)
    }
  })
}

const handleGoogleLogin = () => {
  window.open('https://accounts.google.com/ServiceLogin?Email=' + encodeURIComponent(registrationResult.value.email) + '&continue=https://mail.google.com/mail/', '_blank')
}

onMounted(async () => {
  console.log('RegisterView mounted')
  
  try {
    // 检查token
    const token = localStorage.getItem('token')
    if (!token) {
      throw new Error('未登录')
    }

    // 从路由状态获取用户信息
    const routeState = router.currentRoute.value.state
    if (!routeState || !routeState.user) {
      // 如果没有用户信息，需要重新获取
      const response = await fetch('/api/user/info', {
        headers: {
          'Authorization': `${localStorage.getItem('token_type')} ${token}`
        }
      })

      if (!response.ok) {
        throw new Error('获取用户信息失败')
      }

      const data = await response.json()
      if (data.error) {
        throw new Error(data.error)
      }

      userInfo.value = data.user
    } else {
      userInfo.value = routeState.user
    }

    // 验证用户信息的完整性
    if (!userInfo.value.username || typeof userInfo.value.trust_level !== 'number') {
      throw new Error('用户信息不完整')
    }

    // 如果信任等级不足，跳转回首页
    if (userInfo.value.trust_level < 2) {
      throw new Error('信任等级不足')
    }
    
  } catch (error) {
    console.error('初始化失败:', error.message)
    ElMessage.error(error.message)
    // 清除token
    localStorage.removeItem('token')
    localStorage.removeItem('token_type')
    // 延迟跳转，让用户看到错误消息
    setTimeout(() => {
      router.push('/')
    }, 1500)
  }
})
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  
  .register-card {
    width: 90%;
    max-width: 600px;
    padding: 40px;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    
    .register-header {
      text-align: center;
      margin-bottom: 40px;
      
      .logo {
        width: 180px;
        height: auto;
        margin-bottom: 24px;
        filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
      }
      
      h2 {
        color: #409EFF;
        font-size: 28px;
        margin: 0 0 12px 0;
        font-weight: 600;
      }

      .subtitle {
        color: #606266;
        font-size: 16px;
        margin: 0;
        letter-spacing: 1px;
      }
    }
    
    .register-content {
      .register-form {
        margin-bottom: 30px;

        .form-section {
          background: #f8f9fa;
          border-radius: 8px;
          padding: 24px;
          margin-bottom: 24px;

          .section-title {
            color: #303133;
            font-size: 18px;
            margin: 0 0 20px 0;
            font-weight: 500;
            position: relative;
            padding-left: 12px;

            &::before {
              content: '';
              position: absolute;
              left: 0;
              top: 50%;
              transform: translateY(-50%);
              width: 4px;
              height: 16px;
              background: #409EFF;
              border-radius: 2px;
            }
          }
        }

        .name-row {
          display: flex;
          gap: 20px;
          
          .name-item {
            flex: 1;
          }
        }

        .school-email-input {
          display: flex;
          align-items: center;
          gap: 12px;

          .el-input {
            flex: 1;

            :deep(.el-input__wrapper) {
              box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
            }
          }

          .email-suffix {
            color: #606266;
            font-size: 14px;
            white-space: nowrap;
            user-select: none;
            background: #f4f4f5;
            padding: 8px 12px;
            border-radius: 4px;
          }
        }

        .email-tip {
          font-size: 12px;
          color: #909399;
          margin-top: 8px;
          padding-left: 4px;
        }

        .submit-button {
          width: 100%;
          height: 44px;
          font-size: 16px;
          font-weight: 500;
          letter-spacing: 1px;
          margin-top: 12px;
          background: linear-gradient(135deg, #409EFF 0%, #3a8ee6 100%);
          border: none;
          box-shadow: 0 2px 6px rgba(64, 158, 255, 0.3);
          transition: all 0.3s ease;

          &:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(64, 158, 255, 0.4);
          }
        }
      }

      .user-info {
        background: #f8f9fa;
        border-radius: 8px;
        padding: 24px;

        .section-title {
          color: #303133;
          font-size: 18px;
          margin: 0 0 20px 0;
          font-weight: 500;
          position: relative;
          padding-left: 12px;

          &::before {
            content: '';
            position: absolute;
            left: 0;
            top: 50%;
            transform: translateY(-50%);
            width: 4px;
            height: 16px;
            background: #409EFF;
            border-radius: 2px;
          }
        }

        :deep(.el-descriptions) {
          background: white;
          border-radius: 4px;
          overflow: hidden;
        }

        :deep(.el-tag) {
          margin: 0;
          min-width: 80px;
          text-align: center;
          padding: 4px 8px;
        }
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 480px) {
  .register-container {
    padding: 16px;
    
    .register-card {
      padding: 24px;
      
      .register-header {
        margin-bottom: 30px;
        
        .logo {
          width: 140px;
        }
        
        h2 {
          font-size: 24px;
        }

        .subtitle {
          font-size: 14px;
        }
      }

      .register-content {
        .register-form {
          .form-section {
            padding: 16px;

            .name-row {
              flex-direction: column;
              gap: 0;
            }
          }
        }
      }
    }
  }
}

.success-dialog-content {
  text-align: center;
  padding: 20px 0;

  .success-icon {
    font-size: 48px;
    margin-bottom: 16px;
  }

  h3 {
    color: #303133;
    margin: 0 0 24px 0;
  }

  .account-info {
    text-align: left;
    background: #f8f9fa;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 24px;

    p {
      margin: 8px 0;
      word-break: break-all;
    }

    .password-tip {
      color: #E6A23C;
      font-size: 12px;
      margin-top: 8px;
    }
  }

  .dialog-actions {
    display: flex;
    gap: 12px;
    justify-content: center;

    .el-button {
      min-width: 120px;

      .el-icon {
        margin-right: 4px;
      }
    }
  }
}

:deep(.el-dialog__header) {
  margin-right: 0;
  border-bottom: 1px solid #dcdfe6;
}
</style> 