<template>
  <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules">
    <h1>注 册</h1>
    <el-form-item prop="username">
      <el-input
        v-model="registerForm.username"
        placeholder="邮箱/手机号"
        size="large"
        class="input-field"
      >
        <template #prefix>
          <el-icon><UserFilled /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="code">
      <el-input
        v-model="registerForm.code"
        placeholder="验证码"
        size="large"
        @keypress="onlyNumbers"
        class="input-field"
      >
        <template #append>
          <el-button
            :disabled="isButtonDisabled || isCounting"
            :class="{ 'custom-button': !isButtonDisabled && !isCounting }"
            @click="getCode"
          >
            {{ countdownText }}
          </el-button>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password">
      <el-input
        v-model="registerForm.password"
        placeholder="密码"
        type="password"
        show-password
        @paste.prevent
        @copy.prevent
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
        @paste.prevent
        @copy.prevent
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
      <el-button class="submit-button" type="primary" @click="submitRegister">
        注 册
      </el-button>
    </el-form-item>
    <el-form-item>
      <el-button
        class="submit-button"
        type="primary"
        @click="changeShowForm('login')"
      >
        返 回
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";

const registerForm = reactive({
  username: "",
  code: "",
  password: "",
  confirmPassword: "",
});

const isButtonDisabled = ref(true);
const isCounting = ref(false);
const countdownTime = ref(60);
const countdownText = ref("获取验证码");

const onlyNumbers = (event: KeyboardEvent) => {
  const key = event.key;
  if (!/^\d$/.test(key) && key !== "Backspace" && key !== "Tab") {
    event.preventDefault();
  }
};

// 验证手机号或邮箱的规则
const validateUsername = (rule, value, callback) => {
  const phoneReg = /^1[3-9][0-9]\d{8}$/; // 手机号正则表达式
  const emailReg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

  if (!value) {
    isButtonDisabled.value = true;
    return callback(new Error("请输入邮箱或手机号"));
  } else if (phoneReg.test(value) || emailReg.test(value)) {
    isButtonDisabled.value = false;
    callback();
  } else {
    isButtonDisabled.value = true;
    return callback(new Error("请输入有效的邮箱或手机号"));
  }
};

// 验证密码复杂度的规则
const validatePasswordComplexity = (rule, value, callback) => {
  const complexReg =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/; // 至少8个字符，包含大写字母、小写字母、数字和特殊字符
  if (!value) {
    return callback(new Error("请输入密码"));
  } else if (!complexReg.test(value)) {
    return callback(
      new Error("密码至少8个字符，包含大、小写字母、数字和特殊字符")
    );
  }
  callback();
};

const registerRules = {
  username: [{ required: true, validator: validateUsername, trigger: "blur" }],
  password: [
    { required: true, validator: validatePasswordComplexity, trigger: "blur" },
  ],
  code: [{ required: true, message: "请输入验证码", trigger: "blur" }],
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

// 获取验证码
const getCode = () => {
  registerFormRef.value?.validateField("username", (valid) => {
    if (valid) {
      if (!isButtonDisabled.value && !isCounting.value) {
        isCounting.value = true;
        countdownText.value = `${countdownTime.value}秒`;
        localStorage.setItem("signCountdownTime", String(countdownTime.value));
        localStorage.setItem("isSignCounting", "true");

        const timer = setInterval(() => {
          countdownTime.value -= 1;
          countdownText.value = `${countdownTime.value}秒`;

          if (countdownTime.value <= 0) {
            clearInterval(timer);
            isCounting.value = false;
            isButtonDisabled.value = false;
            countdownText.value = "获取验证码";
            localStorage.removeItem("signCountdownTime");
            localStorage.removeItem("isSignCounting");
          } else {
            localStorage.setItem(
              "signCountdownTime",
              String(countdownTime.value)
            );
          }
        }, 1000);
      }
    }
  });
};

// 组件挂载时恢复倒计时状态
onMounted(() => {
  const savedCountdown = localStorage.getItem("signCountdownTime");
  const savedIsCounting = localStorage.getItem("isSignCounting");
  if (savedIsCounting === "true") {
    isCounting.value = true;
    isButtonDisabled.value = true;
    countdownTime.value = Number(savedCountdown);
    countdownText.value = `${countdownTime.value}秒`;

    const timer = setInterval(() => {
      countdownTime.value -= 1;
      countdownText.value = `${countdownTime.value}秒`;
      localStorage.setItem("signCountdownTime", String(countdownTime.value));
      if (countdownTime.value <= 0) {
        clearInterval(timer);
        isCounting.value = false;
        isButtonDisabled.value = false;
        countdownText.value = "获取验证码";
        localStorage.removeItem("signCountdownTime");
        localStorage.removeItem("isSignCounting");
      }
    }, 1000);
  }
});

const emits = defineEmits(["changeShowForm"]);

const changeShowForm = (type: string) => {
  emits("changeShowForm", type);
};
</script>

<style lang="scss" scoped>
.input-field {
  width: 340px;
}

.submit-button {
  width: 100%;
}

.custom-button {
  color: #000000 !important;
}
.button-container {
  width: 100%;
  display: flex;
  justify-content: space-between;
}
</style>
