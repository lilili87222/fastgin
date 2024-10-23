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
          :key="showForm"
        />
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import BgSvg from "@/assets/login/login-box-bg.svg";

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
const changeShowForm = (value: string) => {
  showForm.value = value;
};
</script>

<style lang="scss">
@import "./login.scss";
</style>
