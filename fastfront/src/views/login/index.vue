<template>
  <div class="login-container">
    <div class="lefts" style="flex: 1">
      <div class="top-row">
        <img class="logo" src="@/assets/logo.png" alt="logo" />
        <span>Admin</span>
      </div>
      <div class="bg-container">
        <BgSvg class="bg" />
        <h3>后台管理系统</h3>
      </div>
    </div>
    <div class="rights">
      <transition name="slide" mode="out-in">
        <component
          :is="showFormComponent"
          @changeShowForm="changeShowForm"
          @clear="clear"
          @closeDia="closeDia"
          @openDia="openDia"
          :key="showForm"
          ref="showRef"
        />
      </transition>
    </div>
    <Code style="border-radius: 20px" ref="codeRef" @handle="handle"></Code>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import BgSvg from "@/assets/login/login-box-bg.svg";

import Code from "./components/Code.vue";
import LoginForm from "./components/LoginForm.vue";
import SignForm from "./components/SignForm.vue";
import ForgetForm from "./components/ForgetForm.vue";

const showForm = ref("login");
const showFormComponent = computed(() => {
  switch (showForm.value) {
    case "sign":
      return SignForm;
    case "forget":
      return ForgetForm;
    default:
      return LoginForm;
  }
});

const codeRef = ref();
const openDia = (type: string) => {
  codeRef.value.open(type);
};
const closeDia = () => {
  codeRef.value.closeDia();
};
const clear = () => {
  codeRef.value.clear();
};

const showRef = ref();
const handle = (captcha: any) => {
  showRef.value.handle(captcha);
};
const changeShowForm = (value: string) => {
  showForm.value = value;
};
</script>

<style lang="scss">
@import "./login.scss";
</style>
