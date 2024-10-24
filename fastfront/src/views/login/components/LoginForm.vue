<template>
  <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules">
    <h1>登 录</h1>
    <el-form-item prop="username">
      <el-input
        v-model="loginForm.username"
        placeholder="邮箱/手机号"
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
        v-model="loginForm.password"
        placeholder="密码"
        type="password"
        show-password
        size="large"
        @keydown.enter="submitSignIn"
        class="input-field"
      >
        <template #prefix>
          <el-icon><Lock /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item>
      <div class="checkbox-container">
        <el-checkbox v-model="rememberMe">记住我</el-checkbox>
        <el-button type="primary" link @click="changeShowForm('forget')">
          忘记密码?
        </el-button>
        <el-button type="primary" link @click="changeShowForm('sign')">
          立即注册
        </el-button>
      </div>
    </el-form-item>
    <el-form-item>
      <el-button
        :loading="loading"
        class="submit-button"
        type="primary"
        @click.prevent="submitSignIn"
        >登 录</el-button
      >
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { ElForm } from "element-plus";
import store from "@/store";
import Cookies from "js-cookie";

const loginForm = reactive({
  username: "admin@admin.com",
  password: "123456",
});

const rememberMe = ref(false);

// 验证手机号或邮箱的规则
const validateUsername = (rule, value, callback) => {
  const phoneReg = /^1[3-9][0-9]\d{8}$/; // 手机号正则表达式
  const emailReg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/; // 邮箱正则表达式
  const regex = /@admin\.com$/;
  if (!value) {
    return callback(new Error("请输入手机号/邮箱"));
  } else if (
    phoneReg.test(value) ||
    emailReg.test(value) ||
    regex.test(value)
  ) {
    callback();
  } else {
    return callback(new Error("请输入有效的手机号或邮箱"));
  }
};

onMounted(() => {
  const rememberMeCookie = Cookies.get("rememberMe");
  if (rememberMeCookie) {
    const rememberMeData = JSON.parse(rememberMeCookie);
    loginForm.username = rememberMeData.username;
    loginForm.password = rememberMeData.password;
  }
});

const loginRules = {
  username: [{ required: true, validator: validateUsername, trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
};

const redirect = ref();
const otherQuery = ref({});
const router = useRouter();
const route = useRoute();

const loading = ref(false);
const loginFormRef = ref<InstanceType<typeof ElForm> | null>(null);
const submitSignIn = () => {
  loginFormRef.value?.validate((valid: boolean) => {
    if (valid) {
      emits("openDia", "login");
    }
  });
};

const handle = (captcha) => {
  loading.value = true;
  let params = { ...loginForm, ...captcha };
  store
    .user()
    .login(params)
    .then(() => {
      rememberMe.value
        ? Cookies.set("rememberMe", JSON.stringify(loginForm))
        : Cookies.remove("rememberMe");
      router.push({
        path: redirect.value || "/",
        query: otherQuery.value,
      });
    })
    .catch(() => {
      emits("clear");
    })
    .finally(() => {
      loading.value = false;
    });
};
defineExpose({ handle });
interface QueryType {
  [propname: string]: string;
}
const getOtherQuery = (query: QueryType) => {
  return Object.keys(query).reduce((acc: QueryType, cur) => {
    if (cur !== "redirect") {
      acc[cur] = query[cur];
    }
    return acc;
  }, {});
};

watch(
  () => route.query,
  (query) => {
    if (query) {
      redirect.value = query.redirect || "";
      otherQuery.value = getOtherQuery(query as QueryType);
    }
  },
  { immediate: true }
);

const emits = defineEmits(["changeShowForm", "openDia", "clear"]);
const changeShowForm = (value: string) => {
  emits("changeShowForm", value);
};
</script>

<style lang="scss" scoped>
.input-field {
  width: 340px;
}

.checkbox-container {
  width: 100%;
  display: flex;
  justify-content: space-between;
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
