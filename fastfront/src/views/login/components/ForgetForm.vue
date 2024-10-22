<template>
  <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules">
    <h1>忘记密码</h1>
    <el-form-item prop="username">
      <el-input
        v-model="registerForm.username"
        placeholder="用户名"
        size="large"
        class="input-field"
      >
        <template #prefix>
          <el-icon><UserFilled /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password">
      <el-input
        v-model="registerForm.password"
        placeholder="密码"
        type="password"
        show-password
        size="large"
        class="input-field"
      >
        <template #prefix>
          <el-icon><Lock /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="confirmPassword">
      <el-input
        v-model="registerForm.confirmPassword"
        placeholder="确认密码"
        type="password"
        show-password
        size="large"
        @keydown.enter="submitRegister"
        class="input-field"
      >
        <template #prefix>
          <el-icon><Lock /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item>
      <el-button class="submit-button" type="primary" @click="submitRegister"
        >注 册</el-button
      >
    </el-form-item>
    <el-form-item>
      <div class="button-container">
        <el-button @click="changeShowForm">已有账号，登录</el-button>
      </div>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";

const registerForm = reactive({
  username: "",
  password: "",
  confirmPassword: "",
});

const registerRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
  confirmPassword: [
    { required: true, message: "请确认密码", trigger: "blur" },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error("两次输入的密码不一致"));
        } else {
          callback();
        }
      },
      trigger: "blur",
    },
  ],
};

const registerFormRef = ref();
const submitRegister = () => {
  registerFormRef.value?.validate((valid) => {
    if (valid) {
      console.log("注册", registerForm);
    }
  });
};

const emits = defineEmits(["changeShowForm"]);
const changeShowForm = () => {
  emits("changeShowForm", "login");
};
</script>

<style lang="scss" scoped>
.input-field {
  width: 340px;
}

.submit-button {
  width: 100%;
}

.button-container {
  width: 100%;
  display: flex;
  justify-content: space-between;
}
</style>
