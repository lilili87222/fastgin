<template>
  <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules">
    <h1>登 录</h1>
    <el-form-item prop="username">
      <el-input
        v-model="loginForm.username"
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
        <el-button type="primary" link @click="toForget"> 忘记密码? </el-button>
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
    <el-form-item>
      <div class="button-container">
        <el-button @click="changeShowForm">注 册</el-button>
      </div>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import type { ElForm } from "element-plus";
import store from "@/store";

const loginForm = reactive({
  username: "admin",
  password: "123456",
});

const rememberMe = ref(false);

const loginRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
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
      loading.value = true;
      store
        .user()
        .login(loginForm)
        .then(() => {
          router.push({
            path: redirect.value || "/",
            query: otherQuery.value,
          });
        })
        .finally(() => {
          loading.value = false;
        });
    }
  });
};

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

const emits = defineEmits(["changeShowForm"]);
const changeShowForm = () => {
  emits("changeShowForm", "sign");
};

const toForget = () => {
  emits("changeShowForm", "forget");
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
